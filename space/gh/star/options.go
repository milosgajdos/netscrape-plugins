package star

// Options provides GitHub scraper options.
type Options struct {
	// User s GitHub username
	User string
	// Paging size for GitHub API results
	Paging int
	// Workers for mapping repos
	Workers int
}

// Option is GitHub scraper option.
type Option func(*Options)

// Paging configures GitHub results paging
func Paging(p int) Option {
	return func(o *Options) {
		o.Paging = p
	}
}

// User configures GitHub username.
func User(u string) Option {
	return func(o *Options) {
		o.User = u
	}
}

// Workers configures GitHub API results paging.
func Workers(w int) Option {
	return func(o *Options) {
		o.Workers = w
	}
}
