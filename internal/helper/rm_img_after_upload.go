package helper

import (
	"os"

	"go.uber.org/zap"
)

func RmImgAfterUpload(cfg *Config, logger *zap.SugaredLogger, url, key string) error {
	if cfg.RmImgAfterUpload {
		logger.Infof("removing local screenshot of %v", url)
		if err := os.Remove(key); err != nil {
			return err
		}
	}

	return nil
}
