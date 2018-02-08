package nfsrpc

import (
	"bytes"
	"io"
	"net"
	"net/rpc"
	"sync"

	"errors"
	"github.com/rasky/go-xdr/xdr2"
)

type clientCodec struct {
	conn         io.ReadWriteCloser // network connection
	recordReader io.Reader          // reader for RPC record
	notifyClose  chan<- io.ReadWriteCloser

	// Sun RPC responses include Seq (XID) but not ServiceMethod (procedure
	// number). Go package net/rpc expects both. So we save ServiceMethod
	// when sending the request and look it up when filling rpc.Response
	pending *sync.Map // maps Seq (XID) to ServiceMethod
	//mutex   sync.Mutex        // protects pending
	//pending map[uint64]string // maps Seq (XID) to ServiceMethod
	cred *OpaqueAuth
	verf *OpaqueAuth
}

// NewClientCodec returns a new rpc.ClientCodec using Sun RPC on conn.
// If a non-nil channel is passed as second argument, the conn is sent on
// that channel when Close() is called on conn.
func NewClientCodec(conn io.ReadWriteCloser, cred *OpaqueAuth, verf *OpaqueAuth, notifyClose chan<- io.ReadWriteCloser) rpc.ClientCodec {
	if cred == nil {
		cred = &OpaqueAuth{}
	}
	if verf == nil {
		verf = &OpaqueAuth{}
	}
	return &clientCodec{
		conn:        conn,
		notifyClose: notifyClose,
		pending:     &sync.Map{},
		cred:        cred,
		verf:        verf,
	}
}

// NewClient returns a new rpc.Client which internally uses Sun RPC codec
func NewClient(conn io.ReadWriteCloser, cred *OpaqueAuth, verf *OpaqueAuth) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn, cred, verf, nil))
}

// Dial connects to a Sun-RPC server at the specified network address
func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, nil, nil), err
}

func (c *clientCodec) WriteRequest(req *rpc.Request, param interface{}) error {
	procedureID, ok := GetProcedureID(req.ServiceMethod)
	if !ok {
		return ErrProcUnavail
	}
	c.pending.Store(req.Seq, req.ServiceMethod)
	msg := RPCMsg{
		Xid:  uint32(req.Seq),
		Type: RPCMsgTypeCall,
		CBody: CallBody{
			RPCVersion: RPCProtocolVersion,
			Program:    procedureID.ProgramNumber,
			Version:    procedureID.ProgramVersion,
			Procedure:  procedureID.ProcedureNumber,
			Cred:       *c.cred,
			Verf:       *c.verf,
		},
	}
	buffer := new(bytes.Buffer)

	// xdr格式编码
	if _, err := xdr.Marshal(buffer, &msg); err != nil {
		return err
	}

	if param != nil {
		// // xdr格式编码 参数
		if _, err := xdr.Marshal(buffer, &param); err != nil {
			return err
		}
	}
	if _, err := Write(c.conn, buffer.Bytes()); err != nil {
		if err == io.EOF && c.notifyClose != nil {
			c.notifyClose <- c.conn // 等待已有连接关闭
		}
		return err
	}
	return nil
}

func (c *clientCodec) ReadResponseHeader(resp *rpc.Response) error {
	readBytes, err := Read(c.conn)
	if err != nil {
		if err == io.EOF && c.notifyClose != nil {
			c.notifyClose <- c.conn
		}
		return err
	}
	c.recordReader = bytes.NewReader(readBytes)

	// 解码消息到reply
	var reply RPCMsg
	if _, err = xdr.Unmarshal(c.recordReader, &reply); err != nil {
		return err
	}

	// 响应完成从队列中删除
	resp.Seq = uint64(reply.Xid)
	if t, ok := c.pending.Load(resp.Seq); ok {
		resp.ServiceMethod = t.(string)
		c.pending.Delete(resp.Seq)
	}

	if err := c.checkReplyForErr(&reply); err != nil {
		return err
	}

	return nil
}

func (c *clientCodec) checkReplyForErr(reply *RPCMsg) error {

	if reply.Type != RPCMsgTypeReply {
		return ErrInvalidRPCMessageType
	}

	switch reply.RBody.Stat {
	case MsgAccepted:
		switch reply.RBody.Areply.Stat {
		case Success:
		case ProgMismatch:
			return ErrProgMismatch{
				reply.RBody.Areply.MismatchInfo.Low,
				reply.RBody.Areply.MismatchInfo.High}
		case ProgUnavail:
			return ErrProgUnavail
		case ProcUnavail:
			return ErrProcUnavail
		case GarbageArgs:
			return ErrGarbageArgs
		case SystemErr:
			return ErrSystemErr
		default:
			return ErrInvalidMsgAccepted
		}
	case MsgDenied:
		switch reply.RBody.Rreply.Stat {
		case RPCMismatch:
			return ErrRPCMismatch{
				reply.RBody.Rreply.MismatchInfo.Low,
				reply.RBody.Rreply.MismatchInfo.High}
		case AuthError:
			switch reply.RBody.Rreply.AuthStat {
			case AuthBadcred:
				return errors.New("bad credential (seal broken)")
			case AuthBadverf:
				return errors.New(" bad verifier (seal broken)")
			case AuthRejectedcred:
				return errors.New("client must begin new session")
			case AuthRejectedVerf:
				return errors.New("verifier expired or replayed")
			case AuthTooweak:
				return errors.New("rejected for security reasons")
			case AuthInvalidresp:
				return errors.New("bogus response verifier")
			case AuthFailed:
				return errors.New("reason unknown")
			}
			return ErrAuthError
		default:
			return ErrInvalidMsgDeniedType
		}
	default:
		return ErrInvalidRPCRepyType
	}

	return nil
}

func (c *clientCodec) ReadResponseBody(result interface{}) error {
	if result == nil {
		// read and drain it out ?
		return nil
	}

	if _, err := xdr.Unmarshal(c.recordReader, &result); err != nil {
		return err
	}

	return nil
}

func (c *clientCodec) Close() error {
	return c.conn.Close()
}
