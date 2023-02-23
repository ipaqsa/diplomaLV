package pkg

type ConfigT struct {
	Port      string `yaml:"port"`
	Agent     string `yaml:"agent"`
	Ingress   string `yaml:"ingress"`
	Keyword   string `yaml:"keyword"`
	SKEY_SIZE uint   `yaml:"skey_size"`
	AKEY_SIZE uint   `yaml:"akey_size"`
}

var Config = ConfigT{}
