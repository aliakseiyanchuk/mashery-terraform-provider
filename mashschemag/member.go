package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

var MemberResourceSchemaBuilder = tfmapper.NewSchemaBuilder[tfmapper.Orphan, masherytypes.MemberIdentifier, masherytypes.Member]().
	Identity(&tfmapper.JsonIdentityMapper[masherytypes.MemberIdentifier]{
		IdentityFunc: func() masherytypes.MemberIdentifier { return masherytypes.MemberIdentifier{} },
		ValidateIdentFunc: func(inp masherytypes.MemberIdentifier) bool {
			return len(inp.MemberId) > 0 && len(inp.Username) > 0
		},
	}).
	RootIdentity(&tfmapper.RootParentIdentity{})

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 member identifier",
			},
		},
	})
}

// -----------------
// Addressable object properties

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 member name",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MashPackUpdated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was updated",
			},
		},
	})
}

// -----------------------------
// Object-specific fields

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Username
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberUserName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mashery V3 member user name",
			},
		},
	})
}
func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Email
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberUserEmail,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Mashery V3 member user email address",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.DisplayName
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberDisplayName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member display name",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Uri
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.ObjectUri,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member URI",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Blog
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberBlog,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member blog",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Im
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberIM,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member instant messenger contact",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Imsvc
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberIMService,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member instant messenger service",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Phone
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberPhone,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member phone",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Company
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberCompany,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member company",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Address1
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberAddress1,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member address 1",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Address2
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberAddress2,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member address 2",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Locality
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberLocality,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member locality",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.Region
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberRegion,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member region",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.PostalCode
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberPostalCode,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member postal code",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.CountryCode
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberCountryCode,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member country code",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.FirstName
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberFirstName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member first name",
			},
		},
	}).Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.LastName
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberLastName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member last name",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.AreaStatus
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberAreaStatus,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 member area status",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) *string {
			return &in.ExternalId
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.ObjectExternalId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mashery V3 member external identified",
			},
		},
	})
}

func init() {
	MemberResourceSchemaBuilder.Add(&tfmapper.StringPtrFieldMapper[masherytypes.Member]{
		Locator: func(in *masherytypes.Member) **string {
			return &in.PasswdNew
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Member]{
			Key: mashschema.MemberInitialPassword,
			Schema: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return true
				},
				DiffSuppressOnRefresh: true,
				Description:           "Mashery V3 member external identified",
			},
		},
	})
}
