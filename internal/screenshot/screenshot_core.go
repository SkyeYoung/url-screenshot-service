package screenshot

import (
	"fmt"
	"path"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/playwright-community/playwright-go"
)

func ScreenshotCore(bctx BrowserCtx, url, folderPath string) (string, error) {
	logger := helper.GetLogger("server").Named("screenshot")
	logger.Infof("screenshotting `%v`", url)

	page := *bctx.GetPage()
	img := ""

	if _, e := page.Goto(url, playwright.PageGotoOptions{
		Referer: playwright.String("https://www.google.com/"),
	}); e != nil {
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

	logger.Infof("screenshotted `%v`", url)
	return img, nil
}
