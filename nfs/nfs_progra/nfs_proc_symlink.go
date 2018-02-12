package nfs_progra

/*
Procedure 10: SYMLINK - Create a symbolic link

   SYNOPSIS

      SYMLINK3res NFSPROC3_SYMLINK(SYMLINK3args) = 10;

      struct symlinkdata3 {
           sattr3    symlink_attributes;
           nfspath3  symlink_data;
      };

      struct SYMLINK3args {
           diropargs3    where;
           symlinkdata3  symlink;
      };

      struct SYMLINK3resok {
           post_op_fh3   obj;
           post_op_attr  obj_attributes;
           wcc_data      dir_wcc;
      };

      struct SYMLINK3resfail {
           wcc_data      dir_wcc;
      };

      union SYMLINK3res switch (nfsstat3 status) {
      case NFS3_OK:
           SYMLINK3resok   resok;
      default:
           SYMLINK3resfail resfail;
      };

   DESCRIPTION

      Procedure SYMLINK creates a new symbolic link. On entry,
      the arguments in SYMLINK3args are:

      where
         The location of the symbolic link to be created:

         dir
            The file handle for the directory in which the
            symbolic link is to be created.

         name
            The name that is to be associated with the created
            symbolic link. Refer to General comments on
            filenames on page 30.

      symlink
         The symbolic link to create:

         symlink_attributes
            The initial attributes for the symbolic link.

         symlink_data
            The string containing the symbolic link data.

      On successful return, SYMLINK3res.status is NFS3_OK and
      SYMLINK3res.resok contains:

      obj
         The file handle for the newly created symbolic link.

      obj_attributes
         The attributes for the newly created symbolic link.

      dir_wcc
         Weak cache consistency data for the directory,
         where.dir. For a client that requires only the
         post-SYMLINK directory attributes, these can be found
         in dir_wcc.after.

      Otherwise, SYMLINK3res.status contains the error on
      failure and SYMLINK3res.resfail contains the following:

      dir_wcc
         Weak cache consistency data for the directory,
         where.dir. For a client that requires only the
         post-SYMLINK directory attributes, these can be found
         in dir_wcc.after. Even though the SYMLINK failed, full
         wcc_data is returned to allow the client to determine
         whether the failing SYMLINK changed the directory.

   IMPLEMENTATION

      Refer to General comments on filenames on page 30.

      For symbolic links, the actual file system node and its
      contents are expected to be created in a single atomic
      operation.  That is, once the symbolic link is visible,
      there must not be a window where a READLINK would fail or
      return incorrect data.

   ERRORS

      NFS3ERR_IO
      NFS3ERR_ACCES
      NFS3ERR_EXIST
      NFS3ERR_NOTDIR
      NFS3ERR_NOSPC
      NFS3ERR_ROFS
      NFS3ERR_NAMETOOLONG
      NFS3ERR_DQUOT
      NFS3ERR_STALE
      NFS3ERR_BADHANDLE
      NFS3ERR_NOTSUPP
      NFS3ERR_SERVERFAULT
*/

type SymLinkData3 struct {
	SymLinkAttrs SAttr3
	SymLinkData  string
}

type symLinkResOk struct {
	Obj      PostOpFh3
	ObjAttrs PostOpAttr
}

type SymLinkRes struct {
	Status NFS3Stat     `xdr:"union"`
	ResOk  symLinkResOk `xdr:"unioncase=0"`
	DirWcc WccData
}

type SymLinkArgs struct {
	Where   DirOpArgs3
	SymLink SymLinkData3
}
