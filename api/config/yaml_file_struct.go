package config

type Document struct {
	Loggers    []Logger    `yaml:"loggers"`
	Providers  []Provider  `yaml:"providers"`
	Servers    []Server    `yaml:"servers"`
	SQLOptions []SQLOption `yaml:"sqloptions"`
	Backends   []Backend   `yaml:"backends"`
}

type Logger struct {
	Kind string            `yaml:"kind"`
	Env  map[string]string `yaml:"env"`
}

type Provider struct {
	Kind string            `yaml:"kind"`
	Env  map[string]string `yaml:"env"`
}

type Server struct {
	Kind string            `yaml:"kind"`
	Env  map[string]string `yaml:"env"`
}

type Backend struct {
	Kind string            `yaml:"kind"`
	Env  map[string]string `yaml:"env"`
}

type SQLOption struct {
	Kind string            `yaml:"kind"`
	Env  map[string]string `yaml:"env"`
}
