package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceCacheResource *ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache]

func init() {
	ServiceCacheResource = &ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ServiceIdentifier, masherytypes.ServiceCache]{
		Schema: mashschemag.ServiceCacheResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceCacheResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/services/(.+)/cache", masherytypes.ServiceIdentifier{}, func(items []string) masherytypes.ServiceIdentifier {
			return masherytypes.ServiceIdentifier{ServiceId: items[1]}
		}),

		UpsertableFunc: func() masherytypes.ServiceCache {
			return masherytypes.ServiceCache{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) (masherytypes.ServiceCache, bool, error) {
			return client.GetServiceCache(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, ident masherytypes.ServiceIdentifier, m masherytypes.ServiceCache) (masherytypes.ServiceCache, masherytypes.ServiceIdentifier, error) {
			// Align the API of the Mashery V3 client
			if createdCache, err := client.CreateServiceCache(ctx, ident, m); err != nil {
				return masherytypes.ServiceCache{}, masherytypes.ServiceIdentifier{}, err
			} else {
				rvIdent := masherytypes.ServiceIdentifier{
					ServiceId: ident.ServiceId,
				}
				return createdCache, rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier, m masherytypes.ServiceCache) (masherytypes.ServiceCache, error) {
			m.ParentServiceId = identifier

			if updatedCache, err := client.UpdateServiceCache(ctx, m); err != nil {
				return masherytypes.ServiceCache{}, err
			} else {
				return updatedCache, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) error {
			return client.DeleteServiceCache(ctx, identifier)
		},

		// Offending count is not required for Service cache, as it can be deleted at any moment.
	}
}
