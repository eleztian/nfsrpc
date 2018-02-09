package main

import (
	"bs-2018/mynfs/methods"
	"bs-2018/mynfs/nfsrpc"
	"fmt"
	"log"
	"net"
)

var host = "111.231.215.178"

func main() {
	nfsrpc.RegisterProcedure(nfsrpc.Procedure{
		nfsrpc.ProcedureID{
			methods.MOUNT_PROG,
			methods.MOUNT_VERS,
			uint32(methods.MOUNTPROC3_EXPORT),
		},
		"Mount.Export",
	})
	nfsrpc.RegisterProcedure(nfsrpc.Procedure{
		nfsrpc.ProcedureID{
			methods.MOUNT_PROG,
			methods.MOUNT_VERS,
			uint32(methods.MOUNTPROC3_DUMP),
		},
		"Mount.Dump",
	})
	nfsrpc.RegisterProcedure(nfsrpc.Procedure{
		nfsrpc.ProcedureID{
			methods.MOUNT_PROG,
			methods.MOUNT_VERS,
			uint32(methods.MOUNTPROC3_MNT),
		},
		"Mount.Mnt",
	})
	conn, err := net.Dial("tcp", host+":1011")
	if err != nil {
		fmt.Println(err)
	}

	client := nfsrpc.NewClient(conn,nil, nil)
	fmt.Println(methods.Export(client))
	fmt.Println(methods.Dump(client))
	fmt.Println(methods.Mnt(client, "/home/ubuntu/go"))
	maps, err := nfsrpc.PmapGetMaps(host)
	port, err := nfsrpc.PmapGetPort(host, 100005, 3, 6)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(port)
	if err != nil {
		log.Fatal("sunrpc.PmapGetMaps() failed: " + err.Error())
	}

	protocols := make(map[uint32]string, 2)
	protocols[uint32(6)] = "tcp"
	protocols[uint32(17)] = "udp"

	fmt.Printf("\tprogram\tvers\tproto\tport\t\n")
	for _, m := range maps {
		fmt.Printf("\t%d\t%d\t%s\t%d\n", m.Program, m.Version, protocols[m.Protocol], m.Port)
	}

}
