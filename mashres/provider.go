package mashres

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
)

// Provider Mashery Terraform Provider mashschema definition
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: mashery.ProviderConfigSchema,
		ResourcesMap: map[string]*schema.Resource{
			"mashery_service":       ServiceResource.ResourceSchema(),
			"mashery_service_cache": ServiceCacheResource.ResourceSchema(),
			"mashery_service_oauth": ServiceOAuthResource.ResourceSchema(),
			//"mashery_service_error_set":            resourceMasheryErrorSet(),
			//"mashery_processor_chain":              resourceMasheryProcessorChain(),
			"mashery_service_endpoint": ServiceEndpointResource.ResourceSchema(),
			//"mashery_endpoint_method":              EndpointMethodResource.TFDataSourceSchema(),
			//"mashery_endpoint_method_filter":       EndpointMethodFilterResponse.TFDataSourceSchema(),
			//"mashery_package":                      PackageResource.TFDataSourceSchema(),
			//"mashery_package_plan":                 PackagePlanResource.TFDataSourceSchema(),
			//"mashery_package_plan_service":         PackagePlanServiceResource.TFDataSourceSchema(),
			//"mashery_package_plan_endpoint":        PackagePlanServiceEndpointResource.TFDataSourceSchema(),
			//"mashery_package_plan_endpoint_method": resourceMasheryPlanMethod(),
			//"mashery_member":                       MemberResource.TFDataSourceSchema(),
			//"mashery_application":                  ApplicationResource.TFDataSourceSchema(),
			//"mashery_package_key":                  PackageKeyResource.TFDataSourceSchema(),
			"mashery_unique_path": mashery.ResourceMasheryUniquePath(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"mashery_system_domains":     systemDomainsDataSource.TFDataSourceSchema(),
			//"mashery_public_domains":     publicDomainsDataSource.TFDataSourceSchema(),
			"mashery_email_template_set": EmailTemplateSetDataSource.DataSourceSchema(),
			"mashery_role":               RoleDataSource.DataSourceSchema(),
		},
		ConfigureContextFunc: mashery.ProviderConfigure,
	}
}
