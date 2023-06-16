package screenshot

import (
	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/playwright-community/playwright-go"
)

func browserCtx(callback func(page playwright.Page) error) error {
	logger := helper.GetLogger("server").Named("screenshot")
	pw, err := playwright.Run()

	if err != nil {
		logger.Fatalf("could not launch playwright: %v", err)
	}
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{})
	if err != nil {
		logger.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage(playwright.BrowserNewContextOptions{
		Viewport: &playwright.BrowserNewContextOptionsViewport{
			Width:  playwright.Int(1920),
			Height: playwright.Int(1080),
		},
		UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
	})
	if err != nil {
		logger.Fatalf("could not create page: %v", err)
	}

	e := callback(page)

	if err = browser.Close(); err != nil {
		logger.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		logger.Fatalf("could not stop Playwright: %v", err)
	}

	return e
}
