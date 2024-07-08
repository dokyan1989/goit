package api

type config struct {
	Env  string
	Host string
	Port int
}

func defaultConfig() *config {
	return &config{
		Env:  "dev",
		Host: "",
		Port: 3000,
	}
}

type Option func(*config)

func WithEnv(env string) Option {
	return func(cfg *config) {
		cfg.Env = env
	}
}

func WithHost(host string) Option {
	return func(cfg *config) {
		cfg.Host = host
	}
}

func WithPort(port int) Option {
	return func(cfg *config) {
		cfg.Port = port
	}
}
