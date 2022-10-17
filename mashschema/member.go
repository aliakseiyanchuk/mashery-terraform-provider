package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
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

// Filling in member boilerplate settings
func fillMemberBoilerplate() {
	addComputedString(&MemberMapper.schema, MashMemberCreated, "Date/time the member was created")
	addComputedString(&MemberMapper.schema, MashMemberUpdated, "Date/time the member was updated")

	addRequiredString(&MemberMapper.schema, MashMemberDisplayName, "Name to use in portal blog comments, discussion forums")
	addOptionalString(&MemberMapper.schema, MashMemberUri, "URI of the website of the user")
	addOptionalString(&MemberMapper.schema, MashMemberBlog, "URI of the blog of the user")
	addOptionalString(&MemberMapper.schema, MashMemberIm, "IM handle")
	addOptionalString(&MemberMapper.schema, MashMemberImSvc, "IM service, e.g. Hipchat, Google Hangouts, etc.")
	addOptionalString(&MemberMapper.schema, MashMemberPhone, "Phone number of the user")
	addOptionalString(&MemberMapper.schema, MashMemberCompany, "Company name")
	addOptionalString(&MemberMapper.schema, MashMemberAddress1, "Address 1")
	addOptionalString(&MemberMapper.schema, MashMemberAddress2, "Address 2")
	addOptionalString(&MemberMapper.schema, MashMemberLocality, "Locale/city of the user")
	addOptionalString(&MemberMapper.schema, MashMemberRegion, "Region/country of the user")
	addOptionalString(&MemberMapper.schema, MashMemberPostalCode, "Postal/zip code")
	addOptionalString(&MemberMapper.schema, MashMemberCountryCode, "Code of the country of the user")
	addOptionalString(&MemberMapper.schema, MashMemberFirstName, "First name")
	addOptionalString(&MemberMapper.schema, MashMemberLastName, "Last name")

	addOptionalString(&MemberMapper.schema, MashMemberExternalId, "ID of the user in an external system, e.g. Salesforce")

}

var MemberMapper *MemberMapperImpl

type MemberMapperImpl struct {
	ResourceMapperImpl
}

func (mmi *MemberMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.Member, V3ObjectIdentifier, diag.Diagnostics) {
	mid := masherytypes.MemberIdentifier{}
	CompoundIdFrom(&mid, d.Id())

	if mid.Username == "" {
		mid.Username = extractSetOrPrefixedString(d, MashMemberUserName, MashMemberUserNamePrefix)
	}

	return masherytypes.Member{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id: mid.MemberId,
		},
		Username:    mid.Username,
		Email:       ExtractString(d, MashMemberEmail, "not-set@terraform-managed.io"),
		DisplayName: ExtractString(d, MashMemberDisplayName, ""),
		Uri:         ExtractString(d, MashMemberUri, ""),
		Blog:        ExtractString(d, MashMemberBlog, ""),
		Im:          ExtractString(d, MashMemberIm, ""),
		Imsvc:       ExtractString(d, MashMemberImSvc, ""),
		Phone:       ExtractString(d, MashMemberPhone, ""),
		Company:     ExtractString(d, MashMemberCompany, ""),
		Address1:    ExtractString(d, MashMemberAddress1, ""),
		Address2:    ExtractString(d, MashMemberAddress2, ""),
		Locality:    ExtractString(d, MashMemberLocality, ""),
		Region:      ExtractString(d, MashMemberRegion, ""),
		PostalCode:  ExtractString(d, MashMemberPostalCode, ""),
		CountryCode: ExtractString(d, MashMemberCountryCode, ""),
		FirstName:   ExtractString(d, MashMemberFirstName, ""),
		LastName:    ExtractString(d, MashMemberLastName, ""),
		AreaStatus:  ExtractString(d, MashMemberAreaStatus, "waiting"),
		ExternalId:  ExtractString(d, MashMemberExternalId, ""),
	}, nil, nil
}

func (mmi *MemberMapperImpl) PersistTyped(inp masherytypes.Member, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashMemberUserName: inp.Username,
		MashMemberCreated:  inp.Created.ToString(),
		MashMemberUpdated:  inp.Updated.ToString(),

		// Email and display name will not be updated; these are mandatory fields.
		MashMemberEmail:       inp.Email,
		MashMemberDisplayName: inp.DisplayName,

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

	return mmi.persistMap(inp.Identifier(), data, d)
}

func init() {
	MemberMapper = &MemberMapperImpl{
		ResourceMapperImpl{
			v3ObjectName: "member",
			schema: map[string]*schema.Schema{
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
			},
		},
	}

	fillMemberBoilerplate()

	MemberMapper.v3Identity = func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		rv := masherytypes.MemberIdentifier{}
		if CompoundIdFrom(&rv, d.Id()) {
			return rv, nil
		} else {
			return rv, diag.Diagnostics{MemberMapper.lackingIdentificationDiagnostic("id")}
		}
	}
	MemberMapper.upsertFunc = func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
		return MemberMapper.UpsertableTyped(d)
	}

	MemberMapper.persistFunc = func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
		ptr := rv.(*masherytypes.Member)
		return MemberMapper.PersistTyped(*ptr, d)
	}
}
