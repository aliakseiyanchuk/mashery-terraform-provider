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

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) (*masherytypes.ServiceCache, error) {
			return client.GetServiceCache(ctx, identifier.ServiceId)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, ident masherytypes.ServiceIdentifier, m masherytypes.ServiceCache) (*masherytypes.ServiceCache, *masherytypes.ServiceIdentifier, error) {
			// Align the API of the Mashery V3 client
			if createdCache, err := client.CreateServiceCache(ctx, ident.ServiceId, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := masherytypes.ServiceIdentifier{
					ServiceId: ident.ServiceId,
				}
				return createdCache, &rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier, m masherytypes.ServiceCache) (*masherytypes.ServiceCache, error) {
			if updatedCache, err := client.UpdateServiceCache(ctx, identifier.ServiceId, m); err != nil {
				return nil, err
			} else {
				return updatedCache, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) error {
			return client.DeleteServiceCache(ctx, identifier.ServiceId)
		},

		// Offending count is not required for Service cache, as it can be deleted at any moment.
	}
}
