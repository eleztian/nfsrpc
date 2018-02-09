package nfs

import (
	"bs-2018/mynfs/methods"
	"bs-2018/mynfs/nfsrpc"
	"fmt"
	"net"
	"net/rpc"
	"testing"
	"bytes"
	"github.com/rasky/go-xdr/xdr2"
	"time"
)

var host = "111.231.215.178"

func GetNfsClient(host string) (*rpc.Client, error) {
	port, err := nfsrpc.PmapGetPort(host, NFS_PROG, NFS_VERS, nfsrpc.IPProtoTCP)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)

	// xdr格式编码
	if _, err := xdr.Marshal(buffer, &nfsrpc.AuthsysParms{
		1,
		"ubuntu",
		500,
		500,
		0,
	}); err != nil {
		return nil,err
	}
	//bs := buffer.Bytes()
	client := nfsrpc.NewClient(conn,nil,nil,
		//&nfsrpc.OpaqueAuth{
		//	Flavor: nfsrpc.AuthNone,
		//	Body:bs,
		//},
		//&nfsrpc.OpaqueAuth{
		//	Flavor: nfsrpc.AuthNone,
		//	Body:bs,
		//},
	)
	return client, nil
}

func GetHandle(name string) (*NFS_FH3, error) {
	c, _ := net.Dial("tcp", host+":1011")
	client := nfsrpc.NewClient(c, nil, nil)
	s, err := methods.Mnt(client, name)
	if err != nil {
		return nil, err
	}
	return &NFS_FH3{s.MountInfo.FHandle}, nil
}

func TestGetAttr(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	result := GetAttrRes{}
	err = client.Call("NFS.GetAttr",
		GetAttr3Args{fh.Data},
		&result)
	if err != nil {
		t.Error("NFS.Null: ", err)
	}
	fmt.Println(result)
}

func TestSetAttr(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/t.go")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fh)
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	result := SetAttr3Res{}
	err = client.Call("NFS.SetAttr",
		SetAttr3Aargs{
			Object: NFS_FH3{fh.Data},
			NewAttr: SAttr3{
				MTime: SetMTime{2, uint64(time.Now().Unix())},
			},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.ResOk.ObjWcc.After)
}

func TestLookUp(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fh)
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	result := LookUp3Res{}
	err = client.Call("NFS.Lookup",
		LookUp3Args{
			What:DirOpArgs3{Dir:*fh, Name:"t1.go"},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}