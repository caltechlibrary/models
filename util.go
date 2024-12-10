package models

import (
	"regexp"
	"sort"
)

// IsValidVarname tests a sting confirms it conforms to Model's naming rule.
func IsValidVarname(s string) bool {
	if len(s) == 0 {
		return false
	}
	// NOTE: variable names must start with a latter and maybe followed by
	// one or more letters, digits and underscore.
	vRe := regexp.MustCompile(`^([a-zA-Z]|[a-zA-Z][0-9a-zA-Z\_]+)$`)
	return vRe.Match([]byte(s))
}

// getAttributeIds returns a list of attribue keys in a maps[string]interface{} structure
func getAttributeIds(m map[string]string) []string {
	ids := []string{}
	for k, _ := range m {
		if k != "" {
			ids = append(ids, k)
		}
	}
	if len(ids) > 0 {
		sort.Strings(ids)
	}
	return ids
}
