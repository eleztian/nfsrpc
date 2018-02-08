package main

import (
	"fmt"
	"log"
	"bs-2018/mynfs/nfsrpc"
)

var host = "111.231.215.178"

func main() {

	maps, err := nfsrpc.PmapGetMaps(host)
	port , err := nfsrpc.PmapGetPort(host, 100005, 3, 6)
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
