package helper

import "encoding/base64"

func EncodeImgName(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func DecodeImgName(name string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func EncodeImgNameAddExt(name string) string {
	return WrapImgExt(EncodeImgName(name))
}
