package config

import "github.com/BurntSushi/toml"

type Config struct {
	LLM        LLMConfig `toml:"llm"`
	DB         DBConfig  `toml:"db"`
	Categories []string  `toml:"categories"`
}

type LLMConfig struct {
	APIKey string `toml:"api_key"`
	Model  string `toml:"model"`
}

type DBConfig struct {
	Address string `toml:"address"`
}

func GetConfig(file string) (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
