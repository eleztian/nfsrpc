package nfsrpc

/*
  rfc5531 section-8.4.2
*/

// A RPC message can be of two types: call or reply
const (
	RPCMsgTypeCall  MsgType = 0
	RPCMsgTypeReply MsgType = 1
)

// RPCProtocolVersion is the version of RPC protocol as described in RFC 5531
const RPCProtocolVersion = 2

// RPCMsg represents a complete RPC message (call or reply)
// union switch (msg_type mtype) {
//		case CALL:
//			call_body cbody;
//		case REPLY:
//			reply_body rbody;
// } body;
type RPCMsg struct {
	Xid   uint32
	Type  MsgType   `xdr:"union"`
	CBody CallBody  `xdr:"unioncase=0"`
	RBody ReplyBody `xdr:"unioncase=1"`
}

// MsgType is an enumeration representing the type of RPC message
type MsgType int32

// CallBody represents the body of a RPC Call
type CallBody struct {
	RPCVersion uint32     // must be equal to 2
	Program    uint32     // Remote program
	Version    uint32     // Remote program's version
	Procedure  uint32     // Procedure number
	Cred       OpaqueAuth // Authentication credential
	Verf       OpaqueAuth // Authentication verifier
}

// OpaqueAuth is a structure with AuthFlavor enumeration followed by up to
// 400 bytes that are opaque to (uninterpreted by) the RPC protocol
// implementation.
type OpaqueAuth struct {
	Flavor AuthFlavor
	Body   []byte `xdr:"opaque=true"`
}

const OpaueBodyMaxLength = 400

// AuthFlavor represents the type of authentication used
type AuthFlavor int32

// Sun-assigned authentication flavor numbers
const (
	AuthNone  AuthFlavor = iota // No authentication
	AuthSys                     // Unix style (uid+gids)
	AuthShort                   // Short hand unix style
	AuthDh                      // DES style (encrypted timestamp)
	AuthKerb                    // Keberos Auth
	AuthRSA                     // RSA authentication
	RPCSecGss                   // GSS-based RPC security
)

type AuthsysParms struct {
	Stamp       uint32
	MachineName string // max 255
	Uid         uint32
	Gid         uint32
	Gids        uint32 // max 16 a counted array of groups that contain the caller as a member.
}

// ReplyBody represents a generic RPC reply to a `Call`
type ReplyBody struct {
	Stat   ReplyStat     `xdr:"union"`
	Areply AcceptedReply `xdr:"unioncase=0"`
	Rreply RejectedReply `xdr:"unioncase=1"`
}

// ReplyStat is an enumeration representing the type of reply
type ReplyStat int32

// A reply to a call message can take two forms: the message was either
// accepted or rejected
const (
	MsgAccepted ReplyStat = 0
	MsgDenied   ReplyStat = 1
)

// AcceptedReply contains reply accepted by the RPC server. Note that there
// could be an error even though the call was accepted.
type AcceptedReply struct {
	Verf         OpaqueAuth
	Stat         AcceptStat    `xdr:"union"`
	MismatchInfo MismatchReply `xdr:"unioncase=2"` // ProgMismatch
	// procedure-specific results start here
}

// AcceptStat is an enumeration representing the status of procedure called
type AcceptStat int32

// Given that a call message was accepted, the following is the status of an
// attempt to call a remote procedure
const (
	Success      AcceptStat = iota // RPC executed successfully
	ProgUnavail                    // Remote hasn't exported the program
	ProgMismatch                   // Remote can't support version number
	ProcUnavail                    // Program can't support procedure
	GarbageArgs                    // Procedure can't decode params
	SystemErr                      // Other errors
)

// MismatchReply is used in ProgMismatch and RPCMismatch cases to denote
// lowest and highest version of RPC version or program version supported
type MismatchReply struct {
	Low  uint32
	High uint32
}

// RejectedReply represents a reply to a call rejected by the RPC server. The
/// call can be ejected for two reasons: either the server is not running a
// compatible version of the RPC protocol (RPCMismatch) or the server rejects
// the identity of the caller (AuthError)
type RejectedReply struct {
	Stat         RejectStat    `xdr:"union"`
	MismatchInfo MismatchReply `xdr:"unioncase=0"` // RPCMismatch
	AuthStat     AuthStat      `xdr:"unioncase=1"` // AuthError
}

// RejectStat is an enumeration representing the reason for rejection
type RejectStat int32

// Why call was rejected
const (
	RPCMismatch RejectStat = 0 // RPC version number != 2
	AuthError   RejectStat = 1 // Remote can't authenticate caller
)

// AuthStat represents the reason for authentication failure
type AuthStat int32

// Why authentication failed
const (
	AuthOk AuthStat = iota // Success

	// failed at remote end
	AuthBadcred      // Bad credential (seal broken)
	AuthRejectedcred // Client must begin new session
	AuthBadverf      // Bad verifier (seal broken)
	AuthRejectedVerf // Verifier expired or replayed
	AuthTooweak      // Rejected for security reasons

	// failed locally
	AuthInvalidresp // Bogus response verifier
	AuthFailed      // Reason unknown

	// AUTH_KERB errors: deprecated. see [RFC2695]
	AuthKerbGenric // kerberos generic error
	AuthTimeexpire // time of credential expired
	AuthTktFile    // problem with ticket file
	AuthDecode     // can't decode authenticator
	AuthNetAddr    // wrong net address in ticket

	// RPCSEC_GSS GSS related errors
	RpcsesGssCredproblem // no credentials for user
	RpcsecGssCtxproblem  // problem with context
)
