package nfs_progra

/*
Procedure 8: CREATE - Create a file

   SYNOPSIS

      CREATE3res NFSPROC3_CREATE(CREATE3args) = 8;

      enum createmode3 {
           UNCHECKED = 0,
           GUARDED   = 1,
           EXCLUSIVE = 2
      };

      union createhow3 switch (createmode3 mode) {
      case UNCHECKED:
      case GUARDED:
           sattr3       obj_attributes;
      case EXCLUSIVE:
           createverf3  verf;
      };

      struct CREATE3args {
           diropargs3   where;
           createhow3   how;
      };

      struct CREATE3resok {
           post_op_fh3   obj;
           post_op_attr  obj_attributes;
           wcc_data      dir_wcc;
      };

      struct CREATE3resfail {
           wcc_data      dir_wcc;
      };

      union CREATE3res switch (nfsstat3 status) {
      case NFS3_OK:
           CREATE3resok    resok;
      default:
           CREATE3resfail  resfail;
      };

   DESCRIPTION

      Procedure CREATE creates a regular file. On entry, the
      arguments in CREATE3args are:

      where
         The location of the file to be created:

         dir
            The file handle for the directory in which the file
            is to be created.

         name
            The name that is to be associated with the created
            file.  Refer to General comments on filenames on
            page 30.

      When creating a regular file, there are three ways to
      create the file as defined by:

      how
         A discriminated union describing how the server is to
         handle the file creation along with the appropriate
         attributes:

      mode
         One of UNCHECKED, GUARDED, and EXCLUSIVE. UNCHECKED
         means that the file should be created without checking
         for the existence of a duplicate file in the same
         directory. In this case, how.obj_attributes is a sattr3
         describing the initial attributes for the file. GUARDED
         specifies that the server should check for the presence
         of a duplicate file before performing the create and
         should fail the request with NFS3ERR_EXIST if a
         duplicate file exists. If the file does not exist, the
         request is performed as described for UNCHECKED.
         EXCLUSIVE specifies that the server is to follow
         exclusive creation semantics, using the verifier to
         ensure exclusive creation of the target. No attributes
         may be provided in this case, since the server may use
         the target file metadata to store the createverf3
         verifier.

      On successful return, CREATE3res.status is NFS3_OK and the
      results in CREATE3res.resok are:

      obj
         The file handle of the newly created regular file.

      obj_attributes
         The attributes of the regular file just created.

      dir_wcc
         Weak cache consistency data for the directory,
         where.dir. For a client that requires on the
         post-CREATE directory attributes, these can be found in
         dir_wcc.after.

      Otherwise, CREATE3res.status contains the error on failure
      and CREATE3res.resfail contains the following:

      dir_wcc
         Weak cache consistency data for the directory,
         where.dir. For a client that requires only the
         post-CREATE directory attributes, these can be found in
         dir_wcc.after. Even though the CREATE failed, full
         wcc_data is returned to allow the client to determine
         whether the failing CREATE resulted in any change to
         the directory.

   IMPLEMENTATION

      Unlike the NFS version 2 protocol, in which certain fields
      in the initial attributes structure were overloaded to
      indicate creation of devices and FIFOs in addition to
      regular files, this procedure only supports the creation
      of regular files. The MKNOD procedure was introduced in
      the NFS version 3 protocol to handle creation of devices
      and FIFOs. Implementations should have no reason in the
      NFS version 3 protocol to overload CREATE semantics.

      One aspect of the NFS version 3 protocol CREATE procedure
      warrants particularly careful consideration: the mechanism
      introduced to support the reliable exclusive creation of
      regular files. The mechanism comes into play when how.mode
      is EXCLUSIVE.  In this case, how.verf contains a verifier
      that can reasonably be expected to be unique.  A
      combination of a client identifier, perhaps the client
      network address, and a unique number generated by the
      client, perhaps the RPC transaction identifier, may be
      appropriate.

      If the file does not exist, the server creates the file
      and stores the verifier in stable storage. For file
      systems that do not provide a mechanism for the storage of
      arbitrary file attributes, the server may use one or more
      elements of the file metadata to store the verifier. The
      verifier must be stored in stable storage to prevent
      erroneous failure on retransmission of the request. It is
      assumed that an exclusive create is being performed
      because exclusive semantics are critical to the
      application. Because of the expected usage, exclusive
      CREATE does not rely solely on the normally volatile
      duplicate request cache for storage of the verifier. The
      duplicate request cache in volatile storage does not
      survive a crash and may actually flush on a long network
      partition, opening failure windows.  In the UNIX local
      file system environment, the expected storage location for
      the verifier on creation is the metadata (time stamps) of
      the file. For this reason, an exclusive file create may
      not include initial attributes because the server would
      have nowhere to store the verifier.

      If the server can not support these exclusive create
      semantics, possibly because of the requirement to commit
      the verifier to stable storage, it should fail the CREATE
      request with the error, NFS3ERR_NOTSUPP.

      During an exclusive CREATE request, if the file already
      exists, the server reconstructs the file's verifier and
      compares it with the verifier in the request. If they
      match, the server treats the request as a success. The
      request is presumed to be a duplicate of an earlier,
      successful request for which the reply was lost and that
      the server duplicate request cache mechanism did not
      detect. If the verifiers do not match, the request is
      rejected with the status, NFS3ERR_EXIST.

      Once the client has performed a successful exclusive
      create, it must issue a SETATTR to set the correct file
      attributes.  Until it does so, it should not rely upon any
      of the file attributes, since the server implementation
      may need to overload file metadata to store the verifier.

      Use of the GUARDED attribute does not provide exactly-once
      semantics.  In particular, if a reply is lost and the
      server does not detect the retransmission of the request,
      the procedure can fail with NFS3ERR_EXIST, even though the
      create was performed successfully.

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

type CreateMode3 int

const (
	UNCHECKED = 0
	GUARDED   = 1
	EXCLUSIVE = 2
)

type createResOk struct {
	Obj      PostOpFh3
	ObjAttrs PostOpAttr
}

type Create3Res struct {
	Status NFS3Stat    `xdr:"union"`
	ResOk  createResOk `xdr:"unioncase=0"`
	DirWcc WccData
}

type Create3Args struct {
	Where DirOpArgs3
	How   Create3How
}

type Create3How struct {
	Mode     CreateMode3 `xdr:"union"`
	ObjAttrs SAttr3      `xdr:"unioncase=1"`
	Verf     CreateVerf3 `xdr:"unioncase=2"`
}
