package tfmapper

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"terraform-provider-mashery/funcsupport"
	"terraform-provider-mashery/mashschema"
)

type JsonIdentityFieldMapper[Ident any, MType any] struct {
	FieldMapperBase

	IdentityFunc      funcsupport.Supplier[Ident]
	Locator           LocatorFunc[MType, Ident]
	ValidateIdentFunc func(Ident) bool
}

func (sfm *JsonIdentityFieldMapper[Ident, MType]) PrepareMapper() *JsonIdentityFieldMapper[Ident, MType] {
	sfm.Schema.ValidateDiagFunc = sfm.ValidateDiag
	return sfm
}

func (sfm *JsonIdentityFieldMapper[Ident, MType]) ValueToSchema(i interface{}, state *schema.ResourceData) error {
	str := wrapJSON(i)
	return state.Set(sfm.Key, str)
}

func (sfm *JsonIdentityFieldMapper[Ident, MType]) ValidateDiag(i interface{}, _ cty.Path) diag.Diagnostics {
	rv := diag.Diagnostics{}

	if str, ok := i.(string); ok {
		ident := sfm.IdentityFunc()
		if err := unwrapJSON(str, &ident); err != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("supplied value is not a valid wrapped json"),
				Detail:   fmt.Sprintf("could not parse supplied value as type %s", reflect.TypeOf(ident).Name()),
			})
		} else if !sfm.ValidateIdentFunc(ident) {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("supplied identity is not valid"),
			})
		}
	} else {
		rv = append(rv, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("supplied value for field %s is not string", sfm.Key),
		})
	}

	return rv
}

func (sfm *JsonIdentityFieldMapper[Ident, MType]) RemoteToSchema(_ *MType, _ *schema.ResourceData) *diag.Diagnostic {
	// Nothing to do
	return nil
}

func (sfm *JsonIdentityFieldMapper[Ident, MType]) SchemaToRemote(state *schema.ResourceData, remote *MType) {
	if sfm.Schema.Computed && !sfm.Schema.Optional {
		return
	}

	ident := sfm.IdentityFunc()
	val := mashschema.ExtractString(state, sfm.Key, "")
	_ = unwrapJSON(val, &ident)

	*sfm.Locator(remote) = ident
}
