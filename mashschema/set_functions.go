package mashschema

import (
	"github.com/hashicorp/terraform/helper/hashcode"
)

func StringHashcode(i interface{}) int {
	return hashcode.String(i.(string))
}
