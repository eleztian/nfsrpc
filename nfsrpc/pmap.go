package nfsrpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"strings"
)

const (
	pmapPort                 = 111
	portmapperProgramNumber  = 100000
	portmapperProgramVersion = 2
	PrograNmae               = "Pmap"
)

// Protocol is a type representing the protocol (TCP or UDP) over which the
// program/server being registered listens on.
type Protocol uint32

const (
	// IPProtoTCP is the protocol number for TCP/IP
	IPProtoTCP Protocol = 6
	// IPProtoUDP is the protocol number for UDP/IP
	IPProtoUDP Protocol = 17
)

var defaultAddress = "127.0.0.1:" + strconv.Itoa(pmapPort)

// PortMapping is a mapping between (program, version, protocol) to port number
type PortMapper struct {
	Program  uint32
	Version  uint32
	Protocol uint32
	Port     uint32
}

func PamapInit() {
	// This is ordered as per procedure number
	methods := []string{
		"ProcNull",
		"ProcSet",
		"ProcUnset",
		"ProcGetPort",
		"ProcDump",
		"ProcCallIt",
	}

	producerId := ProcedureID{
		ProgramNumber:  portmapperProgramNumber,
		ProgramVersion: portmapperProgramVersion,
	}

	for id, procName := range methods {
		producerId.ProcedureNumber = uint32(id)
		if err := RegisterProcedure(
			Procedure{
				ID:   producerId,
				Name: PrograNmae + "." + procName,
			},
		); err != nil {
			panic(err)
		}
	}

	log.Println("\tPmap Register over")
}

func NewPortMapperClient(host string) *rpc.Client {
	if host == "" {
		host = defaultAddress
	} else if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:%d", host, pmapPort)
	}
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil
	}
	return rpc.NewClientWithCodec(NewClientCodec(conn, nil, nil, nil))
}

type portMappingList struct {
	Map  PortMapper
	Next *portMappingList `xdr:"optional"`
}

type getMapsReply struct {
	Next *portMappingList `xdr:"optional"`
}

// PmapGetPort returns the port number on which the program specified is
// awaiting call requests. If host is empty string, localhost is used.
func PmapGetPort(host string, programNumber, programVersion uint32, protocol Protocol) (uint32, error) {

	var port uint32

	client := NewPortMapperClient(host)
	if client == nil {
		return port, errors.New("Could not create pmap client")
	}
	defer client.Close()

	mapping := &PortMapper{
		Program:  programNumber,
		Version:  programVersion,
		Protocol: uint32(protocol),
	}

	err := client.Call("Pmap.ProcGetPort", mapping, &port)
	return port, err
}

// PmapGetMaps returns a list of PortMapping entries present in portmapper's
// database. If host is empty string, localhost is used.
func PmapGetMaps(host string) ([]PortMapper, error) {

	var mappings []PortMapper
	var result getMapsReply

	client := NewPortMapperClient(host)
	if client == nil {
		return nil, errors.New("Could not create pmap client")
	}
	defer client.Close()

	err := client.Call("Pmap.ProcDump", nil, &result)
	if err != nil {
		return nil, err
	}

	if result.Next != nil {
		trav := result.Next
		for {
			entry := PortMapper(trav.Map)
			mappings = append(mappings, entry)
			trav = trav.Next
			if trav == nil {
				break
			}
		}
	}

	return mappings, nil
}

// PmapSet creates port mapping of the program specified. It return true on
// success and false otherwise.
func PmapSet(host string, programNumber, programVersion uint32, protocol Protocol, port uint32) (bool, error) {

	var result bool

	client := NewPortMapperClient(host)
	if client == nil {
		return false, errors.New("Could not create pmap client")
	}
	defer client.Close()

	mapping := &PortMapper{
		Program:  programNumber,
		Version:  programVersion,
		Protocol: uint32(protocol),
		Port:     port,
	}

	err := client.Call("Pmap.ProcSet", mapping, &result)
	return result, err
}

// PmapUnset will unregister the program specified. It returns true on success
// and false otherwise.
func PmapUnset(host string, programNumber, programVersion uint32) (bool, error) {

	var result bool

	client := NewPortMapperClient(host)
	if client == nil {
		return false, errors.New("Could not create pmap client")
	}
	defer client.Close()

	mapping := &PortMapper{
		Program: programNumber,
		Version: programVersion,
	}

	err := client.Call("Pmap.ProcUnset", mapping, &result)
	return result, err
}
