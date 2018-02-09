package nfs

// GETATTR3res NFSPROC3_GETATTR(GETATTR3args) = 1;

type GetAttr3Args struct {
	Object Opaque
}

type GetAttr3Resok struct {
	ObjAttributes Fattr3
}

type GetAttrRes struct {
	NFSStat3 NFS3Stat `xdr:"union"`
	Resok GetAttr3Resok `xdr:"unioncase=0"`
}

func (c *Client) GetAttr(args *GetAttr3Args) (*GetAttrRes, error){
	var result GetAttrRes
	err := c.Call("NFS.GetAttr", args, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
