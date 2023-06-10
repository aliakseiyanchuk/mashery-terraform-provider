package mashschema

const (
	MashPlanId                                = "plan_id"
	MashPlanRef                               = "plan_ref"
	MashPlanCreated                           = "created"
	MashPlanUpdated                           = "updated"
	MashPlanName                              = "name"
	MashPlanDescription                       = "description"
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
	MashPlanEmailTemplateSetId            = "email_template_set"
	MashPlanAdminEmailTemplateSetId       = "admin_email_template_set"
)

var RateLimitPeriodEnum = []string{MashDurationMinute, MashDurationHourly, MashDurationDay, MashDurationMonth}
