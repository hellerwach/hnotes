package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Name = "config.yaml"

var DirPath = filepath.Join(os.Getenv("HOME"), ".config/hnotes")
var FilePath = filepath.Join(DirPath, Name)

type Config struct {
	Port int `yaml:"port"`
}

func init() {
	viper.SetConfigFile(FilePath)
	viper.AddConfigPath(DirPath)
	viper.SetConfigType("yaml")
}

func MustRead() {
	if err := Read(); err != nil {
		logrus.Fatalln(err)
	}
}

func Read() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Could not read config: %v", err)
	}
	return nil
}
