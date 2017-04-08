package smservice

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceList map[string]*ServiceConfig `yaml:"servicelist"`
	Errormsg    map[string]string         `yaml:"errormsg"`
	RedisConf   map[string]string         `yaml:"redisconf"`
	MysqlConf   map[string]string         `yaml:"mysqlconf"`
}

type ServiceConfig struct {
	Agent       string `yaml:"agent"`
	Tpl         string `yaml:"smstpl"`
	Signame     string `yaml:"signname"`
	Callback    string `yaml:"callback"`
	Maxsendnums string `yaml:"maxsendnums"`
	Validtime   string `yaml:"validtime"`
}

var (
	config     Config
	configFile = "./conf.yaml"
)

func (cfg *Config) ParseConfigData(data []byte) error {
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) ParseConfigFile(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	return cfg.ParseConfigData(data)
}

func LoadConfig() {
	// 加载配置文件
	config = Config{}
	config.ParseConfigFile(configFile)
}
