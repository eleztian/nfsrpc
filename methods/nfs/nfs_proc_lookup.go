package nfs

/*
 SYNOPSIS

      LOOKUP3res NFSPROC3_LOOKUP(LOOKUP3args) = 3;

      struct LOOKUP3args {
           diropargs3  what;
      };

      struct LOOKUP3resok {
           nfs_fh3      object;
           post_op_attr obj_attributes;
           post_op_attr dir_attributes;
      };

      struct LOOKUP3resfail {
           post_op_attr dir_attributes;
      };

      union LOOKUP3res switch (nfsstat3 status) {
      case NFS3_OK:
           LOOKUP3resok    resok;
      default:
           LOOKUP3resfail  resfail;
      };

   DESCRIPTION

      Procedure LOOKUP searches a directory for a specific name
      and returns the file handle for the corresponding file
      system object. On entry, the arguments in LOOKUP3args
      are:

      what
         Object to look up:

         dir
            The file handle for the directory to search.

         name
            The filename to be searched for. Refer to General
            comments on filenames on page 30.

      On successful return, LOOKUP3res.status is NFS3_OK and
      LOOKUP3res.resok contains:

      object
         The file handle of the object corresponding to
         what.name.

      obj_attributes
         The attributes of the object corresponding to
         what.name.

      dir_attributes
         The post-operation attributes of the directory,
         what.dir.

      Otherwise, LOOKUP3res.status contains the error on failure and
      LOOKUP3res.resfail contains the following:

      dir_attributes
         The post-operation attributes for the directory,
         what.dir.

   IMPLEMENTATION

      At first glance, in the case where what.name refers to a
      mount point on the server, two different replies seem
      possible. The server can return either the file handle for
      the underlying directory that is mounted on or the file
      handle of the root of the mounted directory.  This
      ambiguity is simply resolved. A server will not allow a
      LOOKUP operation to cross a mountpoint to the root of a
      different filesystem, even if the filesystem is exported.
      This does not prevent a client from accessing a hierarchy
      of filesystems exported by a server, but the client must
      mount each of the filesystems individually so that the
      mountpoint crossing takes place on the client.  A given
      server implementation may refine these rules given
      capabilities or limitations particular to that
      implementation. Refer to [X/OpenNFS] for a discussion on
      exporting file systems.

      Two filenames are distinguished, as in the NFS version 2
      protocol.  The name, ".", is an alias for the current
      directory and the name, "..", is an alias for the parent
      directory; that is, the directory that includes the
      specified directory as a member. There is no facility for
      dealing with a multiparented directory and the NFS
      protocol assumes a hierarchical organization, organized as
      a single-rooted tree.

      Note that this procedure does not follow symbolic links.
      The client is responsible for all parsing of filenames
      including filenames that are modified by symbolic links
      encountered during the lookup process.

   ERRORS

      NFS3ERR_IO
      NFS3ERR_NOENT
      NFS3ERR_ACCES
      NFS3ERR_NOTDIR
      NFS3ERR_NAMETOOLONG
      NFS3ERR_STALE
      NFS3ERR_BADHANDLE
      NFS3ERR_SERVERFAULT
*/

type lookUp3Resok struct {
	Object   NFS_FH3
	ObjAttrs PostOpAttr
}

type LookUp3Args struct {
	What DirOpArgs3
}

type LookUp3Res struct {
	Status   NFS3Stat     `xdr:"union"`
	ResOk    lookUp3Resok `xdr:"unioncase=0"`
	DirAttrs PostOpAttr
}
