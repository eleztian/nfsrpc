package nfs_progra

import "fmt"

func (ns NFS3Stat) Error() error {
	errMsg := ""
	switch ns {
	case NFS3_OK:
		return nil
	case NFS3ERR_PERM:
		errMsg = "Not owner"
	case NFS3ERR_NOENT:
		errMsg = "No such file or directory"
	case NFS3ERR_IO:
		errMsg = "I/O error"
	case NFS3ERR_NXIO:
		errMsg = "I/O error"
	case NFS3ERR_ACCES:
		errMsg = "Permission denied"
	case NFS3ERR_EXIST:
		errMsg = "File exists"
	case NFS3ERR_XDEV:
		errMsg = "Attempt to do a cross-device hard link"
	case NFS3ERR_NODEV:
		errMsg = "No such device"
	case NFS3ERR_NOTDIR:
		errMsg = "Not a directory"
	case NFS3ERR_ISDIR:
		errMsg = "Is a directory"
	case NFS3ERR_INVAL:
		errMsg = "Invalid argument or unsupported argument for an operation"
	case NFS3ERR_FBIG:
		errMsg = "File too large"
	case NFS3ERR_NOSPC:
		errMsg = "No space left on device"
	case NFS3ERR_ROFS:
		errMsg = "Read-only file system"
	case NFS3ERR_MLINK:
		errMsg = "Too many hard links"
	case NFS3ERR_NAMETOOLONG:
		errMsg = "The filename in an operation was too long"
	case NFS3ERR_NOTEMPTY:
		errMsg = "An attempt was made to remove a directory that was not empty"
	case NFS3ERR_DQUOT:
		errMsg = "Resource (quota) hard limit exceeded"
	case NFS3ERR_STALE:
		errMsg = "Invalid file handle"
	case NFS3ERR_REMOTE:
		errMsg = "Too many levels of remote in path"
	case NFS3ERR_BADHANDLE:
		errMsg = "Illegal NFS file handle"
	case NFS3ERR_NOT_SYNC:
		errMsg = "Update synchronization mismatch was detected during a SETATTR operation"
	case NFS3ERR_BAD_COOKIE:
		errMsg = "READDIR or READDIRPLUS cookie is stale"
	case NFS3ERR_NOTSUPP:
		errMsg = "Operation is not supported"
	case NFS3ERR_TOOSMALL:
		errMsg = "Buffer or request is too small"
	case NFS3ERR_SERVERFAULT:
		errMsg = "An error occurred on the server which does not map to any of the legal NFS version 3 protocol error values"
	case NFS3ERR_BADTYPE:
		errMsg = "An attempt was made to create an object of a type not supported by the server"
	case NFS3ERR_JUKEBOX:
		errMsg = "he server initiated the request, but was not able to complete it in a timely fashion"
	default:
		errMsg = "unknown error"

	}
	return fmt.Errorf("NFS [%d] %s", ns, errMsg)
}
