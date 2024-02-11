package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SetKeyWithDiag(state *schema.ResourceData, key string, v interface{}) *diag.Diagnostic {
	if err := state.Set(key, v); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied value for field %s was not accepted: %s", key, err.Error()),
			AttributePath: cty.GetAttrPath(key),
		}
	} else {
		return nil
	}
}
