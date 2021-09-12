package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

const (
	MashMemberUserName       = "username"
	MashMemberUserNamePrefix = "username_prefix"
	MashMemberCreated        = "created"
	MashMemberUpdated        = "updated"
	MashMemberEmail          = "email"
	MashMemberDisplayName    = "display_name"
	MashMemberUri            = "uri"
	MashMemberBlog           = "blog"
	MashMemberIm             = "im"
	MashMemberImSvc          = "imsvc"
	MashMemberPhone          = "phone"
	MashMemberCompany        = "company"
	MashMemberAddress1       = "address_1"
	MashMemberAddress2       = "address_2"
	MashMemberLocality       = "locality"
	MashMemberRegion         = "region"
	MashMemberPostalCode     = "postal_code"
	MashMemberCountryCode    = "country_code"
	MashMemberFirstName      = "first_name"
	MashMemberLastName       = "last_name"
	MashMemberAreaStatus     = "area_status"
	MashMemberExternalId     = "external_id"
	MashMemberRoles          = "roles" // not implemented
)

var memberAreaStatusEnum = []string{"active", "waiting", "disabled"}

var MemberSchema = map[string]*schema.Schema{
	MashMemberUserName: {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "Member username name.",
	},
	MashMemberUserNamePrefix: {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{MashMemberUserName},
		Description:   "Prefix to generate a unique user name",
	},
	MashMemberEmail: {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Registration email",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			if str, ok := i.(string); ok {
				if match, err := regexp.MatchString(".+@.+\\..+", str); !match || err != nil {
					return diag.Diagnostics{diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "Malformed input",
						Detail:        "Input is not a well-formed email address",
						AttributePath: path,
					}}
				}
				return diag.Diagnostics{}
			} else {
				return diag.Errorf("input is not a string, but %s", i)
			}
		},
	},
	MashMemberAreaStatus: {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Area status",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return validateStringValueInSet(i, path, &memberAreaStatusEnum)
		},
	},
	// Roles are not implemented.
}

// Filling in member boilerplate settings
func fillMemberBoilerplate() {
	addComputedString(&MemberSchema, MashMemberCreated, "Date/time the member was created")
	addComputedString(&MemberSchema, MashMemberUpdated, "Date/time the member was updated")

	addRequiredString(&MemberSchema, MashMemberDisplayName, "Name to use in portal blog comments, discussion forums")
	addOptionalString(&MemberSchema, MashMemberUri, "URI of the website of the user")
	addOptionalString(&MemberSchema, MashMemberBlog, "URI of the blog of the user")
	addOptionalString(&MemberSchema, MashMemberIm, "IM handle")
	addOptionalString(&MemberSchema, MashMemberImSvc, "IM service, e.g. Hipchat, Google Hangouts, etc.")
	addOptionalString(&MemberSchema, MashMemberPhone, "Phone number of the user")
	addOptionalString(&MemberSchema, MashMemberCompany, "Company name")
	addOptionalString(&MemberSchema, MashMemberAddress1, "Address 1")
	addOptionalString(&MemberSchema, MashMemberAddress2, "Address 2")
	addOptionalString(&MemberSchema, MashMemberLocality, "Locale/city of the user")
	addOptionalString(&MemberSchema, MashMemberRegion, "Region/country of the user")
	addOptionalString(&MemberSchema, MashMemberPostalCode, "Postal/zip code")
	addOptionalString(&MemberSchema, MashMemberCountryCode, "Code of the country of the user")
	addOptionalString(&MemberSchema, MashMemberFirstName, "First name")
	addOptionalString(&MemberSchema, MashMemberLastName, "Last name")

	addOptionalString(&MemberSchema, MashMemberExternalId, "ID of the user in an external system, e.g. Salesforce")

}

func init() {
	fillMemberBoilerplate()
}

type MemberIdentifier struct {
	MemberId string
	Username string
}

func (m *MemberIdentifier) Id() string {
	return CreateCompoundId(m.MemberId, m.Username)
}

func (m *MemberIdentifier) From(id string) {
	ParseCompoundId(id, &m.MemberId, &m.Username)
}

func MashMemberUpsertable(d *schema.ResourceData) v3client.MasheryMember {
	mid := MemberIdentifier{}
	mid.From(d.Id())

	if mid.Username == "" {
		mid.Username = extractSetOrPrefixedString(d, MashMemberUserName, MashMemberUserNamePrefix)
	}

	return v3client.MasheryMember{
		AddressableV3Object: v3client.AddressableV3Object{
			Id: mid.MemberId,
		},
		Username:    mid.Username,
		Email:       extractString(d, MashMemberEmail, "not-set@terraform-managed.io"),
		DisplayName: extractString(d, MashMemberDisplayName, ""),
		Uri:         extractString(d, MashMemberUri, ""),
		Blog:        extractString(d, MashMemberBlog, ""),
		Im:          extractString(d, MashMemberIm, ""),
		Imsvc:       extractString(d, MashMemberImSvc, ""),
		Phone:       extractString(d, MashMemberPhone, ""),
		Company:     extractString(d, MashMemberCompany, ""),
		Address1:    extractString(d, MashMemberAddress1, ""),
		Address2:    extractString(d, MashMemberAddress2, ""),
		Locality:    extractString(d, MashMemberLocality, ""),
		Region:      extractString(d, MashMemberRegion, ""),
		PostalCode:  extractString(d, MashMemberPostalCode, ""),
		CountryCode: extractString(d, MashMemberCountryCode, ""),
		FirstName:   extractString(d, MashMemberFirstName, ""),
		LastName:    extractString(d, MashMemberLastName, ""),
		AreaStatus:  extractString(d, MashMemberAreaStatus, "waiting"),
		ExternalId:  extractString(d, MashMemberExternalId, ""),
	}
}

func V3MemberToResourceData(inp *v3client.MasheryMember, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashMemberUserName: inp.Username,
		MashMemberCreated:  inp.Created.ToString(),
		MashMemberUpdated:  inp.Updated.ToString(),

		// Email and display name will not be updated; these are mandatory fields.
		//MashMemberEmail: inp.Email,
		//MashMemberDisplayName: inp.DisplayName,

		MashMemberUri:         inp.Uri,
		MashMemberBlog:        inp.Blog,
		MashMemberIm:          inp.Im,
		MashMemberImSvc:       inp.Imsvc,
		MashMemberPhone:       inp.Phone,
		MashMemberCompany:     inp.Company,
		MashMemberAddress1:    inp.Address1,
		MashMemberAddress2:    inp.Address2,
		MashMemberLocality:    inp.Locality,
		MashMemberRegion:      inp.Region,
		MashMemberPostalCode:  inp.PostalCode,
		MashMemberCountryCode: inp.CountryCode,
		MashMemberFirstName:   inp.FirstName,
		MashMemberLastName:    inp.LastName,
		MashMemberAreaStatus:  inp.AreaStatus,
		MashMemberExternalId:  inp.ExternalId,
	}

	return SetResourceFields(data, d)
}
