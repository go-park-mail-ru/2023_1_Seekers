package common

import (
	"path/filepath"
)

var audioTemplate = "audio.tpl"
var docsTemplate = "docs.tpl"
var defaultTemplate = "download.tpl"
var imgTemplate = "img.tpl"
var videoTemplate = "video.tpl"

var ext2tpl = map[string]string{
	// audio
	".mp3": audioTemplate,
	".wav": audioTemplate,
	".ogg": audioTemplate,
	// documents -> google.docs
	".json": docsTemplate,
	".hpp":  docsTemplate,
	".cpp":  docsTemplate,
	".py":   docsTemplate,
	".html": docsTemplate,
	".txt":  docsTemplate,
	".go":   docsTemplate,
	".doc":  docsTemplate,
	".docx": docsTemplate,
	".ppt":  docsTemplate,
	".pptx": docsTemplate,
	".xls":  docsTemplate,
	".xlsx": docsTemplate,
	".tif":  docsTemplate,
	".pdf":  docsTemplate,
	// images
	".png":  imgTemplate,
	".jpeg": imgTemplate,
	".jpg":  imgTemplate,
	".gif":  imgTemplate,
	".svg":  imgTemplate,
	".webp": imgTemplate,
	// videos
	".mp4":  videoTemplate,
	".mov":  videoTemplate,
	".webm": videoTemplate,
}

func GetTplFile(fileName string) string {
	ext := filepath.Ext(fileName)

	tpl, ok := ext2tpl[ext]
	if !ok {
		return defaultTemplate
	}

	return tpl
}
