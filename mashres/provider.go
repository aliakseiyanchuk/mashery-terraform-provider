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
			"mashery_service_endpoint":                     ServiceEndpointResource.ResourceSchema(),
			"mashery_endpoint_method":                      ServiceEndpointMethodResource.ResourceSchema(),
			"mashery_endpoint_method_filter":               ServiceEndpointMethodFilterResource.ResourceSchema(),
			"mashery_package":                              PackageResource.ResourceSchema(),
			"mashery_package_plan":                         PackagePlanResource.ResourceSchema(),
			"mashery_package_plan_service":                 PackagePlanServiceResource.ResourceSchema(),
			"mashery_package_plan_service_endpoint":        PackagePlanServiceEndpointResource.ResourceSchema(),
			"mashery_package_plan_service_endpoint_method": PackagePlanServiceEndpointMethodResource.ResourceSchema(),
			//"mashery_member":                       MemberResource.TFDataSourceSchema(),
			//"mashery_application":                  ApplicationResource.TFDataSourceSchema(),
			//"mashery_package_key":                  PackageKeyResource.TFDataSourceSchema(),
			"mashery_unique_path": mashery.ResourceMasheryUniquePath(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"mashery_system_domains":     systemDomainsDataSource.TFDataSourceSchema(),
			//"mashery_public_domains":     publicDomainsDataSource.TFDataSourceSchema(),
			"mashery_organization":       OrganizationDataSource.DataSourceSchema(),
			"mashery_email_template_set": EmailTemplateSetDataSource.DataSourceSchema(),
			"mashery_role":               RoleDataSource.DataSourceSchema(),
		},
		ConfigureContextFunc: mashery.ProviderConfigure,
	}
}
