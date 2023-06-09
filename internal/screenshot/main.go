package screenshot

import (
	"errors"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/playwright-community/playwright-go"
)

func Screenshot(url string) (path string, err error) {
	logger := helper.GetLogger()

	p := helper.WrapImgExt(helper.EncodeImgName(path))

	browserCtx(func(page playwright.Page) {
		if _, e := page.Goto(path, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		}); e != nil {
			err = errors.New("could not goto url: " + path)
			logger.Warn(err)
			return
		}

		if _, err = page.Screenshot(playwright.PageScreenshotOptions{
			Path:    playwright.String(p),
			Quality: playwright.Int(50),
			Type:    helper.PlaywrightImgExt(),
		}); err != nil {
			err = errors.New("could not create screenshot of " + path + ", err:" + err.Error())
			logger.Error(err)
		}
	})

	return p, err
}
