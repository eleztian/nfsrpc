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
	case nfsrpc.RPCSecGss:
	}
	return &nfsrpc.OpaqueAuth{Flavor: authWay}
}

