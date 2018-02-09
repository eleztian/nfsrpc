package nfs

type SattrGuard3 struct {
	Check    bool    `xdr:"union"`
	ObjCTime NFSTime `xdr:"unioncase=true"`
}

type SetAttr3Aargs struct {
	Object  NFS_FH3
	NewAttr SAttr3
	Guard   SattrGuard3
}

type SetAttr3Resok struct {
	ObjWcc WccData
}

type SetAttr3ResFail struct {
	ObjWcc WccData
}

type SetAttr3Res struct {
	Status  NFS3Stat        `xdr:"union"`
	ResOk   SetAttr3Resok
	//ResFail SetAttr3ResFail `xdr:"unioncase!=0"`
}