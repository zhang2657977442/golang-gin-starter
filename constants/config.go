package constants

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/zhang2657977442/golang-gin-starter/utils/log"
)

const (
	defaultLogFileSize    = 100
	defaultLogRetainDays  = 7
	defaultLogRetainCount = 10
)

var single *Config

func Conf() *Config {
	return single
}

type Config struct {
	Global *GlobalCfg `toml:"global,omitempty" json:"global"`
	Log    *LogConfig `toml:"log,omitempty" json:"log"`
}

type LogConfig struct {
	Path           string `toml:"path,omitempty" json:"path"`
	Level          string `toml:"level,omitempty" json:"level"`
	MaxFileSize    int    `toml:"max_file_size,omitempty" json:"max_file_size"`
	MaxRetainDays  int    `toml:"max_retain_days,omitempty" json:"max_retain_days"`
	MaxRetainCount int    `toml:"max_retain_count,omitempty" json:"max_retain_count"`
	MouduleName    string `toml:"moudule_name,omitempty" json:"module_name"`
}

type GlobalCfg struct {
	Port    uint32 `toml:"port,omitempty" json:"port"`
	FileDir string `toml:"file_dir,omitempty" json:"file_dir"`
}

func (c *LogConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("log path is not specified")
	}
	if c.Level == "" {
		c.Level = "info"
	}
	if c.MaxFileSize <= 0 {
		c.MaxFileSize = defaultLogFileSize
	}
	if c.MaxRetainDays <= 0 {
		c.MaxRetainDays = defaultLogRetainDays
	}
	if c.MaxRetainCount <= 0 {
		c.MaxRetainCount = defaultLogRetainCount
	}
	return nil
}

func InitConfig(path string) error {
	single = &Config{
		Global: &GlobalCfg{},
		Log:    &LogConfig{},
	}

	LoadConfig(single, path)

	return Conf().Validate()
}

func LoadConfig(conf *Config, path string) {
	if len(path) == 0 {
		log.Error("Strat Backend", "Init", "configPath file is empty")
		os.Exit(-1)
	}

	if _, err := toml.DecodeFile(path, conf); err != nil {
		log.Error("Strat Backend", "Init", "decode:[%s] failed, err:[%s]", path, err.Error())
		os.Exit(-1)
	}
}

func (config *Config) Validate() error {
	if err := config.Log.Validate(); err != nil {
		return err
	}

	return config.validatePath()
}

func createDir(dir string) error {
	stat, err := os.Stat(dir)
	if err == nil {
		if !stat.IsDir() { // exist but not a dir
			err = errors.New("already exists but not a dir")
		}
	} else if os.IsNotExist(err) { // create if not exist
		err = os.MkdirAll(dir, 0755)
	}
	return err
}

func (config *Config) validatePath() error {
	if err := createDir(config.Log.Path); err != nil {
		return fmt.Errorf("create log dir[%s] failed: %v", config.Log.Path, err)
	}
	return nil
}
