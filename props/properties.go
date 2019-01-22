package props

import (
	"fmt"

	"mass/utils/file"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ResultPath string `yaml:"result_path"`
	Auth       *struct {
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		PrivateKey string `yaml:"private_key"`
	} `yaml:"auth"`
}

func ParseConfig(config string) (*Config, error) {
	if config == "" {
		return nil, fmt.Errorf("")
	}

	b, err := file.ToBytes(config)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	c := new(Config)
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return c, nil
}
