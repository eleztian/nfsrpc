package nfs

/*
Procedure 20: PATHCONF - Retrieve POSIX information

   SYNOPSIS

      PATHCONF3res NFSPROC3_PATHCONF(PATHCONF3args) = 20;

      struct PATHCONF3args {
           nfs_fh3   object;
      };

      struct PATHCONF3resok {
           post_op_attr obj_attributes;
           uint32       linkmax;
           uint32       name_max;
           bool         no_trunc;
           bool         chown_restricted;
           bool         case_insensitive;
           bool         case_preserving;
      };

      struct PATHCONF3resfail {
           post_op_attr obj_attributes;
      };

      union PATHCONF3res switch (nfsstat3 status) {
      case NFS3_OK:
           PATHCONF3resok   resok;
      default:
           PATHCONF3resfail resfail;
      };

   DESCRIPTION

      Procedure PATHCONF retrieves the pathconf information for
      a file or directory. If the FSF_HOMOGENEOUS bit is set in
      FSFINFO3resok.properties, the pathconf information will be
      the same for all files and directories in the exported
      file system in which this file or directory resides. On
      entry, the arguments in PATHCONF3args are:

      object
         The file handle for the file system object.

      On successful return, PATHCONF3res.status is NFS3_OK and
      PATHCONF3res.resok contains:

      obj_attributes
         The attributes of the object specified by object.

      linkmax
         The maximum number of hard links to an object.

      name_max
         The maximum length of a component of a filename.

      no_trunc
         If TRUE, the server will reject any request that
         includes a name longer than name_max with the error,
         NFS3ERR_NAMETOOLONG. If FALSE, any length name over
         name_max bytes will be silently truncated to name_max
         bytes.

      chown_restricted
         If TRUE, the server will reject any request to change
         either the owner or the group associated with a file if
         the caller is not the privileged user. (Uid 0.)

      case_insensitive
         If TRUE, the server file system does not distinguish
         case when interpreting filenames.

      case_preserving
         If TRUE, the server file system will preserve the case
         of a name during a CREATE, MKDIR, MKNOD, SYMLINK,
         RENAME, or LINK operation.

      Otherwise, PATHCONF3res.status contains the error on
      failure and PATHCONF3res.resfail contains the following:

      obj_attributes
         The attributes of the object specified by object.

   IMPLEMENTATION

      In some implementations of the NFS version 2 protocol,
      pathconf information was obtained at mount time through
      the MOUNT protocol.  The proper place to obtain it, is as
      here, in the NFS version 3 protocol itself.

   ERRORS

      NFS3ERR_STALE
      NFS3ERR_BADHANDLE
      NFS3ERR_SERVERFAULT
*/

type PathConf3Args struct {
	Object NFS_FH3
}

type PathConf3ResOk struct {
	LinkMax         uint32
	NameMax         uint32
	NoTrunc         bool
	ChownRestricted bool
	CaseInsensitive bool
	CasePreserving  bool
}

type PathConfRes struct {
	Status   NFS3Stat `xdr:"union"`
	ObjAttrs PostOpAttr
	ResOk    PathConf3ResOk `xdr:"unioncase=0"`
}
