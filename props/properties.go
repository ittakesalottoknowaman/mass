package props

import (
	"mass/utils/file"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ResultPath string `yaml:"result_path"`
	Auth       []*struct {
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		PrivateKey string `yaml:"private_key"`
	} `yaml:"auth"`
}

func ParseConfig(config string) (*Config, error) {
	b, err := file.ToBytes(config)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
