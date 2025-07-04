package spacemq

import "github.com/redis/go-redis/v9"

type Client struct {
	client redis.UniversalClient
}

type ClientOpts struct{}

func NewClient(c redis.UniversalClient, opts ...ClientOpts) *Client {
	return &Client{}
}
