package nfs_progra

import (
	"bs-2018/mynfs/nfsrpc"
	"fmt"
	"net"
	"testing"
	"time"
)

var host = "111.231.215.178"

func GetNfsClient(host string) (*Client, error) {
	return NewNFSClientAuth(host, 1, &nfsrpc.AuthsysParms{
		100231,
		"ubuntu",
		500,
		500,
		0,
	})
}

func GetHandle(name string) (*NFS_FH3, error) {
	c, _ := net.Dial("tcp", host+":1011")
	client := nfsrpc.NewClient(c, nil, nil)
	s, err := nfsrpc.Mnt(client, name)
	if err != nil {
		return nil, err
	}
	return &NFS_FH3{s.MountInfo.FHandle}, nil
}

func TestGetAttr(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/t2")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	result := GetAttrRes{}
	err = client.Call("NFS.GetAttr",
		GetAttr3Args{fh.Data},
		&result)
	if err != nil {
		t.Error("NFS.Null: ", err)
	}
	fmt.Println(result.Resok.ObjAttributes.Type)
}

func TestSetAttr(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/test.go")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fh)
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	result := SetAttr3Res{}
	err = client.Call("NFS.SetAttr",
		SetAttr3Aargs{
			Object: NFS_FH3{fh.Data},
			NewAttr: SAttr3{
				Uid: SetUid{true, 0},
				Gid: SetGid{true, 0},
			},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.ObjWcc.After.Attributes.MTime)
	tm := time.Unix(int64(result.ObjWcc.After.Attributes.MTime.Seconds), int64(result.ObjWcc.After.Attributes.MTime.NSeconds))
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
	fmt.Println(result.ObjWcc.After)
	tm = time.Unix(int64(result.ObjWcc.Before.Attributes.CTime.Seconds), int64(result.ObjWcc.Before.Attributes.CTime.NSeconds))
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
}

func TestLookUp(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fh)
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	result := LookUp3Res{}
	err = client.Call("NFS.Lookup",
		LookUp3Args{
			What: DirOpArgs3{Dir: *fh, Name: "t.go"},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.ResOk)

}

func TestAccess(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/t.go")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	result := Access3Res{}
	err = client.Call("NFS.Access",
		Access3Args{
			Object: *fh,
			Access: ACCESS3_MODIFY,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("access: ", result.ResOk.Access)
}

func TestReadLink(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/t2")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	result := ReadLink3Res{}
	err = client.Call("NFS.ReadLink",
		ReadLink3Args{
			SymLink: *fh,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("SymLinkData: ", result)
}

func TestRead(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/t.go")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	result := Read3Res{}
	err = client.Call("NFS.Read",
		Read3Args{
			File:  *fh,
			Count: 5,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Read: ", result.ResOk.Count, string(result.ResOk.Data))
	err = client.Call("NFS.Read",
		Read3Args{
			File:   *fh,
			Offset: 5,
			Count:  100,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Read: ", result.ResOk.Count, string(result.ResOk.Data))
}

func TestWrite(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/t.go")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := Write3Res{}
	err = client.Call("NFS.Write",
		Write3Args{
			File:   *fh,
			Offset: 84,
			Count:  uint32(len([]byte("\njust a test\n"))),
			Stable: UNSTABLE,
			Data:   []byte("\njust a test\n"),
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Write: ", result)
}

func TestCreate(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := Create3Res{}
	err = client.Call("NFS.Create",
		Create3Args{
			Where: DirOpArgs3{*fh, "newfile_test_nfs2"},
			How:   Create3How{Mode: 1, ObjAttrs: SAttr3{}},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Create: ", result)
}

func TestMkDir(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := MkDirRes{}
	err = client.Call("NFS.MkDir",
		MkDirArgs{
			Where: DirOpArgs3{*fh, "newfile_test_nfs"},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Create: ", result)
}

// TODO: CAN NOT PASS
func TestSymLink(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/newfile_test_nfs")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := SymLinkRes{}
	err = client.Call("NFS.SymLink",
		SymLinkArgs{
			Where:   DirOpArgs3{*fh, "newfile_test_nfs2"},
			SymLink: SymLinkData3{SAttr3{}, "newfile_test_nfs3"},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Create: ", result)
}

func TestRemove(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := Remove3Res{}
	err = client.Call("NFS.Remove",
		Remove3Args{
			Object: DirOpArgs3{*fh, "t.go"},
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Create: ", result)
}

func TestReadDir(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := ReadDirRes{}
	err = client.Call("NFS.ReadDir",
		ReadDir3Args{
			Dir:   *fh,
			Count: 2000, // The size must include all XDR overhead.
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	d := result.ResOk.Reply.Entries
	for {
		if d == nil {
			break
		}
		fmt.Println(d.Name, d.Fileid, d.Cookie)
		d = d.NextEntry
	}
	//err = client.Call("NFS.ReadDir",
	//	ReadDir3Args{
	//		Dir:*fh,
	//		Cookie:0,
	//		CookieVerf:result.ResOk.CookieVerf,
	//		Count:20,
	//	},
	//	&result)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println("Create: ",result)
}

func TestReadDirPlus(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := ReadDirPlusRes{}
	err = client.Call("NFS.ReadDirPlus",
		ReadDirPlus3Args{
			Dir:      *fh,
			DirCount: 20,
			MaxCount: 2000, // The size must include all XDR overhead.
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
	d := result.ResOk.Reply.Entries
	for {
		if d == nil {
			break
		}
		fmt.Println(d.Name, d.Fileid, d.Cookie, d.NameAttrs, d.NameHandle.Handle.Data)
		d = d.NextEntry
	}
	//err = client.Call("NFS.ReadDir",
	//	ReadDir3Args{
	//		Dir:*fh,
	//		Cookie:0,
	//		CookieVerf:result.ResOk.CookieVerf,
	//		Count:20,
	//	},
	//	&result)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println("Create: ",result)
}

func TestFsStat(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := FsStatRes{}
	err = client.Call("NFS.FsStat",
		FsStat3Args{
			*fh,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.ResOk.TBytes)
}

func TestFsInfo(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := FsInfoRes{}
	err = client.Call("NFS.FsStat",
		FsInfoArgs{
			*fh,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.ResOk)
}

func TestPathConf(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src")
	if err != nil {
		t.Error(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := PathConfRes{}
	err = client.Call("NFS.PathConf",
		PathConf3Args{
			*fh,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result.ResOk)
}

func TestCommit(t *testing.T) {
	fh, err := GetHandle("/home/ubuntu/go/src/test.go")
	if err != nil {
		t.Fatal(err)
	}
	client, err := GetNfsClient(host)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println(fh)
	result := Commit3Res{}
	err = client.Call("NFS.Commit",
		Commit3Args{
			*fh,
			0,
			20,
		},
		&result)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}
