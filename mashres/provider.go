package mashres

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"time"
)

func ConfigureCaching(ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
	redisServer := d.Get(mashery.ProviderRedisCacheField).(string)
	cacheDurationStr := d.Get(mashery.ProviderCacheDuration).(string)

	var rv diag.Diagnostics

	if len(redisServer) > 0 {

		if initErr := InitRedisClient(redisServer); initErr != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "cannot initialize Redis cache server",
				Detail:   fmt.Sprintf("Redis cache initialisation returned the following error: %s", initErr.Error()),
			})
		}

		var parseErr error
		if DefaultCacheDuration, parseErr = time.ParseDuration(cacheDurationStr); parseErr != nil {
			rv = append(rv, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "incorrect cache duration",
				Detail:   fmt.Sprintf("%s is an incorrect time duration format: %s", cacheDurationStr, parseErr.Error()),
			})
		}

		tflog.Info(ctx, fmt.Sprintf("Mashery provider will use Redis cache at %s with cache duration of %s", redisServer, cacheDurationStr))
	} else {
		tflog.Info(ctx, "Mashery provider doesn't have sufficient information to use Redis caches in data sources")
	}

	return rv
}

// Provider Mashery Terraform Provider mashschema definition
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: mashery.ProviderConfigSchema,
		ResourcesMap: map[string]*schema.Resource{
			"mashery_member":                               MemberResource.ResourceSchema(),
			"mashery_application":                          ApplicationResource.ResourceSchema(),
			"mashery_application_package_key":              ApplicationPackageKeyResource.ResourceSchema(),
			"mashery_service":                              ServiceResource.ResourceSchema(),
			"mashery_service_cache":                        ServiceCacheResource.ResourceSchema(),
			"mashery_service_oauth":                        ServiceOAuthResource.ResourceSchema(),
			"mashery_service_error_set":                    ServiceErrorSetResource.ResourceSchema(),
			"mashery_service_endpoint":                     ServiceEndpointResource.ResourceSchema(),
			"mashery_service_endpoint_method":              ServiceEndpointMethodResource.ResourceSchema(),
			"mashery_service_endpoint_method_filter":       ServiceEndpointMethodFilterResource.ResourceSchema(),
			"mashery_package":                              PackageResource.ResourceSchema(),
			"mashery_package_plan":                         PackagePlanResource.ResourceSchema(),
			"mashery_package_plan_service":                 PackagePlanServiceResource.ResourceSchema(),
			"mashery_package_plan_service_endpoint":        PackagePlanServiceEndpointResource.ResourceSchema(),
			"mashery_package_plan_service_endpoint_method": PackagePlanServiceEndpointMethodResource.ResourceSchema(),
			"mashery_unique_path":                          mashery.ResourceMasheryUniquePath(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mashery_member":             MemberDataSource.DataSourceSchema(),
			"mashery_package":            PackageDataSource.DataSourceSchema(),
			"mashery_package_plan":       PackagePlanDataSource.DataSourceSchema(),
			"mashery_application":        ApplicationDataSource.DataSourceSchema(),
			"mashery_organization":       OrganizationDataSource.DataSourceSchema(),
			"mashery_email_template_set": EmailTemplateSetDataSource.DataSourceSchema(),
			"mashery_role":               RoleDataSource.DataSourceSchema(),
		},
		ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			cacheDiags := ConfigureCaching(ctx, data)
			v3Cl, providerDiags := mashery.ProviderConfigure(ctx, data)

			return v3Cl, append(cacheDiags, providerDiags...)
		},
	}
}
