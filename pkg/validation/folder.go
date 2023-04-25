package validation

import "github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"

func FolderName(folderName string) error {
	if len(folderName) == 0 {
		return errors.ErrInvalidFolderName
	}

	return nil
}
