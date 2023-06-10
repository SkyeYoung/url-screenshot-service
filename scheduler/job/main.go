package job

import (
	"context"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/reugn/go-quartz/quartz"
	"go.uber.org/zap"
)

type JobCtxValKey string

const (
	CFG    JobCtxValKey = "cfg"
	LOGGER JobCtxValKey = "logger"
)

type IJob interface {
	ExecuteCore(logger *zap.SugaredLogger, cfg *helper.Config) (string, error)
}

type JobIns struct {
	quartz.Job
	IJob
}

func (j *JobIns) Key() int {
	return quartz.HashCode(j.Description())
}

func (j *JobIns) Description() string {
	return ""
}

func (j *JobIns) ExecuteCore(logger *zap.SugaredLogger, cfg *helper.Config) (string, error) {
	return "", nil
}

func (j *JobIns) Execute(ctx context.Context) {
	cfg := ctx.Value(CFG).(*helper.Config)
	logger := ctx.Value(LOGGER).(*zap.SugaredLogger)

	logger.Infof("Job `%v` started", j.Job.Description())

	info, err := j.IJob.ExecuteCore(logger, cfg)
	if err != nil {
		logger.Errorf("Job `%v` failed: %v", j.Job.Description(), err)
		return
	}
	logger.Infof("Job `%v` result: %v", j.Job.Description(), info)

	logger.Infof("Job `%v` finished", j.Job.Description())
}
