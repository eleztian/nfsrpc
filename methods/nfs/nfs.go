package nfs

import (
	"bs-2018/mynfs/nfsrpc"
)

func init() {
	methods := []string{
		"Null",
		"Getattr",
		"Setattr",
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
		ProgramNumber:NFS_PROG,
		ProgramVersion:NFS_VERS,
	}
	for id, procName := range methods {
		producerId.ProcedureNumber = uint32(id)
		if err := nfsrpc.RegisterProcedure(
			nfsrpc.Procedure{
				ID:producerId,
				Name:"NFS." + procName,
			},
		); err != nil {
			panic(err)
		}
	}
}
