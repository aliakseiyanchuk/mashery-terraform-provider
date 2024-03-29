package mashschema

const (
	MashPackageId                       = "package_id"
	MashPackageRef                      = "package_ref"
	MashPackCreated                     = "created"
	MashPackUpdated                     = "updated"
	MashPackDescription                 = "description"
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

	MashDurationMinute = "minute"
	MashDurationHourly = "hour"
	MashDurationDay    = "day"
	MashDurationWeek   = "week"
	MashDurationMonth  = "month"
)

var NotifyDeveloperPeriodEnum = []string{MashDurationMinute, MashDurationHourly,
	MashDurationDay, MashDurationWeek, MashDurationMonth}
var NotifyAdminPeriodEnum = []string{MashDurationMinute, MashDurationHourly,
	MashDurationDay, MashDurationWeek, MashDurationMonth}
