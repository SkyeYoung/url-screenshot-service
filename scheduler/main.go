package scheduler

import (
	"context"
	"reflect"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/scheduler/job"
	"github.com/reugn/go-quartz/quartz"
)

type Scheduler interface {
	Start()
}

type scheduler struct {
	cfg *helper.Config
	sch *quartz.Scheduler
}

func New(cfg *helper.Config) Scheduler {
	sch := quartz.NewStdScheduler()
	return &scheduler{
		cfg: cfg,
		sch: &sch,
	}
}

func (s *scheduler) Start() {
	sch := *s.sch
	logger := helper.GetLogger("scheduler")

	ctx := context.WithValue(context.Background(), job.CFG, s.cfg)
	ctx = context.WithValue(ctx, job.LOGGER, logger)

	logger.Info("Starting scheduler...")
	sch.Start(ctx)

	jobs := map[string]quartz.Job{
		"ClearLocalImgJob": new(job.ClearLocalImgJob),
		"UpdateR2ImgJob":   new(job.UpdateR2ImgJob),
	}
	cfgr := reflect.Indirect(reflect.ValueOf(s.cfg))
	for key, job := range jobs {
		info := cfgr.FieldByName(key).Interface().(helper.JobConfig)
		if !info.Disable {
			logger.Infof("Scheduling job `%v`...", job.Description())
			cron, err := quartz.NewCronTrigger(info.Cron)
			if err != nil {
				logger.Fatal(err)
			}
			sch.ScheduleJob(ctx, job, cron)
		}
	}

	defer logger.Info("Scheduler stopped")
	sch.Wait(ctx)
}
