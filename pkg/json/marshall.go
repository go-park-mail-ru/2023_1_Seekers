package json

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2023_1_Seekers/pkg/errors"
	"github.com/mailru/easyjson"
	pkgErrors "github.com/pkg/errors"
)

func MarshalEasyJSON(dataStruct interface{}) ([]byte, error) {
	structMarshaller, ok := dataStruct.(easyjson.Marshaler)
	if !ok {
		return json.Marshal(dataStruct)
	}

	out, err := easyjson.Marshal(structMarshaller)
	if !ok {
		return []byte{}, pkgErrors.WithMessage(errors.ErrInternal, "failed cast to marshall with easyjson")
	}

	return out, err
}
