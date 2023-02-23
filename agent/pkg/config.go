package pkg

type ConfigT struct {
	Port      string   `yaml:"port"`
	DBPath    string   `yaml:"dbpath"`
	SKEY_SIZE uint     `yaml:"skey_size"`
	AKEY_SIZE uint     `yaml:"akey_size"`
	KEYWORD   string   `yaml:"keyword"`
	Services  []string `yaml:"services"`
	Brokers   []string `yaml:"brokers"`
}

var Config = ConfigT{}
