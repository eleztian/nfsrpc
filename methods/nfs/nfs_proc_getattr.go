package nfs

// GETATTR3res NFSPROC3_GETATTR(GETATTR3args) = 1;

type GetAttr3Args struct {
	Object Opaque
}

type GetAttr3Resok struct {
	ObjAttributes
}

func (c *Client) GetAttr() (error){

	return nil
}
