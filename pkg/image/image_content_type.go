package image

import (
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/common"
)

func CheckImageContentType(contentType string) bool {
	if contentType == common.ContentTypePNG || contentType == common.ContentTypeSVG ||
		contentType == common.ContentTypeWEBP || contentType == common.ContentTypeJPEG {
		return true
	}
	return false
}
