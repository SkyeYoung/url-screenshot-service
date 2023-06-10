package job

import (
	"fmt"

	"os"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/pkg/errors"
	"github.com/reugn/go-quartz/quartz"
	"go.uber.org/zap"
)

type ClearLocalImgJob struct {
	JobIns
}

func NewClearLocalImgJob() *ClearLocalImgJob {
	userType := &ClearLocalImgJob{}
	userType.IJob = interface{}(userType).(IJob)
	userType.Job = interface{}(userType).(quartz.Job)
	return userType
}

func (j *ClearLocalImgJob) Description() string {
	return "Clear local images."
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
