package nfs_progra

/*
Procedure 16: READDIR - Read From Directory

   SYNOPSIS

      READDIR3res NFSPROC3_READDIR(READDIR3args) = 16;

      struct READDIR3args {
           nfs_fh3      dir;
           cookie3      cookie;
           cookieverf3  cookieverf;
           count3       count;
      };
      struct entry3 {
           fileid3      fileid;
           filename3    name;
           cookie3      cookie;
           entry3       *nextentry;
      };

      struct dirlist3 {
           entry3       *entries;
           bool         eof;
      };

      struct READDIR3resok {
           post_op_attr dir_attributes;
           cookieverf3  cookieverf;
           dirlist3     reply;
      };

      struct READDIR3resfail {
           post_op_attr dir_attributes;
      };

      union READDIR3res switch (nfsstat3 status) {
      case NFS3_OK:
           READDIR3resok   resok;
      default:
           READDIR3resfail resfail;
      };

   DESCRIPTION

      Procedure READDIR retrieves a variable number of entries,
      in sequence, from a directory and returns the name and
      file identifier for each, with information to allow the
      client to request additional directory entries in a
      subsequent READDIR request. On entry, the arguments in
      READDIR3args are:

      dir
         The file handle for the directory to be read.

      cookie
         This should be set to 0 in the first request to read
         the directory. On subsequent requests, it should be a
         cookie as returned by the server.

      cookieverf
         This should be set to 0 in the first request to read
         the directory. On subsequent requests, it should be a
         cookieverf as returned by the server. The cookieverf
         must match that returned by the READDIR in which the
         cookie was acquired.

      count
         The maximum size of the READDIR3resok structure, in
         bytes.  The size must include all XDR overhead. The
         server is free to return less than count bytes of
         data.

      On successful return, READDIR3res.status is NFS3_OK and
      READDIR3res.resok contains:

      dir_attributes
         The attributes of the directory, dir.

      cookieverf
         The cookie verifier.

      reply
         The directory list:

         entries
            Zero or more directory (entry3) entries.

         eof
            TRUE if the last member of reply.entries is the last
            entry in the directory or the list reply.entries is
            empty and the cookie corresponded to the end of the
            directory. If FALSE, there may be more entries to
            read.

      Otherwise, READDIR3res.status contains the error on
      failure and READDIR3res.resfail contains the following:

      dir_attributes
         The attributes of the directory, dir.

   IMPLEMENTATION

      In the NFS version 2 protocol, each directory entry
      returned included a cookie identifying a point in the
      directory. By including this cookie in a subsequent
      READDIR, the client could resume the directory read at any
      point in the directory.  One problem with this scheme was
      that there was no easy way for a server to verify that a
      cookie was valid. If two READDIRs were separated by one or
      more operations that changed the directory in some way
      (for example, reordering or compressing it), it was
      possible that the second READDIR could miss entries, or
      process entries more than once. If the cookie was no
      longer usable, for example, pointing into the middle of a
      directory entry, the server would have to either round the
      cookie down to the cookie of the previous entry or round
      it up to the cookie of the next entry in the directory.
      Either way would possibly lead to incorrect results and
      the client would be unaware that any problem existed.

      In the NFS version 3 protocol, each READDIR request
      includes both a cookie and a cookie verifier. For the
      first call, both are set to 0.  The response includes a
      new cookie verifier, with a cookie per entry.  For
      subsequent READDIRs, the client must present both the
      cookie and the corresponding cookie verifier.  If the
      server detects that the cookie is no longer valid, the
      server will reject the READDIR request with the status,
      NFS3ERR_BAD_COOKIE. The client should be careful to
      avoid holding directory entry cookies across operations
      that modify the directory contents, such as REMOVE and
      CREATE.

      One implementation of the cookie-verifier mechanism might
      be for the server to use the modification time of the
      directory. This might be overly restrictive, however. A
      better approach would be to record the time of the last
      directory modification that changed the directory
      organization in a way that would make it impossible to
      reliably interpret a cookie. Servers in which directory
      cookies are always valid are free to use zero as the
      verifier always.

      The server may return fewer than count bytes of
      XDR-encoded entries.  The count specified by the client in
      the request should be greater than or equal to FSINFO
      dtpref.

      Since UNIX clients give a special meaning to the fileid
      value zero, UNIX clients should be careful to map zero
      fileid values to some other value and servers should try
      to avoid sending a zero fileid.

   ERRORS

      NFS3ERR_IO
      NFS3ERR_ACCES
      NFS3ERR_NOTDIR
      NFS3ERR_BAD_COOKIE
      NFS3ERR_TOOSMALL
      NFS3ERR_STALE
      NFS3ERR_BADHANDLE
      NFS3ERR_SERVERFAULT
*/

type ReadDir3Args struct {
	Dir        NFS_FH3
	Cookie     uint64
	CookieVerf CookieVerf3
	Count      uint32
}

type Entry3 struct {
	Fileid    uint64
	Name      string
	Cookie    uint64
	NextEntry *Entry3 `xdr:"optional"`
}

type DirList3 struct {
	Entries *Entry3 `xdr:"optional"`
	Eof     bool
}

type readDir3ResOk struct {
	CookieVerf CookieVerf3
	Reply      DirList3
}

type ReadDirRes struct {
	Status   NFS3Stat `xdr:"union"`
	DirAttrs PostOpAttr
	ResOk    readDir3ResOk `xdr:"unioncase=0"`
}
