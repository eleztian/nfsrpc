package nfs

import (
	"bs-2018/mynfs/nfs/nfs_progra"
	"bs-2018/mynfs/nfsrpc"
	"fmt"
	"testing"
)

var host = "111.231.215.178"

func GetNfsClient(host string) (*nfs_progra.Client, error) {
	return nfs_progra.NewNFSClientAuth(host, 1, &nfsrpc.AuthsysParms{
		100231,
		"ubuntu",
		500,
		500,
		0,
	})
}

func TestNew(t *testing.T) {
	nfs, err := New(host, "/home/ubuntu/go", "ubuntu", 500, 500)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nfs.Root.Handle)
}

func TestNFS_GetObjectHandle(t *testing.T) {
	nfs, _ := New(host, "/home/ubuntu/go/src", "ubuntu", 500, 500)
	c, _ := GetNfsClient(host)
	h, err := nfs.GetObjectHandle(c, "test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(h)
	h, err = nfs.GetObjectHandle(c, "test/test1.go")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(h)
	fmt.Println(ObjectCache)

}
