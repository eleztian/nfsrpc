package methods

import "bs-2018/mynfs/nfsrpc"

// MOUNT
// RFC 1813 Section 5.0

const (
	MOUNT_PROG = 100005
	MOUNT_VERS = 3
)

type MountMethod int

// mount
const (
	MOUNTPROC3_NULL MountMethod = iota
	MOUNTPROC3_MNT
	MOUNTPROC3_DUMP
	MOUNTPROC3_UMNT
	MOUNTPROC3_UMNTALL
	MOUNTPROC3_EXPORT
)

// mount error
const (
	MNT3_OK             = 0     /* no error */
	MNT3ERR_PERM        = 1     /* Not owner */
	MNT3ERR_NOENT       = 2     /* No such file or directory */
	MNT3ERR_IO          = 5     /* I/O error */
	MNT3ERR_ACCES       = 13    /* Permission denied */
	MNT3ERR_NOTDIR      = 20    /* Not a directory */
	MNT3ERR_INVAL       = 22    /* Invalid argument */
	MNT3ERR_NAMETOOLONG = 63    /* Filename too long */
	MNT3ERR_NOTSUPP     = 10004 /* Operation not supported */
	MNT3ERR_SERVERFAULT = 10006 /* A failure on the server */
)

const (
	AUTH_WAY   = nfsrpc.AuthNone
	MNTPATHLEN = 1024 /* Maximum bytes in a path name */
	MNTNAMLEN  = 255  /* Maximum bytes in a name */
	FHSIZE3    = 64   /* Maximum bytes in a V3 file handle */
)
