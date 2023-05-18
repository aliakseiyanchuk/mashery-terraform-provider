package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
	"terraform-provider-mashery/tfmapper"
)

var ServiceResource *ResourceTemplate[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]

func init() {
	ServiceResource = &ResourceTemplate[tfmapper.Orphan, masherytypes.ServiceIdentifier, masherytypes.Service]{
		Schema: mashschemag.ServiceResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceResourceSchemaBuilder.Mapper(),

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) (*masherytypes.Service, error) {
			return client.GetService(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, orphan tfmapper.Orphan, m masherytypes.Service) (*masherytypes.Service, *masherytypes.ServiceIdentifier, error) {
			if cratedService, err := client.CreateService(ctx, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := cratedService.Identifier()
				return cratedService, &rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier, m masherytypes.Service) (*masherytypes.Service, error) {
			m.Id = identifier.ServiceId

			if updatedService, err := client.UpdateService(ctx, m); err != nil {
				return nil, err
			} else {
				return updatedService, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) error {
			return client.DeleteService(ctx, identifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) (int64, error) {
			return client.CountEndpointsOf(ctx, identifier)
		},
	}
}
