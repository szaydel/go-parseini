// vars
package ini

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	atoi           = strconv.Atoi
	tolower        = strings.ToLower
	trim           = strings.TrimFunc
	fmtInt         = strconv.FormatInt
	fmtFloat       = strconv.FormatFloat
	fmtBool        = strconv.FormatBool
	pFloat         = strconv.ParseFloat
	regDoubleQuote = regexp.MustCompile(R_DBL_QUOTE)
	regSingleQuote = regexp.MustCompile(R_SNGL_QUOTE)
	regNoQuote     = regexp.MustCompile(R_NO_QUOTE)
	regNoValue     = regexp.MustCompile(R_NO_VAL)
)
