package screenshot

import (
	"errors"
	"path"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/playwright-community/playwright-go"
	"go.uber.org/zap"
)

func Screenshot(url, tmpFolder string) (string, error) {
	logger := helper.GetLogger("server").With(zap.Namespace("screenshot"))

	var img string
	err := browserCtx(func(page playwright.Page) error {
		if _, e := page.Goto(url, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		}); e != nil {
			err := errors.New("could not goto url: " + url)
			logger.Warn(err)
			return err
		}

		img = helper.WrapImgExt(helper.EncodeImgName(url))
		p := path.Join(tmpFolder, img)
		if _, err := page.Screenshot(playwright.PageScreenshotOptions{
			Path:    playwright.String(p),
			Quality: playwright.Int(50),
			Type:    helper.PlaywrightImgExt(),
		}); err != nil {
			err = errors.New("could not create screenshot of " + url + ", err:" + err.Error())
			logger.Error(err)
			return err
		}

		return nil
	})

	return img, err
}
