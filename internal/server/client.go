package server

import "context"

type Client interface {
	Run(ctx context.Context) error
	Close(ctx context.Context) error
	ClientID() string
	UUID() string
}
