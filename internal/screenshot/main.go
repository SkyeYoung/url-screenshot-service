package screenshot

import (
	"github.com/Jeffail/tunny"
	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
)

type Response struct {
	Url string
	Err error
}

type screenshot struct {
	pool *tunny.Pool
	ctx  *BrowserCtx
}

type Screenshot interface {
	GetPool() *tunny.Pool
	Close()
}

func New(cfg *helper.Config) Screenshot {
	ctx := NewCtx(cfg)

	pool := tunny.NewFunc(1, func(payload interface{}) interface{} {
		url := payload.(string)

		defer func() {
			if r := recover(); r != nil {
				ctx.Close()
				ctx = NewCtx(cfg)
			}
		}()

		img, err := ScreenshotCore(ctx, url, cfg.Prefix)
		if err != nil {
			panic(err)
		}

		return &Response{
			Url: img,
			Err: err,
		}
	})

	return &screenshot{
		pool: pool,
		ctx:  &ctx,
	}
}

func (s *screenshot) GetPool() *tunny.Pool {
	return s.pool
}

func (s *screenshot) Close() {
	(*s.ctx).Close()
	s.pool.Close()
}
