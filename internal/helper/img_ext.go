package helper

import "github.com/playwright-community/playwright-go"

func GetImgExt() string {
	return "jpeg"
}

func WrapImgExt(str string) string {
	return str + "." + GetImgExt()
}

func PlaywrightImgExt() *playwright.ScreenshotType {
	v := playwright.ScreenshotType(GetImgExt())
	return &v
}
