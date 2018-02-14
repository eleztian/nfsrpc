package nfsrpc

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
)

// MOUNT
// RFC 1813 Section 5.0

const (
	MOUNT_PROG = 100005
	MOUNT_VERS = 3
	Mount_Name = "Mount"
)

type MountMethod uint32

// mount
const (
	MOUNTPROC3_NULL MountMethod = iota
	MOUNTPROC3_MNT
	MOUNTPROC3_DUMP
	MOUNTPROC3_UMNT
	MOUNTPROC3_UMNTALL
	MOUNTPROC3_EXPORT
)

type MountState uint32

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

func (ms MountState) Error() error {
	errMsg := ""
	switch ms {
	case MNT3_OK:
		return nil
	case MNT3ERR_PERM:
		errMsg = "Not owner"
	case MNT3ERR_NOENT:
		errMsg = "No such file or directory"
	case MNT3ERR_IO:
		errMsg = "I/O error"
	case MNT3ERR_ACCES:
		errMsg = "Permission denied"
	case MNT3ERR_NOTDIR:
		errMsg = "Not a directory"
	case MNT3ERR_INVAL:
		errMsg = "Invalid argument"
	case MNT3ERR_NAMETOOLONG:
		errMsg = "Filename too long"
	case MNT3ERR_NOTSUPP:
		errMsg = "Operation not supported"
	case MNT3ERR_SERVERFAULT:
		errMsg = "A failure on the server"
	default:
		errMsg = "Unknown error"
	}
	return fmt.Errorf("Mount [%d] %s", ms, errMsg)
}

const (
	// The MOUNT service uses AUTH_NONE in the NULL procedure.
	AUTH_WAY   = AuthNone
	MNTPATHLEN = 1024 // Maximum bytes in a path name
	MNTNAMLEN  = 255  // Maximum bytes in a name
	FHSIZE3    = 64   // Maximum bytes in a V3 file handle
)

/*

   typedef opaque fhandle3<FHSIZE3>;
   typedef string dirpath<MNTPATHLEN>;
   typedef string name<MNTNAMLEN>;

*/

type MountRes3 struct {
	Stat      MountState  `xdr:"union"`
	MountInfo MountRes3OK `xdr:"unioncase=0"`
}

type MountRes3OK struct {
	FHandle []byte
	flavors AuthFlavor
}

func mountInit() {
	methods := []string{
		"Null",
		"Mnt",
		"Dump",
		"UMnt",
		"UMntAll",
		"Export",
	}

	producerId := ProcedureID{
		ProgramNumber:  MOUNT_PROG,
		ProgramVersion: MOUNT_VERS,
	}
	for id, procName := range methods {
		producerId.ProcedureNumber = uint32(id)
		if err := RegisterProcedure(
			Procedure{
				ID:   producerId,
				Name: Mount_Name + "." + procName,
			},
		); err != nil {
			panic(err)
		}
	}

	log.Println("\tMount Register over")
}

func Mnt(client *rpc.Client, dirpath string) (*MountRes3, error) {
	if client == nil {
		return nil, errors.New("Could not create pmap client")
	}
	//defer client.Close()
	var result MountRes3
	err := client.Call("Mount.Mnt", dirpath, &result)
	if err != nil {
		return nil, err
	}
	return &result, result.Stat.Error()
}

// MOUNTPROC3_DUMP return
type MountBody struct {
	Name    string // max MNTNAMLEN
	Dirpath string // max MNTPATHLEN
}

type DumpReply struct {
	M    MountBody
	Next *DumpReply `xdr:"optional"`
}

type getDumpReply struct {
	Next *DumpReply `xdr:"optional"`
}

func Dump(client *rpc.Client) ([]MountBody, error) {
	var MountList []MountBody
	var result getDumpReply

	if client == nil {
		return nil, errors.New("Could not create pmap client")
	}
	//defer client.Close()

	err := client.Call("Mount.Dump", nil, &result)
	if err != nil {
		return nil, err
	}

	if result.Next != nil {
		trav := result.Next
		for {
			entry := MountBody(trav.M)
			MountList = append(MountList, entry)
			trav = trav.Next
			if trav == nil {
				break
			}
		}
	}
	return MountList, nil
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
	//defer client.Close()
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
