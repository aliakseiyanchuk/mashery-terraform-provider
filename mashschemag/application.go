package mashschemag

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashschema"
	"terraform-provider-mashery/tfmapper"
)

type ApplicationOfMemberIdentifier struct {
	masherytypes.MemberIdentifier
	masherytypes.ApplicationIdentifier
}

var ApplicationResourceSchemaBuilder = tfmapper.NewSchemaBuilder[masherytypes.MemberIdentifier, ApplicationOfMemberIdentifier, masherytypes.Application]().
	Identity(&tfmapper.JsonIdentityMapper[ApplicationOfMemberIdentifier]{
		IdentityFunc: func() ApplicationOfMemberIdentifier { return ApplicationOfMemberIdentifier{} },
	})

// Application parent identity
func init() {
	mapper := tfmapper.JsonIdentityMapper[masherytypes.MemberIdentifier]{
		Key: mashschema.MemberRef,
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "User reference, to which this application belongs",
		},
		IdentityFunc: func() masherytypes.MemberIdentifier {
			return masherytypes.MemberIdentifier{}
		},
		ValidateIdentFunc: func(inp masherytypes.MemberIdentifier) bool {
			return len(inp.MemberId) > 0 && len(inp.Username) > 0
		},
	}

	ApplicationResourceSchemaBuilder.ParentIdentity(mapper.PrepareParentMapper())
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Id
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 application identifier",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Name
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.MashObjName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mashery V3 application name",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.DateMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *masherytypes.MasheryJSONTime {
			return in.Created
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.MashPackCreated,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date/time the object was created",
			},
		},
	}).Add(&tfmapper.DateMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *masherytypes.MasheryJSONTime {
			return in.Updated
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
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
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Username
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationUserName,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application owner's user name",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Description
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.MashObjDescription,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application description",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Type
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationType,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application type",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *bool {
			return &in.Commercial
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationCommercial,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application commercial flag",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.BoolFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *bool {
			return &in.Ads
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationAds,
			Schema: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application ads flag",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.AdsSystem
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationAdsSystem,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application ads system",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.UsageModel
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationUsageModel,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application usage model",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Tags
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationTags,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application tags",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Notes
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationNotes,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application notes",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.HowDidYouHear
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationHowDidYouHear,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application how-did-you-hear response",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.PreferredProtocol
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationPreferredProtocol,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application preferred protocol (Protocol preference of developer, e.g. REST or SOAP.)",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.PreferredOutput
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationPreferredOutput,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application preferred output",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.ExternalId
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ObjectExternalId,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application external identifier",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.Uri
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ObjectUri,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application URI",
			},
		},
	})
}

func init() {
	ApplicationResourceSchemaBuilder.Add(&tfmapper.StringFieldMapper[masherytypes.Application]{
		Locator: func(in *masherytypes.Application) *string {
			return &in.OAuthRedirectUri
		},
		FieldMapperBase: tfmapper.FieldMapperBase[masherytypes.Application]{
			Key: mashschema.ApplicationOAuthRedirectURI,
			Schema: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mashery V3 application OAuth redirect URI",
			},
		},
	})
}
