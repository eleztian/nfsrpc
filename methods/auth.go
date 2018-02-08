package methods

import "bs-2018/mynfs/nfsrpc"

func NewAuth(authWay nfsrpc.AuthFlavor, body interface{}) *nfsrpc.OpaqueAuth {
	switch authWay {
	case nfsrpc.AuthDh:
	case nfsrpc.AuthKerb:
	case nfsrpc.AuthNone:
	case nfsrpc.AuthRSA:
	case nfsrpc.AuthSys:
	case nfsrpc.AuthShort:
	case nfsrpc.RPCsecGss:
	}
	return &nfsrpc.OpaqueAuth{Flavor: authWay}
}

type AUTH_UNIX struct {
	Stamp       uint32
	Machinename string
	Uid         uint32
	Gid         uint32
	Gids        uint32
}
