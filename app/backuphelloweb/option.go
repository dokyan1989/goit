package backuphelloweb

type options struct {
	Env  string
	Host string
	Port int
}

func defaultConfig() *options {
	return &options{
		Env:  "dev",
		Host: "",
		Port: 3000,
	}
}

type Option func(*options)

func WithEnv(env string) Option {
	return func(opts *options) {
		opts.Env = env
	}
}

func WithHost(host string) Option {
	return func(opts *options) {
		opts.Host = host
	}
}

func WithPort(port int) Option {
	return func(opts *options) {
		opts.Port = port
	}
}
