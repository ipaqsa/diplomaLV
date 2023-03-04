package pkg

type ConfigT struct {
	Port      string   `yaml:"port"`
	Path      string   `yaml:"path"`
	Agent     string   `yaml:"agent"`
	Keyword   string   `yaml:"keyword"`
	SKEY_SIZE uint     `yaml:"skey_size"`
	AKEY_SIZE uint     `yaml:"akey_size"`
	Brokers   []string `yaml:"brokers"`
}

var Config = ConfigT{}
