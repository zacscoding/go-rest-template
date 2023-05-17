package maskingutil

import (
	"regexp"
	"strings"
)

// MaskPassword masks password field
func MaskPassword(value string) string {
	regex := regexp.MustCompile(`^(?P<protocol>.+?//)?(?P<username>.+?):(?P<password>.+?)@(?P<address>.+)$`)
	if !regex.MatchString(value) {
		return value
	}
	matches := regex.FindStringSubmatch(value)
	for i, v := range regex.SubexpNames() {
		if "password" == v {
			value = strings.ReplaceAll(value, matches[i], "****")
		}
	}
	return value
}
