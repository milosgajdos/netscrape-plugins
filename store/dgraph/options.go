package dgraph

import "google.golang.org/grpc"

const (
	// DefaultURL is default dgraph connection URL
	DefaultURL = "localhost:9080"
)

// Options configure dgraph.
type Options struct {
	DialOpts []grpc.DialOption
	Auth     *Auth
}

// Option is dgraph option
type Option func(*Options)

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
