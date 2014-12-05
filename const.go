// const
package ini

const (
	EMPTY_STRING = ""
	R_DBL_QUOTE  = "^([^= \t]+)[ \t]*=[ \t]*\"([^\"]*)\"$"
	R_SNGL_QUOTE = "^([^= \t]+)[ \t]*=[ \t]*'([^']*)'$"
	R_NO_QUOTE   = "^([^= \t]+)[ \t]*=[ \t]*([^#;]+)"
	R_NO_VAL     = "^([^= \t]+)[ \t]*=[ \t]*([#;].*)?"
	POUND        = '#'
	SEMICOLON    = ';'
	LEFT_BRKT    = '['
	RIGHT_BRKT   = ']'
	PERMISSION   = 0644
)
