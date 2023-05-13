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
		Key: mashschema.MashPackageId,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Package Id, to which this plan belongs",
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
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key: mashschema.MashPlanId,
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key: mashschema.MashPlanStatus,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status of this plan",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.EmailTemplateSetId
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key: mashschema.MashPlanEmailTemplateSetId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "email template set id to use",
			},
		},
	})
}

func init() {
	PackagePlanResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key: mashschema.MashPlanName,
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
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key: mashschema.MashPlanDescription,
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key:    mashschema.MashPlanSelfServiceKeyProvisioningEnabled,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.AdminKeyProvisioningEnabled
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key:    mashschema.MashPlanAdminKeyProvisioningEnabled,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *string {
			return &in.Notes
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key:    mashschema.MashPlanNotes,
			Schema: &optionalString,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.QpsLimitExempt
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key:    mashschema.MashPlanQpsLimitKeyOverrideAllowed,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.RateLimitExempt
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key:    mashschema.MashPlanRateLimitExempt,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.RateLimitKeyOverrideAllowed
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key:    mashschema.MashPlanRateLimitKeyOverrideAllowed,
			Schema: &optionalBool,
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Plan]{
		Locator: func(in *masherytypes.Plan) *bool {
			return &in.ResponseFilterOverrideAllowed
		},
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
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
		FieldMapperBase: tfmapper.FieldMapperBase{
			Key: mashschema.MashPlanRateLimitCeiling,
			Schema: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5000,
			},
		},
	})
}
