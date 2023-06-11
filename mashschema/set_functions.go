package mashschema

import (
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
