package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type StringIdentityFieldMapper[MType any] struct {
	FieldMapperBase[MType]

	Locator         LocatorFunc[MType, string]
	PreviousLocator LocatorFunc[MType, string]
}

func (sfm *StringIdentityFieldMapper[MType]) NilRemote(state *schema.ResourceData) *diag.Diagnostic {
	if err := state.Set(sfm.Key, ""); err != nil {
		return &diag.Diagnostic{
			Severity:      diag.Error,
			Detail:        fmt.Sprintf("supplied null-value for field %s was not accepted: %s", sfm.Key, err.Error()),
			AttributePath: cty.GetAttrPath(sfm.Key),
		}
	} else {
		return nil
	}
}

func (sfm *StringIdentityFieldMapper[MType]) RemoteToSchema(v *MType, state *schema.ResourceData) *diag.Diagnostic {
	// If the code has defined a null-checker function, this makes this mapper
	// writeable
	val := sfm.Locator(v)

	if err := state.Set(sfm.Key, val); err != nil {
		return &diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("unable to set string value: %s", err.Error()),
		}
	}

	return nil
}

func (sfm *StringIdentityFieldMapper[MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	*sfm.Locator(remote) = state.Id()
}

func (im *StringIdentityFieldMapper[MType]) Identity(state *schema.ResourceData) (string, error) {
	return state.Id(), nil
}

func (im *StringIdentityFieldMapper[MType]) Assign(ident string, state *schema.ResourceData) error {
	state.SetId(ident)
	return nil
}
