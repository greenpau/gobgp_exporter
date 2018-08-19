/*
   Extended GoBGP API Client by Paul Greenberg (@greenpau).
*/

package gobgpapi

import (
	grpc "google.golang.org/grpc"
	"time"
)

type GobgpApiExtendedClient struct {
	Gobgp GobgpApiClient
	Conn  *grpc.ClientConn
}

func NewGobgpApiExtendedClient(s string, t int) (GobgpApiExtendedClient, error) {
	if t == 0 {
		t = 2
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTimeout(time.Duration(t)*time.Second))
	conn, err := grpc.Dial(s, opts...)
	if err != nil {
		return GobgpApiExtendedClient{}, err
	}
	apiInterface := &gobgpApiClient{conn}
	c := GobgpApiExtendedClient{
		Gobgp: apiInterface,
		Conn:  conn,
	}
	return c, nil
}

func (c *GobgpApiExtendedClient) Close() error {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
