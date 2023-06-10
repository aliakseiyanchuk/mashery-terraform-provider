package mashschema

const (
	MashPackageId                       = "package_id"
	MashPackageRef                      = "package_ref"
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

var NotifyDeveloperPeriodEnum = []string{MashDurationMinute, MashDurationHourly,
	MashDurationDay, MashDurationWeek, MashDurationMonth}
var notifyAdminPeriodEnum = []string{MashDurationMinute, MashDurationHourly,
	MashDurationDay, MashDurationWeek, MashDurationMonth}
