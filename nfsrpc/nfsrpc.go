package nfsrpc

import (
	"log"
	"bytes"
	"fmt"
	"github.com/rasky/go-xdr/xdr2"
)

func init()  {
	PamapInit()
	MountInit()
	log.Printf("RPC init finished.")
}

func Auth(authWay AuthFlavor, authBody interface{}) (*OpaqueAuth, error) {
	buffer := new(bytes.Buffer)
	// xdr格式编码
	if _, err := xdr.Marshal(buffer, authBody); err != nil {
		return nil, err
	}
	bs := buffer.Bytes()
	if len(bs) > OpaueBodyMaxLength {
		return nil, fmt.Errorf("the auth body too long (%d), please not over %d ",
			len(bs), OpaueBodyMaxLength)
	}
	return &OpaqueAuth{authWay, bs}, nil
}