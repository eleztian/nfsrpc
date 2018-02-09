package nfs

type LookUp3Args struct {
	What DirOpArgs3
}

type LookUp3Resok struct {
	Object   NFS_FH3
	ObjAttrs PostOpAttr
	DirAttrs PostOpAttr
}

type LookUp3ResFail struct {
	DirAttrs PostOpAttr
}

type LookUp3Res struct {
	Status  NFS3Stat       `xdr:"union"`
	ResOk   LookUp3Resok   `xdr:"unioncase=0"`
	ResFail LookUp3ResFail `xdr:"unioncase=1"`
}
