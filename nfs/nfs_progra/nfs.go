package nfs_progra

import (
	"bs-2018/mynfs/nfsrpc"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func init() {
	methods := []string{
		"Null",
		"GetAttr",
		"SetAttr",
		"Lookup",
		"Access",
		"ReadLink",
		"Read",
		"Write",
		"Create",
		"MkDir",
		"SymLink",
		"MkNod",
		"Remove",
		"RmDir",
		"Rename",
		"Link",
		"ReadDir",
		"ReadDirPlus",
		"FsStat",
		"FsInfo",
		"PathConf",
		"Commit",
	}
	producerId := nfsrpc.ProcedureID{
		ProgramNumber:  NFS_PROG,
		ProgramVersion: NFS_VERS,
	}
	for id, procName := range methods {
		producerId.ProcedureNumber = uint32(id)
		if err := nfsrpc.RegisterProcedure(
			nfsrpc.Procedure{
				ID:   producerId,
				Name: "NFS." + procName,
			},
		); err != nil {
			panic(err)
		}
	}
	log.Println("NFS Register over.")
}

type Client struct {
	*rpc.Client
}

func NewNfsClient(host string) (*Client, error) {
	port, err := nfsrpc.PmapGetPort(host, NFS_PROG, NFS_VERS, nfsrpc.IPProtoTCP)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		return nil, err
	}

	client := nfsrpc.NewClient(conn, nil, nil)

	return &Client{client}, nil
}

func NewNFSClientAuth(host string, authWay nfsrpc.AuthFlavor, authBody interface{}) (*Client, error) {
	port, err := nfsrpc.PmapGetPort(host, NFS_PROG, NFS_VERS, nfsrpc.IPProtoTCP)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		return nil, err
	}

	auth, err := nfsrpc.Auth(authWay, authBody)
	if err != nil {
		return nil, err
	}
	client := nfsrpc.NewClient(conn, auth, nil)

	return &Client{client}, nil
}
