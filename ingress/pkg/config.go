package pkg

type ConfigT struct {
	Port    string   `yaml:"port"`
	Brokers []string `yaml:"brokers"`
}

var Config ConfigT
