package nfs_progra

type SattrGuard3 struct {
	Check    bool    `xdr:"union"`
	ObjCTime NFSTime `xdr:"unioncase=1"`
}

type SetAttr3Res struct {
	Status NFS3Stat
	ObjWcc WccData
}

type SetAttr3Aargs struct {
	Object  NFS_FH3
	NewAttr SAttr3
	Guard   SattrGuard3
}
