package mashery

import (
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

var AppSchema = map[string]*schema.Schema{
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
		Description: "Username of the membe that the application belongs to",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return ValidateCompoundIdent(i, path, 2)
		},
	},
	MashAppEAV: {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}

type ApplicationIdentifier struct {
	MemberId string
	Username string
	AppId    string
}

func (a *ApplicationIdentifier) Id() string {
	return CreateCompoundId(a.MemberId, a.Username, a.AppId)
}

func (a *ApplicationIdentifier) From(id string) {
	ParseCompoundId(id, &a.MemberId, &a.Username, &a.AppId)
}

func (a *ApplicationIdentifier) IsIdentified() bool {
	return len(a.MemberId) > 0 && len(a.Username) > 0 && len(a.AppId) > 0
}

func MashAppUpsertable(d *schema.ResourceData) masherytypes.MasheryApplication {
	mid := MemberIdentifier{}
	mid.From(extractString(d, MashAppOwner, ""))

	appId := ApplicationIdentifier{}
	appId.From(d.Id())

	return masherytypes.MasheryApplication{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   appId.AppId,
			Name: extractSetOrPrefixedString(d, MashAppName, MashAppNamePrefix),
		},
		Username:          mid.Username,
		Description:       extractString(d, MashAppDescription, ""),
		Type:              extractString(d, MashAppType, ""),
		Commercial:        extractBool(d, MashAppCommercial, false),
		Ads:               extractBool(d, MashAppAds, false),
		AdsSystem:         extractString(d, MashAppAdSystem, ""),
		UsageModel:        extractString(d, MashAppUsageModel, ""),
		Tags:              extractString(d, MashAppTags, ""),
		Notes:             extractString(d, MashAppNotes, ""),
		HowDidYouHear:     extractString(d, MashAppHowDidYouHear, ""),
		PreferredProtocol: extractString(d, MashAppPreferredProtocol, ""),
		PreferredOutput:   extractString(d, MashAppPreferredOutput, ""),
		ExternalId:        extractString(d, MashAppExternalId, ""),
		Uri:               extractString(d, MashAppUri, ""),
		OAuthRedirectUri:  extractString(d, MashAppOAuthRedirectUri, ""),
		Eav:               extractEAVPointer(d, MashAppEAV),
	}
}

func V3AppToResourceData(inp *masherytypes.MasheryApplication, d *schema.ResourceData) diag.Diagnostics {
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

	return SetResourceFields(data, d)
}

// Fill in boilerplate fields of Mashery application.
func fillAppSchemaBoilerplate() {
	addComputedString(&AppSchema, MashAppId, "Original application id")
	addComputedString(&AppSchema, MashAppCreated, "Date/time the object was created")
	addComputedString(&AppSchema, MashAppUpdated, "Date/time the object was last updated")
	addComputedString(&AppSchema, MashAppOwnerUsername, "Owner user name of this application")

	addOptionalString(&AppSchema, MashAppDescription, "Description of the application")
	addOptionalString(&AppSchema, MashAppType, "Type of application")

	addOptionalBoolean(&AppSchema, MashAppCommercial, "Whether or not the application is commercial in nature")
	addOptionalBoolean(&AppSchema, MashAppAds, "Whether or not the application supports ads")

	addOptionalString(&AppSchema, MashAppAdSystem, "Advertisement system")
	addOptionalString(&AppSchema, MashAppUsageModel, "Usage model")

	addOptionalString(&AppSchema, MashAppTags, "Tags, i.e. tracking metadata")
	addOptionalString(&AppSchema, MashAppNotes, "Notes about the application.")
	addOptionalString(&AppSchema, MashAppHowDidYouHear, "How did someone hear about the API?")
	addOptionalString(&AppSchema, MashAppPreferredProtocol, "Protocol preference of developer, e.g. REST or SOAP")
	addOptionalString(&AppSchema, MashAppPreferredOutput, "Output preference of developer, e.g. json or xml.")
	addOptionalString(&AppSchema, MashAppExternalId, "ID of the application in an external system, e.g. Salesforce")
	addOptionalString(&AppSchema, MashAppUri, "URI of the application")
	addOptionalString(&AppSchema, MashAppOAuthRedirectUri, "OAuth 2 redirect URI")
}

func init() {
	fillAppSchemaBoilerplate()
}
