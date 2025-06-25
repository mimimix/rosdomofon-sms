package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	SecretKey      string `yaml:"SECRET_KEY" mapstructure:"SECRET_KEY"`
	ProtectionCode string `yaml:"PROTECTION_CODE" mapstructure:"PROTECTION_CODE"`
	KeyId          int    `yaml:"KEY_ID" mapstructure:"KEY_ID"`
	HttpPort       int    `yaml:"HTTP_PORT" mapstructure:"HTTP_PORT"`
	RefreshToken   string `yaml:"REFRESH_TOKEN" mapstructure:"REFRESH_TOKEN"`
	ModemUrl       string `yaml:"MODEM_URL" mapstructure:"MODEM_URL"`
	LastSmsFile    string `yaml:"LAST_SMS_FILE" mapstructure:"LAST_SMS_FILE"`
	SmsAliveTime   int    `yaml:"SMS_ALIVE_TIME" mapstructure:"SMS_ALIVE_TIME"`
}

func Load() (*Config, error) {
	var cfg Config

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var envPath string
	if strings.HasPrefix(wd, "/app") {
		wd = "/app"
		envPath = filepath.Join(wd, "conf.yml")
	} else {
		wd = filepath.Join(wd)
		envPath = filepath.Join(wd, "../../conf.yml")
	}

	fmt.Printf("Loading config from %s\n", envPath)

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, err
	}

	viper.SetConfigFile(envPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Установка значений по умолчанию
	//setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Валидация обязательных полей
	//if err := validateConfig(&cfg); err != nil {
	//	return nil, err
	//}

	return &cfg, nil
}
