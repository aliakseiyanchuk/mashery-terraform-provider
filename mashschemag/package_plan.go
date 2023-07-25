package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var PackagePlanResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.PackageIdentifier, masherytypes.PackagePlanIdentifier, masherytypes.Plan]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.PackagePlanIdentifier]{
		IdentityFunc: func() masherytypes.PackagePlanIdentifier {
			return masherytypes.PackagePlanIdentifier{}
		},
	})

// Parent package identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.PackageIdentifier]{
		Key: mashschema.MashPackageRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Package reference, to which this plan belongs",
		},
		IdentityFunc: func() masherytypes.PackageIdentifier {
			return masherytypes.PackageIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.PackageIdentifier) bool {
			return len(inp.PackageId) > 0
		},
	}

	PackagePlanResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

// init Plan id; Created and Updated
func init() {
	PackagePlanResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPackagePlanId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan id",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Status
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanStatus,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "active",
				Description: "status of this plan",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.MashPlanStatusEnum)
				},
			},
		},
	}).Add(&tfmapper.StringPtrFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) **string {
			return &in.EmailTemplateSetId
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanDeveloperEmailTemplateSetId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "email template set id to use",
			},
		},
	}).Add(&tfmapper.StringPtrFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) **string {
			return &in.AdminEmailTemplateSetId
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanAdminEmailTemplateSetId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "admin email template set id to use",
			},
		},
	})
}

func init() {
	PackagePlanResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashObjName,
			Schema: &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Plan name",
				ValidateDiagFunc: mashschema.ValidateNonEmptyString,
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Description
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashObjDescription,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
		},
	}).Add(&tfmapper.EAVFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) **masherytypes.EAV {
			return &in.Eav
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanEAV,
			Schema: &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Extended attribute values",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}
func init() {
	optionalBool := schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: false,
	}

	optionalString := schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: false,
	}

	PackagePlanResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.SelfServiceKeyProvisioningEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanSelfServiceKeyProvisioningEnabled,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.AdminKeyProvisioningEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanAdminKeyProvisioningEnabled,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Notes
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanNotes,
			Schema: &optionalString,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.QpsLimitExempt
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanQpsLimitExempt,
			Schema: &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{mashschema.MashPlanQpsLimitCeiling},
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.QpsLimitKeyOverrideAllowed
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanQpsLimitKeyOverrideAllowed,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.RateLimitExempt
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanRateLimitExempt,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.RateLimitKeyOverrideAllowed
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanRateLimitKeyOverrideAllowed,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.ResponseFilterOverrideAllowed
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key:    mashschema.MashPlanResponseFilterOverrideAllowed,
			Schema: &optionalBool,
		},
	})
}

func init() {
	PackagePlanResourceSchemaBuilder.Add(&tfmapper.IntFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *int {
			return &in.MaxNumKeysAllowed
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanMaxNumKeysAllowed,
			Schema: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},
	}).Add(&tfmapper.IntFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *int {
			return &in.NumKeysBeforeReview
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanNumKeysBeforeReview,
			Schema: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.RateLimitPeriod
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanRateLimitPeriod,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					return mashschema.ValidateStringValueInSet(i, path, &mashschema.RateLimitPeriodEnum)
				},
			},
		},
	})
}

// Int64 pointer fields
func init() {
	PackagePlanResourceSchemaBuilder.Add(&tfmapper.Int64PointerFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) **int64 {
			return &in.QpsLimitCeiling
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanQpsLimitCeiling,
			Schema: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},
	}).Add(&tfmapper.Int64PointerFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) **int64 {
			return &in.RateLimitCeiling
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanRateLimitCeiling,
			Schema: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000,
			},
		},
	})
}

// Init portal access roles
func init() {
	// TODO: This looks like a common implementation between services IO docs and plans access control
	// to receive keys. This needs to be refactored into a separate mapper to follow the DRY principle
	PackagePlanResourceSchemaBuilder.Add(&tfmapper.PluggableFiledMapperBase[masherytypes.Plan]{
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Plan]{
			Key: mashschema.MashPlanPortalAccessRoles,
			Schema: &schema.Schema{
				Optional: true,
				Type:     schema.TypeSet,
				Set:      mashschema.StringHashcode,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		NilRemoteToSchemaFunc: func(key string, state *schema.ResourceData) *diag.Diagnostic {
			var emptyArray []string
			return tfmapper.SetKeyWithDiag(state, key, emptyArray)
		},
		RemoteToSchemaFunc: func(remote *masherytypes.Plan, key string, state *schema.ResourceData) *diag.Diagnostic {
			var values []string

			if remote.Roles != nil {
				values = make([]string, len(*remote.Roles))

				for i, v := range *remote.Roles {
					values[i] = v.Id
				}
			}

			return tfmapper.SetKeyWithDiag(state, key, values)
		},
		SchemaToRemoteFunc: func(state *schema.ResourceData, key string, remote *masherytypes.Plan) {
			arr := mashschema.ExtractStringArray(state, key, &[]string{})

			if len(arr) > 0 {
				rolesArr := make([]masherytypes.RolePermission, len(arr))
				for i, v := range arr {
					perm := masherytypes.RolePermission{}
					perm.Id = v
					perm.Action = "register_keys"

					rolesArr[i] = perm
				}

				remote.Roles = &rolesArr
			} else {
				remote.Roles = &[]masherytypes.RolePermission{}
			}
		},
	})
}
