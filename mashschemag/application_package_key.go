package mashschemag

import (
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var ApplicationPackageKeyResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.ApplicationIdentifier, masherytypes.ApplicationPackageKeyIdentifier, masherytypes.ApplicationPackageKey]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ApplicationPackageKeyIdentifier]{
		IdentityFunc: func() masherytypes.ApplicationPackageKeyIdentifier {
			return masherytypes.ApplicationPackageKeyIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ApplicationPackageKeyIdentifier) bool {
			return len(inp.PackageKeyId) > 0 && len(inp.ApplicationId) > 0
		},
	})

// Application parent identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.ApplicationIdentifier]{
		Key: mashschema.ApplicationRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Application reference, to which this application belongs",
			ForceNew:    true,
		},
		IdentityFunc: func() masherytypes.ApplicationIdentifier {
			return masherytypes.ApplicationIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.ApplicationIdentifier) bool {
			return len(inp.ApplicationId) > 0
		},
	}

	ApplicationPackageKeyResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 package key identifier",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.StringPtrFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) **string {
			return &in.Apikey
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyApiKey,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "API key",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.StringPtrFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) **string {
			return &in.Secret
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeySecret,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "API key secret",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.Int64PointerFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) **int64 {
			return &in.RateLimitCeiling
		},
		NilValue: int64(-1),
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyRateLimitCeiling,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Rate limit ceiling of this application",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.Int64PointerFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) **int64 {
			return &in.QpsLimitCeiling
		},
		NilValue: int64(-1),
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyQpsLimitCeiling,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "QPS limit ceiling of this application",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *bool {
			return &in.RateLimitExempt
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyRateLimitExempt,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Rate exempt",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *bool {
			return &in.QpsLimitExempt
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyQpsLimitExempt,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "QPS exempt",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *string {
			return &in.Status
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyStatus,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status of this package key (waiting or active)",
			},
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.PluggableFiledMapperBase[masherytypes.ApplicationPackageKey]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyLimits,
			Schema: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 2,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mashschema.ApplicationPackageKeyLimitPeriod: {
							Type:     schema.TypeString,
							Required: true,
						},
						mashschema.ApplicationPackageKeyLimitSource: {
							Type:     schema.TypeString,
							Required: true,
						},
						mashschema.ApplicationPackageKeyLimitCeiling: {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.ApplicationPackageKey, key string, state *schema.ResourceData) *diag.Diagnostic {
			var v []interface{}

			if remote.Limits != nil {
				if len(*remote.Limits) == 2 {
					for _, l := range *remote.Limits {
						mp := map[string]interface{}{
							mashschema.ApplicationPackageKeyLimitPeriod:  l.Period,
							mashschema.ApplicationPackageKeyLimitSource:  l.Source,
							mashschema.ApplicationPackageKeyLimitCeiling: l.Ceiling,
						}
						v = append(v, mp)
					}
				} else {
					return &diag.Diagnostic{
						Summary: "incomplete limits information",
						Detail:  fmt.Sprintf("need exactly 2 limits to perform mapping, but %d were supplied", len(*remote.Limits)),
					}

				}
			}

			return tfmapper.SetKeyWithDiag(state, key, v)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.ApplicationPackageKey) {
			if limitsArray, exists := state.GetOk(key); exists {
				tfLimits := mashschema.UnwrapStructArrayFromTerraformSet(limitsArray)

				var limits []masherytypes.Limit
				for _, tfl := range tfLimits {
					l := masherytypes.Limit{
						Period:  tfl[mashschema.ApplicationPackageKeyLimitPeriod].(string),
						Source:  tfl[mashschema.ApplicationPackageKeyLimitSource].(string),
						Ceiling: int64(tfl[mashschema.ApplicationPackageKeyLimitCeiling].(int)),
					}

					limits = append(limits, l)
				}

				// Limits are not necessary -- these need to be changed after the migration project.
				// remote.Limits = &limits
			}
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.PluggableFiledMapperBase[masherytypes.ApplicationPackageKey]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyPackagePlanRef,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Package plan reference of this package key",
			},
			ValidateFunc: func(in *schema.ResourceData, key string) (bool, string) {
				if in.HasChange(key) {
					p, c := in.GetChange(key)

					previous := p.(string)
					current := c.(string)
					if len(previous) > 0 {
						previousIdent := masherytypes.PackagePlanIdentifier{}
						currentIdent := masherytypes.PackagePlanIdentifier{}

						errP := tfmapper.UnwrapJSON(previous, &previousIdent)
						errC := tfmapper.UnwrapJSON(current, &currentIdent)

						if errP != nil {
							return false, fmt.Sprintf("previous package plan identity is malformed (%s)", previous)
						}
						if errC != nil {
							return false, fmt.Sprintf("updated package plan identity is malformed (%s)", current)
						}

						if errP != nil || errC != nil {
							return false, "package plan identity is malformed"
						}
						if previousIdent.PackageId != currentIdent.PackageId {
							return false, fmt.Sprintf(
								"moving package key from package %s to package % is not possible; taint this resource to have the package key deleted and re-created",
								previousIdent.PackageId,
								currentIdent.PackageId)
						}
					}
				}
				return true, ""
			},
		},
		RemoteToSchemaFunc: func(remote *masherytypes.ApplicationPackageKey, key string, state *schema.ResourceData) *diag.Diagnostic {
			if remote.Package != nil && remote.Plan != nil {
				ident := masherytypes.PackagePlanIdentifier{
					PackageIdentifier: masherytypes.PackageIdentifier{
						PackageId: remote.Package.Id,
					},
					PlanId: remote.Plan.Id,
				}
				return tfmapper.SetKeyWithDiag(state, key, tfmapper.WrapJSON(ident))
			} else {
				return nil
			}

		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.ApplicationPackageKey) {
			v := mashschema.ExtractString(state, key, "")
			ident := masherytypes.PackagePlanIdentifier{}

			if err := tfmapper.UnwrapJSON(v, &ident); err == nil {
				remote.Package = &masherytypes.Package{
					AddressableV3Object: masherytypes.AddressableV3Object{
						Id: ident.PackageId,
					},
				}
				remote.Plan = &masherytypes.Plan{
					AddressableV3Object: masherytypes.AddressableV3Object{
						Id: ident.PlanId,
					},
				}
			}
		},
	})
}

func init() {
	ApplicationPackageKeyResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.ApplicationPackageKey]{
		Locator: func(in *masherytypes.ApplicationPackageKey) *string {
			return &in.Expires
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.ApplicationPackageKey]{
			Key: mashschema.ApplicationPackageKeyExpires,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Expiry date of this key (if desired)",
				ValidateDiagFunc: mashschema.ValidateRegExp("\\d{4}-\\d{2}-\\d{2}"),
			},
		},
	})
}
