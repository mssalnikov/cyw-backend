package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ConfigManager struct
type ConfigManager struct {
	file string
}

// NewConfig constructor
func NewConfig(file string) *ConfigManager {
	return &ConfigManager{file}
}

// Load config from file
func (cm *ConfigManager) Load() (*Config, error) {
	data, err := ioutil.ReadFile(cm.file)
	if err != nil {
		return nil, err
	}

	var config *Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Config struct
type Config struct {
	Debug    bool `yaml:"debug"`
	DB       DB   `yaml:"db"`
	Host     Host `yaml:"host"`
	AuthHost Host `yaml:"authhost"`
	Auth     Auth `yaml:"auth"`
	Navi     Navi `yaml:"navi"`
}

// DB struct
type DB struct {
	UserName     string `yaml:"username"`
	UserPassword string `yaml:"userpassword"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"database"`
}

// Host struct
type Host struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

// Host struct
type AuthHost struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

// Auth struct
type Auth struct {
	// Auth server
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
	// Facebook
	FBClient string `yaml:"fbclient"`
	FBSecret string `yaml:"fbsecret"`
}

type Navi struct {
	AuthToken string `yaml:"authtoken"`
	ApiUri    string `yaml:"apiuri"`
}
