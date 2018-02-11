package nfs

/*
Procedure 9: MKDIR - Create a directory

   SYNOPSIS

      MKDIR3res NFSPROC3_MKDIR(MKDIR3args) = 9;

      struct MKDIR3args {
           diropargs3   where;
           sattr3       attributes;
      };

      struct MKDIR3resok {
           post_op_fh3   obj;
           post_op_attr  obj_attributes;
           wcc_data      dir_wcc;
      };

      struct MKDIR3resfail {
           wcc_data      dir_wcc;
      };

      union MKDIR3res switch (nfsstat3 status) {
      case NFS3_OK:
           MKDIR3resok   resok;
      default:
           MKDIR3resfail resfail;
      };

   DESCRIPTION

      Procedure MKDIR creates a new subdirectory. On entry, the
      arguments in MKDIR3args are:

      where
         The location of the subdirectory to be created:

         dir
            The file handle for the directory in which the
            subdirectory is to be created.

         name
            The name that is to be associated with the created
            subdirectory. Refer to General comments on filenames
            on page 30.

      attributes
         The initial attributes for the subdirectory.

      On successful return, MKDIR3res.status is NFS3_OK and the
      results in MKDIR3res.resok are:

      obj
         The file handle for the newly created directory.

      obj_attributes
         The attributes for the newly created subdirectory.

      dir_wcc
         Weak cache consistency data for the directory,
         where.dir. For a client that requires only the
         post-MKDIR directory attributes, these can be found in
         dir_wcc.after.

      Otherwise, MKDIR3res.status contains the error on failure
      and MKDIR3res.resfail contains the following:

      dir_wcc
         Weak cache consistency data for the directory,
         where.dir. For a client that requires only the
         post-MKDIR directory attributes, these can be found in
         dir_wcc.after. Even though the MKDIR failed, full
         wcc_data is returned to allow the client to determine
         whether the failing MKDIR resulted in any change to the
         directory.

   IMPLEMENTATION

      Many server implementations will not allow the filenames,
      "." or "..", to be used as targets in a MKDIR operation.
      In this case, the server should return NFS3ERR_EXIST.
      Refer to General comments on filenames on page 30.

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

type MkDir3ResOk struct {
	Obj      PostOpFh3
	ObjAttrs PostOpAttr
}

type MkDirRes struct {
	Status      NFS3Stat `xdr:"union"`
	MkDir3ResOk `xdr:"unioncase=0"`
	DirWcc      WccData
}

type MkDirArgs struct {
	Where DirOpArgs3
	Attrs SAttr3
}
