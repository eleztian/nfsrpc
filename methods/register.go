package methods

import "bs-2018/mynfs/nfsrpc"

func init() {
	nfsrpc.RegisterProcedure(nfsrpc.Procedure{
		nfsrpc.ProcedureID{
			MOUNT_PROG,
			MOUNT_VERS,
			uint32(MOUNTPROC3_EXPORT),
		},
		"Mount.Export",
	})
	nfsrpc.RegisterProcedure(nfsrpc.Procedure{
		nfsrpc.ProcedureID{
			MOUNT_PROG,
			MOUNT_VERS,
			uint32(MOUNTPROC3_DUMP),
		},
		"Mount.Dump",
	})
	nfsrpc.RegisterProcedure(nfsrpc.Procedure{
		nfsrpc.ProcedureID{
			MOUNT_PROG,
			MOUNT_VERS,
			uint32(MOUNTPROC3_MNT),
		},
		"Mount.Mnt",
	})
}
