package mashery

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func mashPackageUpsertable(d *schema.ResourceData) masherytypes.MasheryPackage {
	rv := masherytypes.MasheryPackage{
		AddressableV3Object: masherytypes.AddressableV3Object{
			Id:   d.Id(),
			Name: extractSetOrPrefixedString(d, MashPackName, MashPackNamePrefix),
		},

		Description:                 extractString(d, MashPackDescription, ""),
		NotifyDeveloperPeriod:       extractString(d, MashPackNotifyDeveloperPeriod, "hour"),
		NotifyDeveloperNearQuota:    extractBool(d, MashPackNotifyDeveloperNearQuota, true),
		NotifyDeveloperOverQuota:    extractBool(d, MashPackNotifyDeveloperOverQuota, true),
		NotifyDeveloperOverThrottle: extractBool(d, MashPackNotifyDeveloperOverThrottle, false),
		NotifyAdminPeriod:           extractString(d, MashPackNotifyAdminPeriod, MashDurationDay),
		NotifyAdminNearQuota:        extractBool(d, MashPackNotifyAdminNearQuota, false),
		NotifyAdminOverQuota:        extractBool(d, MashPackNotifyAdminOverQuota, false),
		NotifyAdminOverThrottle:     extractBool(d, MashPackNotifyDeveloperOverThrottle, false),
		NotifyAdminEmails:           strings.Join(ExtractStringArray(d, MashPackNotifyAdminEmails, &EmptyStringArray), ","),
		NearQuotaThreshold:          extractIntPointer(d, MashPackNearQuotaThreshold),
		Eav:                         extractStringMap(d, MashPackEAVs),
		KeyAdapter:                  extractString(d, MashPackKeyAdapter, ""),
		KeyLength:                   extractIntPointer(d, MashPackKeyLength),
		SharedSecretLength:          extractIntPointer(d, MashPackSharedSecretLength),
		Plans:                       nil,
	}

	return rv
}

func splitAddressToSet(str string) []interface{} {
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

func v3PackageToResourceData(pack *masherytypes.MasheryPackage, d *schema.ResourceData) diag.Diagnostics {
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

		MashPackNotifyAdminEmails:  splitAddressToSet(pack.NotifyAdminEmails),
		MashPackNearQuotaThreshold: pack.NearQuotaThreshold,
		MashPackEAVs:               pack.Eav,
		MashPackKeyAdapter:         pack.KeyAdapter,
		MashPackKeyLength:          pack.KeyLength,
		MashPackSharedSecretLength: pack.SharedSecretLength,
	}

	return SetResourceFields(data, d)
}
