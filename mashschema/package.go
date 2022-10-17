package mashschema

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform/helper/hashcode"
	"strings"
)

const (
	MashPackagekId                      = "package_id"
	MashPackName                        = "name"
	MashPackNamePrefix                  = "name_prefix"
	MashPackCreated                     = "created"
	MashPackUpdated                     = "updated"
	MashPackDescription                 = "description"
	MashPackTags                        = "tags"
	MashPackNotifyDeveloperPeriod       = "notify_developer_period"
	MashPackNotifyDeveloperNearQuota    = "notify_developer_near_quota"
	MashPackNotifyDeveloperOverQuota    = "notify_developer_over_quota"
	MashPackNotifyDeveloperOverThrottle = "notify_developer_over_throttle"
	MashPackNotifyAdminPeriod           = "notify_admin_period"
	MashPackNotifyAdminNearQuota        = "notify_admin_near_quota"
	MashPackNotifyAdminOverQuota        = "notify_admin_over_quota"
	MashPackNotifyAdminOverThrottle     = "notify_admin_over_throttle"
	MashPackNotifyAdminEmails           = "notify_admin_emails"
	MashPackNearQuotaThreshold          = "near_quota_threshold"
	MashPackEAVs                        = "extended_attribute_values"
	MashPackKeyAdapter                  = "key_adapter"
	MashPackKeyLength                   = "key_length"
	MashPackSharedSecretLength          = "shared_secret_length"

	// Developer notifications
	// TODO: MOve to shared go file.
	MashDurationMinute = "minute"
	MashDurationHourly = "hour"
	MashDurationDay    = "day"
	MashDurationWeek   = "week"
	MashDurationMonth  = "month"
)

var notifyDeveloperPeriodEnum = []string{MashDurationMinute, MashDurationHourly,
	MashDurationDay, MashDurationWeek, MashDurationMonth}
var notifyAdminPeriodEnum = []string{MashDurationMinute, MashDurationHourly,
	MashDurationDay, MashDurationWeek, MashDurationMonth}

var PackageMapper *PackageMapperImpl

type PackageMapperImpl struct {
	ResourceMapperImpl
}

var PackageSchema = map[string]*schema.Schema{
	MashPackagekId: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Package Id",
	},
	MashPackCreated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date/time the object was created",
	},
	MashPackUpdated: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date/time the object was updated",
	},
	MashPackName: {
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		Description:   "Package name",
		ConflictsWith: []string{MashPackNamePrefix},
	},
	MashPackNamePrefix: {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "Prefix for the package name",
		ConflictsWith: []string{MashPackName},
	},
	MashPackDescription: {
		Type:          schema.TypeString,
		Optional:      true,
		Default:       "Managed by Terraform",
		Description:   "Package description",
		ConflictsWith: []string{MashPackTags},
	},
	MashPackTags: {
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "Tags to associate with this endpoint",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		ConflictsWith: []string{MashPackDescription},
	},
	MashPackNotifyDeveloperPeriod: {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "hour",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return validateStringValueInSet(i, path, &notifyDeveloperPeriodEnum)
		},
	},
	MashPackNotifyDeveloperNearQuota: {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
		//Default:     true,
		Description: "Notify developer when approaching quota",
	},
	MashPackNotifyDeveloperOverQuota: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Notify developer when quota exceeded",
	},
	MashPackNotifyDeveloperOverThrottle: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Notify developer when throttle exceeded",
	},
	MashPackNotifyAdminPeriod: {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		//Default:     "day",
		Description: "Package admin notification",
		ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
			return validateStringValueInSet(i, path, &notifyAdminPeriodEnum)
		},
	},
	MashPackNotifyAdminNearQuota: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Notify admin when approaching quota",
	},
	MashPackNotifyAdminOverQuota: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Notify admin when quota exceeded",
	},
	MashPackNotifyAdminOverThrottle: {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Notify admin when throttle exceeded",
	},
	MashPackNotifyAdminEmails: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Email addresses to send admin notifications",
		// TODO: Reference to the string element set is repeated
		// It could be placed in a shared module.
		Set: func(i interface{}) int {
			return hashcode.String(i.(string))
		},
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashPackNearQuotaThreshold: {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Percentage of quota when approaching limit notifications will be sent",
	},
	MashPackEAVs: {
		Type:        schema.TypeMap,
		Optional:    true,
		Computed:    true,
		Description: "Extended attribute values",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	MashPackKeyAdapter: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Custom adapter for key generation",
	},
	MashPackKeyLength: {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Length of keys for this package",
	},
	MashPackSharedSecretLength: {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Length of shared secret for this package",
	},
}

func (pmi *PackageMapperImpl) UpsertableTyped(d *schema.ResourceData) (masherytypes.Package, V3ObjectIdentifier, diag.Diagnostics) {
	rv := masherytypes.Package{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   d.Id(),
			Name: extractSetOrPrefixedString(d, MashPackName, MashPackNamePrefix),
		},

		Description:                 ExtractString(d, MashPackDescription, ""),
		NotifyDeveloperPeriod:       ExtractString(d, MashPackNotifyDeveloperPeriod, "hour"),
		NotifyDeveloperNearQuota:    ExtractBool(d, MashPackNotifyDeveloperNearQuota, true),
		NotifyDeveloperOverQuota:    ExtractBool(d, MashPackNotifyDeveloperOverQuota, true),
		NotifyDeveloperOverThrottle: ExtractBool(d, MashPackNotifyDeveloperOverThrottle, false),
		NotifyAdminPeriod:           ExtractString(d, MashPackNotifyAdminPeriod, MashDurationDay),
		NotifyAdminNearQuota:        ExtractBool(d, MashPackNotifyAdminNearQuota, false),
		NotifyAdminOverQuota:        ExtractBool(d, MashPackNotifyAdminOverQuota, false),
		NotifyAdminOverThrottle:     ExtractBool(d, MashPackNotifyDeveloperOverThrottle, false),
		NotifyAdminEmails:           strings.Join(ExtractStringArray(d, MashPackNotifyAdminEmails, &EmptyStringArray), ","),
		NearQuotaThreshold:          extractIntPointer(d, MashPackNearQuotaThreshold),
		Eav:                         ExtractStringMap(d, MashPackEAVs),
		KeyAdapter:                  ExtractString(d, MashPackKeyAdapter, ""),
		KeyLength:                   extractIntPointer(d, MashPackKeyLength),
		SharedSecretLength:          extractIntPointer(d, MashPackSharedSecretLength),
		Plans:                       nil,
	}

	return rv, nil, nil
}

func (pmi *PackageMapperImpl) splitAddressToSet(str string) []interface{} {
	if len(str) == 0 {
		return nil
	}

	split := strings.Split(str, ",")
	rv := make([]interface{}, len(split))
	for i, v := range split {
		rv[i] = v
	}

	return rv
}

func (pmi *PackageMapperImpl) PersistTyped(pack masherytypes.Package, d *schema.ResourceData) diag.Diagnostics {
	data := map[string]interface{}{
		MashPackagekId:      pack.Id,
		MashPackCreated:     pack.Created.ToString(),
		MashPackUpdated:     pack.Updated.ToString(),
		MashPackName:        pack.Name,
		MashPackDescription: pack.Description,

		MashPackNotifyDeveloperPeriod:       pack.NotifyDeveloperPeriod,
		MashPackNotifyDeveloperNearQuota:    pack.NotifyDeveloperNearQuota,
		MashPackNotifyDeveloperOverQuota:    pack.NotifyDeveloperOverQuota,
		MashPackNotifyDeveloperOverThrottle: pack.NotifyDeveloperOverThrottle,

		MashPackNotifyAdminPeriod:       pack.NotifyAdminPeriod,
		MashPackNotifyAdminNearQuota:    pack.NotifyAdminNearQuota,
		MashPackNotifyAdminOverQuota:    pack.NotifyAdminOverQuota,
		MashPackNotifyAdminOverThrottle: pack.NotifyAdminOverThrottle,

		MashPackNotifyAdminEmails:  pmi.splitAddressToSet(pack.NotifyAdminEmails),
		MashPackNearQuotaThreshold: pack.NearQuotaThreshold,
		MashPackEAVs:               pack.Eav,
		MashPackKeyAdapter:         pack.KeyAdapter,
		MashPackKeyLength:          pack.KeyLength,
		MashPackSharedSecretLength: pack.SharedSecretLength,
	}

	return pmi.persistMap(pack.Identifier(), data, d)
}

func init() {
	PackageMapper = &PackageMapperImpl{
		ResourceMapperImpl{
			schema: PackageSchema,
			v3Identity: func(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
				rv := masherytypes.PackageIdentifier{
					PackageId: d.Id(),
				}

				rvd := diag.Diagnostics{}
				if len(rv.PackageId) == 0 {
					rvd = append(rvd, PackageMapper.lackingIdentificationDiagnostic("id"))
				}
				return rv, rvd
			},
			upsertFunc: func(d *schema.ResourceData) (Upsertable, V3ObjectIdentifier, diag.Diagnostics) {
				return PackageMapper.UpsertableTyped(d)
			},
			persistFunc: func(rv interface{}, d *schema.ResourceData) diag.Diagnostics {
				ptr := rv.(*masherytypes.Package)
				return PackageMapper.PersistTyped(*ptr, d)
			},
		},
	}
}
