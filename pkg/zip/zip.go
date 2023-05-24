package zip

import (
	"archive/zip"
	"github.com/pkg/errors"
	"time"
)

func Append2Zip(filename string, data []byte, zipw *zip.Writer) error {
	fileHeaders := zip.FileHeader{
		Name:     filename,
		Modified: time.Now(),
	}

	wr, err := zipw.CreateHeader(&fileHeaders)
	if err != nil {
		return errors.Wrap(err, "Failed to create zip entry")
	}

	if _, err := wr.Write(data); err != nil {
		return errors.Wrap(err, "Failed to write zip entry")
	}

	return nil
}
