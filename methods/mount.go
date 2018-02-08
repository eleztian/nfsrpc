package methods

import (
	"bs-2018/mynfs/nfsrpc"
	"errors"
	"net/rpc"
	"fmt"
)

// MOUNT
// RFC 1813 Section 5.0

const (
	MOUNT_PROG = 100005
	MOUNT_VERS = 3
)

type MountMethod int

// mount
const (
	MOUNTPROC3_NULL MountMethod = iota
	MOUNTPROC3_MNT
	MOUNTPROC3_DUMP
	MOUNTPROC3_UMNT
	MOUNTPROC3_UMNTALL
	MOUNTPROC3_EXPORT
)

// mount error
const (
	MNT3_OK             = 0     // no error
	MNT3ERR_PERM        = 1     // Not owner
	MNT3ERR_NOENT       = 2     // No such file or directory
	MNT3ERR_IO          = 5     // I/O error
	MNT3ERR_ACCES       = 13    // Permission denied
	MNT3ERR_NOTDIR      = 20    // Not a directory
	MNT3ERR_INVAL       = 22    // Invalid argument
	MNT3ERR_NAMETOOLONG = 63    // Filename too long
	MNT3ERR_NOTSUPP     = 10004 // Operation not supported
	MNT3ERR_SERVERFAULT = 10006 // A failure on the server
)

const (
	// The MOUNT service uses AUTH_NONE in the NULL procedure.
	AUTH_WAY   = nfsrpc.AuthNone
	MNTPATHLEN = 1024 // Maximum bytes in a path name
	MNTNAMLEN  = 255  // Maximum bytes in a name
	FHSIZE3    = 64   // Maximum bytes in a V3 file handle
)

/*

   typedef opaque fhandle3<FHSIZE3>;
   typedef string dirpath<MNTPATHLEN>;
   typedef string name<MNTNAMLEN>;

*/

type MountRes3OK struct {
	FHandle []byte
	flavors nfsrpc.AuthFlavor
}

// MOUNTPROC3_DUMP return
type Dump_result struct {
	name    string // max MNTNAMLEN
	dirpath string // max MNTPATHLEN
}

type DumpReply struct {
	D    Dump_result
	Next *getDumpReply `xdr:"optional"`
}

type getDumpReply struct {
	Next *DumpReply `xdr:"optional"`
}

func Dump(client *rpc.Client) ([]Dump_result, error) {
	var DumpR []Dump_result
	var result getDumpReply

	if client == nil {
		return nil, errors.New("Could not create pmap client")
	}
	defer client.Close()

	err := client.Call("Mount.Dump", nil, &result)
	if err != nil {
		return nil, err
	}

	if result.Next != nil {
		trav := result.Next
		for {
			entry := Dump_result(trav.D)
			DumpR = append(DumpR, entry)
			trav = trav.Next.Next
			if trav == nil {
				break
			}
		}
	}
	return DumpR, nil
}

type GroupNode struct {
	Name string
	Next *GroupNode `xdr:"optional"`
}

type ExportNode struct {
	DirPath string
	Groups  []string
}

type ExportList struct {
	Node ExportNode
	Next *ExportList `xdr:"optional"`
}

type GetExportReply struct {
	Next *ExportList `xdr:"optional"`
}

func Export(client *rpc.Client) ([]ExportNode, error) {
	if client == nil {
		return nil, errors.New("Could not create pmap client")
	}
	defer client.Close()
	var result GetExportReply
	err := client.Call("Mount.Export", nil, &result)
	if err != nil {
		return nil, err
	}
	var Exports []ExportNode
	fmt.Println(*result.Next)
	if result.Next != nil {
		trav := result.Next
		for {
			entry := ExportNode(trav.Node)
			Exports = append(Exports, entry)
			trav = trav.Next
			if trav == nil {
				break
			}
		}
	}
	return Exports, nil
}
