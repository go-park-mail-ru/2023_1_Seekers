package common

import (
	"fmt"
	"math"
)

func ByteSize2Str(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d Б", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	res := float64(b) / float64(div)
	_, f := math.Modf(res)

	if f == 0 {
		return fmt.Sprintf("%.0f %cБ", res, []rune("КМГТПЕ")[exp])
	} else {
		return fmt.Sprintf("%.1f %cБ", res, []rune("КМГТПЕ")[exp])
	}
}
