package screenshot

import (
	"fmt"
	"path"

	"github.com/Jeffail/tunny"
	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/playwright-community/playwright-go"
)

type Response struct {
	Url string
	Err error
}

func ScreenshotCore(bctx BrowserCtx, url, folderPath string) (string, error) {
	logger := helper.GetLogger("server").Named("screenshot")

	page := bctx.GetPage()
	img := ""

	if _, e := page.Goto(url); e != nil {
		err := fmt.Errorf("could not goto `%v`, err: %v", url, e.Error())
		logger.Warn(err)
		return img, err
	}

	img = helper.WrapImgExt(helper.EncodeImgName(url))
	p := path.Join(folderPath, img)
	if _, e := page.Screenshot(playwright.PageScreenshotOptions{
		Path:    playwright.String(p),
		Quality: playwright.Int(50),
		Type:    helper.PlaywrightImgExt(),
	}); e != nil {
		err := fmt.Errorf("could not create screenshot of `%v`, err: %v", url, e.Error())
		logger.Error(err)
		return img, err
	}

	return img, nil
}

func Pool(folder string) *tunny.Pool {
	ctx := New()
	defer ctx.Close()

	pool := tunny.NewFunc(1, func(payload interface{}) interface{} {
		url := payload.(string)

		img, err := ScreenshotCore(ctx, url, folder)

		defer func() {
			if r := recover(); r != nil {
				ctx.Close()
				ctx = New()
			}
		}()

		return &Response{
			Url: img,
			Err: err,
		}
	})

	defer pool.Close()

	return pool
}
