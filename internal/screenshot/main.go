package screenshot

import (
	"github.com/Jeffail/tunny"
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

func New(folder string) Screenshot {
	ctx := NewCtx()

	pool := tunny.NewFunc(1, func(payload interface{}) interface{} {
		url := payload.(string)

		defer func() {
			if r := recover(); r != nil {
				ctx.Close()
				ctx = NewCtx()
			}
		}()

		img, err := ScreenshotCore(ctx, url, folder)
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
