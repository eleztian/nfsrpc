package nfs

import "github.com/pkg/errors"

// Procedure NULL does not do any work. It is made available to
// allow server response testing and timing.
func (c *Client) Null() error {
	if c == nil {
		return errors.New("invalied client")
	}
	return c.Call("NFS.Null", nil, nil)
	return nil
}
