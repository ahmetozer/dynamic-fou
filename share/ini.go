package share

import "strings"

// Get ini value from string
func IniVal(text, key string) string {
	var val string
	for _, each_ln := range strings.Split(text, "\n") {
		if strings.HasPrefix(each_ln, key+"=") {
			val = strings.Split(each_ln, "=")[1]
			break
		}
	}
	return val
}
