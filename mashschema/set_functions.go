package mashschema

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
	"strings"
)

func StringHashcode(i interface{}) int {
	return schema.HashString(i.(string))
}

func cloneStringArray(arr *[]string) []string {
	rv := make([]string, len(*arr))
	for idx, s := range *arr {
		rv[idx] = s
	}

	return rv
}

func SortedStringOf(arr *[]string) string {
	cp := cloneStringArray(arr)
	sort.Strings(cp)
	return strings.Join(cp, ",")
}

func SortedMapOf[T any](inp *map[string]T) string {
	if inp == nil || len(*inp) == 0 {
		return ""
	}

	keys := make([]string, len(*inp))

	idx := 0
	for k, _ := range *inp {
		keys[idx] = k
		idx++
	}

	sort.Strings(keys)

	sb := strings.Builder{}
	for _, k := range keys {
		var v interface{}
		v = (*inp)[k]
		sb.WriteString(fmt.Sprintf("%s::=%s/", k, v))
	}

	return sb.String()
}
