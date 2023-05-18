package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var ServiceResourceSchemaBuilder = tfmapper.NewSchemaBuilder[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.ServiceIdentifier]{
		IdentityFunc: func() masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{}
		},
	}).
	RootIdentity(&tfmapper.RootParentIdentity{})

func init() {
	ServiceResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 service identifier of this service",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	}).Add(&tfmapper.SerOrPrefixedFieldMapper[masherytypes.Service]{
		StringFieldMapper: tfmapper.StringFieldMapper[masherytypes.Service]{
			Locator: func(in *masherytypes.Service) *string {
				return &in.Name
			},

			FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
				Key: mashschema.MashObjName,
			},
		},

		PrefixKey: mashschema.MashObjNamePrefix,
		CompositeMapperBase: tfmapper.CompositeMapperBase{
			CompositeSchema: map[string]*schema.Schema{
				mashschema.MashObjName: {
					Type:          schema.TypeString,
					Optional:      true,
					Computed:      true,
					Description:   "Service name",
					ConflictsWith: []string{mashschema.MashObjNamePrefix},
				},
				mashschema.MashObjNamePrefix: {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Prefix for the service name",
					ConflictsWith: []string{mashschema.MashObjName},
				},
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *string {
			return &in.EditorHandle
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcEditorHandle,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User id which perform latest modification",
			},
		},
	}).Add(&tfmapper.IntFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *int {
			return &in.RevisionNumber
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcRevisionNumber,
			Schema: &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Count of updates that were applied on this service after update",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *string {
			return &in.RobotsPolicy
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcRobotsPolicy,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Robots policy",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *string {
			return &in.CrossdomainPolicy
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcCrossdomainPolicy,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cross-domain policy",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *string {
			return &in.Description
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcDescription,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Managed by Terraform",
				Description: "Description of this service",
			},
		},
	}).Add(&tfmapper.Int64PointerFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) **int64 {
			return &in.QpsLimitOverall
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcQpsLimitOverall,
			Schema: &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: mashschema.ValidateZeroOrGreater,
				Description:      "Maximum number of calls handled per second (QPS) across all developer keys for the API. Most customers do not set a value for this particular setting.",
			},
		},
	}).Add(&tfmapper.BoolFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *bool {
			return &in.RFC3986Encode
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcServiceRFC3986Encode,
			Schema: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Service]{
		Locator: func(in *masherytypes.Service) *string {
			return &in.Version
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Service]{
			Key: mashschema.MashSvcVersion,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0.0.1/TF",
				Description: "Deployed-defined version designator",
			},
		},
	})
}