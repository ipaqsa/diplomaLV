package pkg

type ConfigT struct {
	Port          string `yaml:"port"`
	Agent         string `yaml:"agent"`
	Keyword       string `yaml:"keyword"`
	SKEY_SIZE     uint   `yaml:"skey_size"`
	AKEY_SIZE     uint   `yaml:"akey_size"`
	CacheCapacity uint   `yaml:"cacheCapacity"`
	CacheTimeout  uint   `yaml:"cacheTimeout"`
}

var Config = ConfigT{}
