package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var systemDomainsDataSource *DatasourceTemplate
var publicDomainsDataSource *DatasourceTemplate

func init() {
	systemDomainsDataSource = &DatasourceTemplate{
		V3ObjectTypeName: "system domain",
		QueryExtractor:   DataSourceDefaultQueryExtractor,
		Query: func(ctx context.Context, cl v3client.Client, query map[string]string) ([]interface{}, error) {
			if v3Objects, err := cl.GetSystemDomains(ctx); err != nil {
				return nil, err
			} else {
				return mashschema.CoerceStringArrayToInterfaceType(v3Objects), nil
			}
		},
		Mapper: mashschema.DomainsMapper,
	}

	publicDomainsDataSource = &DatasourceTemplate{
		V3ObjectTypeName: "public domain",
		QueryExtractor:   DataSourceDefaultQueryExtractor,
		Query: func(ctx context.Context, cl v3client.Client, query map[string]string) ([]interface{}, error) {
			if v3Objects, err := cl.GetPublicDomains(ctx); err != nil {
				return nil, err
			} else {
				return mashschema.CoerceStringArrayToInterfaceType(v3Objects), nil
			}
		},
		Mapper: mashschema.DomainsMapper,
	}
}
