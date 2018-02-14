package nfs

import (
	"bs-2018/mynfs/nfs/nfs_progra"
	"bs-2018/mynfs/nfsrpc"
	"fmt"
	"net"
	"strings"
)

var ObjectCache map[string]Object

func init() {
	ObjectCache = make(map[string]Object)
}

type Object struct {
	Name   string
	Handle nfs_progra.NFS_FH3
	Attrs  nfs_progra.Fattr3
}

type NFS struct {
	Host string
	Root struct {
		Path string
		Object
	}
	AuthWay int32
	Auth    interface{}
}

func New(host, root, username string, uid, gid uint32, authWay int32, auth interface{}) (*NFS, error) {
	nfs := &NFS{
		Host: host,
		AuthWay:authWay,
		Auth: auth,
	}
	conn, err := net.Dial("tcp", host+":1011")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := nfsrpc.NewClient(conn, nfsrpc.Auth(nfsrpc.AuthFlavor(authWay), auth), nil)
	res, err := nfsrpc.Mnt(client, root)
	if err != nil {
		return nil, err
	}
	nfs.Root.Path = root
	nfs.Root.Handle = nfs_progra.NFS_FH3{res.MountInfo.FHandle}

	return nfs, nil
}

func (n *NFS) GetObjectAttr(client *nfs_progra.Client, path string) (*Object, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if object, ok := ObjectCache[path]; ok {
		return &object, nil
	}

	path = path[1:]
	dirs := strings.Split(path, "/")
	dirObject := n.Root.Object
	dirPath := ""
	for _, dir := range dirs {
		res := nfs_progra.LookUp3Res{}
		err := client.Call("NFS.Lookup", nfs_progra.LookUp3Args{
			What: nfs_progra.DirOpArgs3{
				Dir:  dirObject.Handle,
				Name: dir,
			},
		}, &res)
		if err != nil {
			return nil, err
		}
		if err = res.Status.Error(); err != nil {
			return nil, fmt.Errorf("%s : %v", dirPath+"/"+dir, err)
		}
		dirPath += "/" + dir
		dirObject.Handle = res.ResOk.Object
		dirObject.Attrs = res.ResOk.ObjAttrs.Attributes
		dirObject.Name = dir
		ObjectCache[dirPath] = dirObject
	}
	return &dirObject, nil
}
