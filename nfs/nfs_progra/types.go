package nfs_progra

const (
	NFS_PROG = 100003
	NFS_VERS = 3
)

type NFSMethod uint32

const (
	NFSPROC3_NULL    NFSMethod = iota // void (void)
	NFSPROC3_GETATTR                  // GETATTR3res (SETATTR3args)
	NFSPROC3_SETATTR
	NFSPROC3_LOOKUP // LOOKUP3res (LOOKUP3args)
	NFSPROC3_ACCESS // ACCESS3res (ACCESS3args)
	NFSPROC3_READLINK
	NFSPROC3_READ
	NFSPROC3_WRITE
	NFSPROC3_CREATE
	NFSPROC3_MKDIR
	NFSPROC3_SYMLINK
	NFSPROC3_MKNOD
	NFSPROC3_REMOVE
	NFSPROC3_RMDIR
	NFSPROC3_RENAME
	NFSPROC3_LINK
	NFSPROC3_READDIR
	NFSPROC3_READDIRPLUS
	NFSPROC3_FSSTAT
	NFSPROC3_FSINFO
	NFSPROC3_PATHCONF
	NFSPROC3_COMMIT // 21
)

type NFS3Stat uint32

const (
	NFS3_OK             NFS3Stat = 0
	NFS3ERR_PERM                 = 1
	NFS3ERR_NOENT                = 2
	NFS3ERR_IO                   = 5
	NFS3ERR_NXIO                 = 6
	NFS3ERR_ACCES                = 13
	NFS3ERR_EXIST                = 17
	NFS3ERR_XDEV                 = 18
	NFS3ERR_NODEV                = 19
	NFS3ERR_NOTDIR               = 20
	NFS3ERR_ISDIR                = 21
	NFS3ERR_INVAL                = 22
	NFS3ERR_FBIG                 = 27
	NFS3ERR_NOSPC                = 28
	NFS3ERR_ROFS                 = 30
	NFS3ERR_MLINK                = 31
	NFS3ERR_NAMETOOLONG          = 63
	NFS3ERR_NOTEMPTY             = 66
	NFS3ERR_DQUOT                = 69
	NFS3ERR_STALE                = 70
	NFS3ERR_REMOTE               = 71
	NFS3ERR_BADHANDLE            = 10001
	NFS3ERR_NOT_SYNC             = 10002
	NFS3ERR_BAD_COOKIE           = 10003
	NFS3ERR_NOTSUPP              = 10004
	NFS3ERR_TOOSMALL             = 10005
	NFS3ERR_SERVERFAULT          = 10006
	NFS3ERR_BADTYPE              = 10007
	NFS3ERR_JUKEBOX              = 10008
)

type Opaque []byte

// The nfs_fh3 is the variable-length opaque object returned by the
// server on LOOKUP, CREATE, SYMLINK, MKNOD, LINK, or READDIRPLUS
// operations, which is used by the client on subsequent operations
// to reference the file.
type NFS_FH3 struct {
	Data Opaque // max NFS3_FHSIZE
}

// gives the type of a file.
type FType uint32

const (
	NF3REG FType = iota + 1
	NF3DIR
	NF3BLK
	NF3CHR
	NF3LNK
	NF3SOCK
	NF3FIFO
)

// The interpretation of the two words depends on the type of file
// system object.
type SpecData struct {
	SpecData1 uint32
	SpecData2 uint32
}

// The nfstime3 structure gives the number of seconds and
// nanoseconds since midnight January 1, 1970 Greenwich Mean Time.
type NFSTime struct {
	Seconds  uint32
	NSeconds uint32
}

// This structure defines the attributes of a file system object.
type Fattr3 struct {
	Type    FType
	Mode    uint32
	NLink   uint32
	Uid     uint32
	Gid     uint32
	Size    uint64
	Used    uint64
	Rdev    SpecData
	FsId    uint64
	FiledId uint64
	ATime   NFSTime
	MTime   NFSTime
	CTime   NFSTime
}

// This structure is used for returning attributes in those
// operations that are not directly involved with manipulating
// attributes.
type PostOpAttr struct {
	AttributesFollow bool   `xdr:"union"`
	Attributes       Fattr3 `xdr:"unioncase=1"`
}

// This is the subset of pre-operation attributes needed to better
// support the weak cache consistency semantics.
type WccAttr struct {
	// Size is the file size in bytes of the object before the operation.
	Size  uint64
	MTime NFSTime
	CTime NFSTime
}
type PreOpAttr struct {
	AttributesFollow bool    `xdr:"union"`
	Attributes       WccAttr `xdr:"unioncase=1"`
}

//When a client performs an operation that modifies the state of a
//file or directory on the server, it cannot immediately determine
//from the post-operation attributes whether the operation just
//performed was the only operation on the object since the last
//time the client received the attributes for the object. This is
//important, since if an intervening operation has changed the
//object, the client will need to invalidate any cached data for
//the object (except for the data that it just wrote).
// 当客户机执行修改服务器上文件或目录的状态的操作时，
// 它不能立即从操作后的属性中判断，是否执行的操作是自上次客户端接收到对象的属性
// 以来唯一的操作。这很重要，因为如果一个中间操作改变了对象，
// 那么客户端将需要为对象的缓存数据无效(除了它刚刚写的数据之外)。
type WccData struct {
	Before PreOpAttr
	After  PostOpAttr
}

// This is the structure used to return a file handle from the
// CREATE, MKDIR, SYMLINK, MKNOD, and READDIRPLUS requests.
type PostOpFh3 struct {
	HandleFollows bool    `xdr:"union"`
	Handle        NFS_FH3 `xdr:"unioncase=1"`
}

type TimeHow uint32

const (
	DONT_CHANGE TimeHow = iota
	SET_TO_SERVER_TIME
	SET_TO_CLIENT_TIME
)

type SetMode struct {
	SetIt bool   `xdr:"union"`
	Mode  uint32 `xdr:"unioncase=1"`
}
type SetUid struct {
	SetIt bool   `xdr:"union"`
	Uid   uint32 `xdr:"unioncase=1"`
}
type SetGid struct {
	SetIt bool   `xdr:"union"`
	Gid   uint32 `xdr:"unioncase=1"`
}
type SetSize struct {
	SetIt uint   `xdr:"union"`
	Size  uint64 `xdr:"unioncase=0"`
}
type SetATime struct {
	SetIt TimeHow `xdr:"union"`
	ATime uint64  `xdr:"unioncase=2"`
}
type SetMTime struct {
	SetIt TimeHow `xdr:"union"`
	MTime uint64  `xdr:"unioncase=2"`
}

// The sattr3 structure contains the file attributes that can be
// set from the client.
type SAttr3 struct {
	Mode  SetMode
	Uid   SetUid
	Gid   SetGid
	Size  SetSize
	ATime SetATime
	MTime SetMTime
}

// The diropargs3 structure is used in directory operations.
type DirOpArgs3 struct {
	Dir  NFS_FH3
	Name string
}

const NFS3_WRITEVERFSIZE = 8

type WriteVerf3 [NFS3_WRITEVERFSIZE]byte

const NFS3_CREADTEVERFSIZE = 8

type CreateVerf3 [NFS3_CREADTEVERFSIZE]byte

const NFS3_COOKIEVERFSIZE = 8

type CookieVerf3 [NFS3_COOKIEVERFSIZE]byte
