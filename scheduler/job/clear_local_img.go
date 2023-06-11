package job

import (
	"context"
	"fmt"

	"os"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ClearLocalImgJob struct {
	JobIns
}

func (j *ClearLocalImgJob) Name() string {
	return "ClearLocalImgJob"
}

func (j *ClearLocalImgJob) Description() string {
	return "Clear local images."
}

func (j *ClearLocalImgJob) Execute(ctx context.Context) {
	ExecuteWrapper(ctx, j)
}

func (j *ClearLocalImgJob) ExecuteCore(logger *zap.SugaredLogger, cfg *helper.Config) (string, error) {
	folderPath := cfg.Prefix
	files, err := os.ReadDir(folderPath)
	if err != nil {

		return "", errors.Wrap(err, fmt.Sprintf("Error reading directory `%v`", folderPath))
	}

	for _, file := range files {
		if !file.IsDir() {
			p := folderPath + "/" + file.Name()
			err := os.Remove(p)
			if err != nil {
				return "", errors.Wrap(err, fmt.Sprintf("Error deleting file `%v`", p))
			}
		}
	}

	return fmt.Sprintf("removed %v image files", len(files)), nil
}
