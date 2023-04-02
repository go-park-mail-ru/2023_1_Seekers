package image

import "github.com/go-park-mail-ru/2023_1_Seekers/pkg"

func CheckImageContentType(contentType string) bool {
	if contentType == pkg.ContentTypePNG || contentType == pkg.ContentTypeSVG ||
		contentType == pkg.ContentTypeWEBP || contentType == pkg.ContentTypeJPEG {
		return true
	}
	return false
}
