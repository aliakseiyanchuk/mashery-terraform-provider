package mashschema

const (
	MashPackagePlanId                         = "package_plan_id"
	MashPackagePlanRef                        = "package_plan_ref"
	MashPlanEAV                               = "extended_attribute_values"
	MashPlanSelfServiceKeyProvisioningEnabled = "self_service_provisioning"
	MashPlanAdminKeyProvisioningEnabled       = "admin_provisioning"
	MashPlanNotes                             = "notes"
	MashPlanMaxNumKeysAllowed                 = "max_keys"
	MashPlanNumKeysBeforeReview               = "keys_before_review"
	MashPlanPortalAccessRoles                 = "portal_access_roles"

	MashPlanQpsLimitCeiling            = "qps"
	MashPlanQpsLimitExempt             = "unlimited_qps"
	MashPlanQpsLimitKeyOverrideAllowed = "qps_override"

	MashPlanRateLimitCeiling            = "quota"
	MashPlanRateLimitExempt             = "unlimited_quota"
	MashPlanRateLimitKeyOverrideAllowed = "quota_override"
	MashPlanRateLimitPeriod             = "quota_period"

	MashPlanResponseFilterOverrideAllowed = "response_filter_override"
	MashPlanStatus                        = "status"
	MashPlanDeveloperEmailTemplateSetId   = "developer_email_template_set"
	MashPlanAdminEmailTemplateSetId       = "admin_email_template_set"
)

var MashPlanStatusEnum = []string{"active", "inactive"}
var RateLimitPeriodEnum = []string{MashDurationMinute, MashDurationHourly, MashDurationDay, MashDurationMonth}
