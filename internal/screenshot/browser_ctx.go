package screenshot

import (
	"log"

	"github.com/SkyeYoung/url-screenshot-service/internal/helper"
	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"
)

type browserCtx struct {
	pw      *playwright.Playwright
	browser *playwright.Browser
	page    *playwright.Page
}

type BrowserCtx interface {
	GetPage() playwright.Page
	ClosePage()
	Close()
}

func New() BrowserCtx {
	logger := helper.GetLogger("server").Named("screenshot")
	pw, err := playwright.Run()

	if err != nil {
		logger.Panicf("could not launch playwright: %v", err)
	}

	browser, err := pw.Firefox.Launch()
	if err != nil {
		logger.Panicf("could not launch browser: %v", err)
	}

	return &browserCtx{
		pw:      pw,
		browser: &browser,
		page:    nil,
	}
}

func (b *browserCtx) GetPage() playwright.Page {
	if b.page == nil {
		logger := helper.GetLogger("server").Named("screenshot")

		page, err := (*b.browser).NewPage(playwright.BrowserNewContextOptions{
			Viewport: &playwright.BrowserNewContextOptionsViewport{
				Width:  playwright.Int(1920),
				Height: playwright.Int(1080),
			},
			UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
		})
		if err != nil {
			logger.Panicf("could not create page: %v", err)
		}
		b.page = &page
	}
	return *b.page
}

func (b *browserCtx) ClosePage() {
	logger := helper.GetLogger("server").Named("screenshot")

	if b.page != nil {
		if err := (*b.page).Close(); err != nil {
			logger.Panic(errors.Wrapf(err, "could not close page"))
		}
	}
}

func (b *browserCtx) Close() {
	b.ClosePage()

	logger := helper.GetLogger("server").Named("screenshot")

	if err := (*b.browser).Close(); err != nil {
		log.Panic(errors.Wrapf(err, "could not close browser"))
	}
	if err := b.pw.Stop(); err != nil {
		logger.Panic(errors.Wrapf(err, "could not stop playwright"))
	}
}
