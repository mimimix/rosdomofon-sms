package modem

import (
	"domofon-api/pkg/huaweimodem"
	"log"

	"domofon-api.gg/config"

	"go.uber.org/zap"
)

func New(config *config.Config) *huaweimodem.Device {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	modem, err := huaweimodem.NewDevice(sugar, config.ModemUrl, "", "")
	if err != nil {
		log.Fatalf("Failed to create modem: %v", err)
	}

	err = modem.Login()
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	return modem
}
