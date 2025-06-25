package modem

import (
	"domofon-api/pkg/huaweimodem"
	"log"
	"time"

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

	maxAttempts := 10
	retryInterval := 5 * time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = modem.Login()
		if err == nil {
			break
		}

		lastErr = err
		if attempt < maxAttempts {
			log.Printf("Login attempt %d/%d failed, retrying in %v: %v", attempt, maxAttempts, retryInterval, err)
			time.Sleep(retryInterval)
		}
	}

	if lastErr != nil {
		log.Fatalf("Failed to login after %d attempts: %v", maxAttempts, lastErr)
	}

	return modem
}
