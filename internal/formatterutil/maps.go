package formatterutil

import (
	"sort"
	"strings"
)

const DefaultListSeparator string = ","

func JoinMapKeysGeneric[T any](m map[string]T, separator string) string {
	if len(m) == 0 {
		return ""
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return strings.Join(keys, separator)
}
