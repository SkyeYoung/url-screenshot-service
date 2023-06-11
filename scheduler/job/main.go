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
	Name() string
	ExecuteCore(logger *zap.SugaredLogger, cfg *helper.Config) (string, error)
}

type IJobIns interface {
	quartz.Job
	IJob
}

type JobIns struct {
	IJobIns
}

func (j *JobIns) Key() int {
	return quartz.HashCode(j.Description())
}

func ExecuteWrapper(ctx context.Context, job IJobIns) {
	cfg := ctx.Value(CFG).(*helper.Config)
	logger := ctx.Value(LOGGER).(*zap.SugaredLogger).Named(job.Name())

	logger.Infof("Job `%v` started", job.Description())

	info, err := job.ExecuteCore(logger, cfg)
	if err != nil {
		logger.Errorf("Job `%v` failed: %v", job.Description(), err)
		return
	}
	logger.Infof("Job `%v` result: %v", job.Description(), info)

	logger.Infof("Job `%v` finished", job.Description())
}
