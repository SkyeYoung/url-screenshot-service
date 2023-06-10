package scheduler

import (
	"context"
	"sync"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/scheduler/job"
	"github.com/reugn/go-quartz/quartz"
)

type Scheduler interface {
	Start(wg *sync.WaitGroup)
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

func (s *scheduler) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	sch := *s.sch
	logger := helper.GetLogger("scheduler")

	ctx := context.WithValue(context.Background(), job.CFG, s.cfg)
	ctx = context.WithValue(ctx, job.LOGGER, logger)

	logger.Info("Starting scheduler...")
	sch.Start(ctx)

	// jobs := []func() *job.JobIns{
	// 	"ClearLocalImgJob": job.NewClearLocalImgJob,
	// 	"UpdateR2ImgJob":   job.NewUpdateR2ImgJob,
	// }
	// cfgr := reflect.Indirect(reflect.ValueOf(s.cfg))
	// for k := range jobs {
	// 	info := cfgr.FieldByName(k).Interface().(helper.JobConfig)
	// 	if !info.Disable {
	// 		job := v()
	// 		logger.Infof("Scheduling job `%v`...", job.Description())
	// 		cron, err := quartz.NewCronTrigger(info.Cron)
	// 		if err != nil {
	// 			logger.Fatal(err)
	// 		}
	// 		sch.ScheduleJob(ctx, job, cron)
	// 	}
	// }
	if !s.cfg.ClearLocalImgJob.Disable {
		job := job.NewClearLocalImgJob()
		logger.Infof("Scheduling job `%v`...", job.Description())
		cron, err := quartz.NewCronTrigger(s.cfg.ClearLocalImgJob.Cron)
		if err != nil {
			logger.Fatal(err)
		}
		sch.ScheduleJob(ctx, job, cron)
	}

	defer logger.Info("Scheduler stopped")
	sch.Wait(ctx)
}
