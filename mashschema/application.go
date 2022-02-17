package mashschema

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MashAppId         = "application_id"
	MashAppName       = "name"
	MashAppNamePrefix = "name_prefix"
	MashAppCreated    = "created"
	MashAppUpdated    = "updated"

	MashAppOwner             = "owner"
	MashAppOwnerUsername     = "owner_username"
	MashAppDescription       = "description"
	MashAppType              = "type"
	MashAppCommercial        = "commercial"
	MashAppAds               = "ads"
	MashAppAdSystem          = "ads_system"
	MashAppUsageModel        = "usage_model"
	MashAppTags              = "tags"
	MashAppNotes             = "notes"
	MashAppHowDidYouHear     = "how_did_you_hear"
	MashAppPreferredProtocol = "preferred_protocol"
	MashAppPreferredOutput   = "preferred_output"
	MashAppExternalId        = "external_id"
	MashAppUri               = "uri"
	MashAppOAuthRedirectUri  = "oauth_redirect_uri"
	MashAppEAV               = "eav"
)

var ApplicationMapper *ApplicationMapperImpl

type ApplicationMapperImpl struct {
	MapperImpl
}

func (am *ApplicationMapperImpl) CreateIdentifierTyped() *ApplicationIdentifier {
	return &ApplicationIdentifier{}
}

type ApplicationIdentifier struct {
	MemberIdentifier
	AppId string `json:"appId"`
}

func (ai *ApplicationIdentifier) Self() interface{} {
	return ai
}

func (ai *ApplicationMapperImpl) GetIdentifier(d *schema.ResourceData) *ApplicationIdentifier {
	rv := &ApplicationIdentifier{}
	CompoundIdFrom(rv, d.Id())

	return rv
}

func (ai *ApplicationMapperImpl) GetOwnerIdentifier(d *schema.ResourceData) *MemberIdentifier {
	rv := &MemberIdentifier{}
	CompoundIdFrom(rv, ExtractString(d, MashAppOwner, ""))

	return rv
}

func (ai *ApplicationMapperImpl) UpsertableTyped(_ context.Context, d *schema.ResourceData) (masherytypes.MasheryApplication, diag.Diagnostics) {
	mid := MemberIdentifier{}
	CompoundIdFrom(&mid, ExtractString(d, MashAppOwner, ""))

	appId := ApplicationIdentifier{}
	CompoundIdFrom(&appId, d.Id())

	return masherytypes.MasheryApplication{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   appId.AppId,
			Name: extractSetOrPrefixedString(d, MashAppName, MashAppNamePrefix),
		},
		Username:          mid.Username,
		Description:       ExtractString(d, MashAppDescription, ""),
		Type:              ExtractString(d, MashAppType, ""),
		Commercial:        extractBool(d, MashAppCommercial, false),
		Ads:               extractBool(d, MashAppAds, false),
		AdsSystem:         ExtractString(d, MashAppAdSystem, ""),
		UsageModel:        ExtractString(d, MashAppUsageModel, ""),
		Tags:              ExtractString(d, MashAppTags, ""),
		Notes:             ExtractString(d, MashAppNotes, ""),
		HowDidYouHear:     ExtractString(d, MashAppHowDidYouHear, ""),
		PreferredProtocol: ExtractString(d, MashAppPreferredProtocol, ""),
		PreferredOutput:   ExtractString(d, MashAppPreferredOutput, ""),
		ExternalId:        ExtractString(d, MashAppExternalId, ""),
		Uri:               ExtractString(d, MashAppUri, ""),
		OAuthRedirectUri:  ExtractString(d, MashAppOAuthRedirectUri, ""),
		Eav:               extractEAVPointer(d, MashAppEAV),
	}, nil
}

func (ai *ApplicationMapperImpl) PersistTyped(ctx context.Context, rawInp interface{}, d *schema.ResourceData) diag.Diagnostics {
	inp := rawInp.(*masherytypes.MasheryApplication)
	data := map[string]interface{}{
		MashAppId:                inp.Id,
		MashAppName:              inp.Name,
		MashAppCreated:           inp.Created.ToString(),
		MashAppUpdated:           inp.Updated.ToString(),
		MashAppOwnerUsername:     inp.Username,
		MashAppDescription:       inp.Description,
		MashAppType:              inp.Type,
		MashAppCommercial:        inp.Commercial,
		MashAppAds:               inp.Ads,
		MashAppAdSystem:          inp.AdsSystem,
		MashAppUsageModel:        inp.UsageModel,
		MashAppTags:              inp.Tags,
		MashAppNotes:             inp.Notes,
		MashAppHowDidYouHear:     inp.HowDidYouHear,
		MashAppPreferredProtocol: inp.PreferredProtocol,
		MashAppPreferredOutput:   inp.PreferredOutput,
		MashAppExternalId:        inp.ExternalId,
		MashAppUri:               inp.Uri,
		MashAppOAuthRedirectUri:  inp.OAuthRedirectUri,
		MashAppEAV:               inp.Eav,
	}

	return ai.SetResourceFields(ctx, data, d)
}

// Fill in boilerplate fields of Mashery application.
func fillAppSchemaBoilerplate() {
	addComputedString(&ApplicationMapper.schema, MashAppId, "Original application id")
	addComputedString(&ApplicationMapper.schema, MashAppCreated, "Date/time the object was created")
	addComputedString(&ApplicationMapper.schema, MashAppUpdated, "Date/time the object was last updated")
	addComputedString(&ApplicationMapper.schema, MashAppOwnerUsername, "Owner user name of this application")

	addOptionalString(&ApplicationMapper.schema, MashAppDescription, "Description of the application")
	addOptionalString(&ApplicationMapper.schema, MashAppType, "Type of application")

	addOptionalBoolean(&ApplicationMapper.schema, MashAppCommercial, "Whether or not the application is commercial in nature")
	addOptionalBoolean(&ApplicationMapper.schema, MashAppAds, "Whether or not the application supports ads")

	addOptionalString(&ApplicationMapper.schema, MashAppAdSystem, "Advertisement system")
	addOptionalString(&ApplicationMapper.schema, MashAppUsageModel, "Usage model")

	addOptionalString(&ApplicationMapper.schema, MashAppTags, "Tags, i.e. tracking metadata")
	addOptionalString(&ApplicationMapper.schema, MashAppNotes, "Notes about the application.")
	addOptionalString(&ApplicationMapper.schema, MashAppHowDidYouHear, "How did someone hear about the API?")
	addOptionalString(&ApplicationMapper.schema, MashAppPreferredProtocol, "Protocol preference of developer, e.g. REST or SOAP")
	addOptionalString(&ApplicationMapper.schema, MashAppPreferredOutput, "Output preference of developer, e.g. json or xml.")
	addOptionalString(&ApplicationMapper.schema, MashAppExternalId, "ID of the application in an external system, e.g. Salesforce")
	addOptionalString(&ApplicationMapper.schema, MashAppUri, "URI of the application")
	addOptionalString(&ApplicationMapper.schema, MashAppOAuthRedirectUri, "OAuth 2 redirect URI")
}

func init() {
	ApplicationMapper = &ApplicationMapperImpl{
		MapperImpl: MapperImpl{
			schema: map[string]*schema.Schema{
				MashAppName: {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Application name",
				},
				MashAppNamePrefix: {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{MashAppName},
					Description:   "Prefix for the application names",
				},
				MashAppOwner: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Username of the member that the application belongs to",
					ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
						return ValidateCompoundIdent(i, path, func() interface{} {
							return &MemberIdentifier{}
						})
					},
				},
				MashAppEAV: {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			}, // Schema
		},
	}

	fillAppSchemaBoilerplate()

	ApplicationMapper.upsertFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return ApplicationMapper.UpsertableTyped(ctx, d)
	}

	ApplicationMapper.persistFunc = func(ctx context.Context, rv interface{}, d *schema.ResourceData) diag.Diagnostics {
		return ApplicationMapper.PersistTyped(ctx, rv.(*masherytypes.MasheryApplication), d)
	}

	ApplicationMapper.identifier = func() interface{} {
		return ApplicationMapper.CreateIdentifierTyped()
	}
}
