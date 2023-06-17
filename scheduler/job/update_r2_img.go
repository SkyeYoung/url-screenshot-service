package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/SkyeYoung/url-screenshot-service/internal/r2"
	"github.com/SkyeYoung/url-screenshot-service/internal/screenshot"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type UpdateR2ImgJob struct {
	JobIns
}

func (j *UpdateR2ImgJob) Name() string {
	return "UpdateR2ImgJob"
}

func (j *UpdateR2ImgJob) Description() string {
	return "Update R2 images."
}

func (j *UpdateR2ImgJob) Execute(ctx context.Context) {
	ExecuteWrapper(ctx, j)
}

func (j *UpdateR2ImgJob) ExecuteCore(logger *zap.SugaredLogger, cfg *helper.Config) (string, error) {
	r2 := r2.New(cfg)

	logger.Infof("listing objects with prefix %v need update", cfg.Prefix)
	objs, err := r2.ListObjects(&cfg.Prefix, func(obj *s3.Object) bool {
		return obj.LastModified.Before(time.Now().Add(-time.Hour * 24))
	})

	if err != nil {
		return "", err
	}

	ss := screenshot.New(cfg.Prefix)
	defer ss.Close()

	for _, obj := range objs {
		logger.Infof("checking image key %v", *obj.Key)
		url := strings.TrimPrefix(*obj.Key, cfg.Prefix+"/")
		url = strings.TrimSuffix(url, "."+helper.GetImgExt())
		url, _ = helper.DecodeImgName(url)
		logger.Infof("decoded url: %v", url)
		url, err = helper.GetValidUrl(url)

		if err != nil {
			logger.Warnf(err.Error())
			continue
		}

		logger.Infof("trying to get screeshot of %v", url)
		if res := ss.GetPool().Process(url).(*screenshot.Response); res.Err != nil {
			logger.Warnf(res.Err.Error())
			continue
		}

		info, err := r2.UploadObject(obj.Key)
		if err != nil {
			return "", err
		}
		logger.Infof("screenshot uploaded to %v", info.Location)

		if err := helper.RmImgAfterUpload(cfg, logger, url, *obj.Key); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("updated %v image files", len(objs)), nil
}
