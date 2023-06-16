package screenshot

import (
	"errors"
	"path"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/playwright-community/playwright-go"
)

func Screenshot(url, folderPath string) (string, error) {
	logger := helper.GetLogger("server").Named("screenshot")

	var img string
	err := browserCtx(func(page playwright.Page) error {
		if _, e := page.Goto(url); e != nil {
			err := errors.New("could not goto url: " + url + ", err:" + e.Error())
			logger.Warn(err)
			return err
		}

		img = helper.WrapImgExt(helper.EncodeImgName(url))
		p := path.Join(folderPath, img)
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
