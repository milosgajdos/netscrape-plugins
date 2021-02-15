package dgraph

import (
	"github.com/milosgajdos/netscrape/pkg/uuid"
	"google.golang.org/grpc"
)

const (
	// DefaultURL is default dgraph connection URL
	DefaultURL = "localhost:9080"
)

// Options configure dgraph.
type Options struct {
	UID      uuid.UID
	DialOpts []grpc.DialOption
	Auth     *Auth
}

// Option is dgraph option
type Option func(*Options)

// WithUID sets UID Options.
func WithUID(u uuid.UID) Option {
	return func(o *Options) {
		o.UID = u
	}
}

// WithDialOpts configure dgraph dial options.
func WithDialOpts(d ...grpc.DialOption) Option {
	return func(o *Options) {
		o.DialOpts = d
	}
}

// WithAuth configures dgraph Auth.
func WithAuth(a *Auth) Option {
	return func(o *Options) {
		o.Auth = a
	}
}
