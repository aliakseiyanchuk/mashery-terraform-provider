package mashery

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const MashSvcCacheTtl = "cache_ttl"

var ServiceCacheSchema = map[string]*schema.Schema{
	MashSvcCacheTtl: {
		Type:        schema.TypeInt,
		Description: "Time till which the data is stored in cache",
		Required:    true,
	},
}

var DataSourceServiceCacheScheme = map[string]*schema.Schema{
	MashSvcId: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Service Id to look up",
	},
}

// The data source would inherit cache ttl settings in a read-only mode.
func init() {
	computedTTL := cloneAsComputed(ServiceCacheSchema)
	inheritAll(&DataSourceServiceCacheScheme, &computedTTL)
}
