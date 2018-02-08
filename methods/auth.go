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

type AuthsysParms struct {
	Stamp       uint32
	MachineName string // max 255
	Uid         uint32
	Gid         uint32
	Gids        uint32 // max 16 a counted array of groups that contain the caller as a member.
}