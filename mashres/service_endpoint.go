package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceEndpointResource *ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]

func init() {
	ServiceEndpointResource = &ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ServiceEndpointIdentifier, masherytypes.Endpoint]{
		Schema: mashschemag.ServiceEndpointResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceEndpointResourceSchemaBuilder.Mapper(),

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier) (*masherytypes.Endpoint, error) {
			return client.GetEndpoint(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, serviceIdent masherytypes.ServiceIdentifier, m masherytypes.Endpoint) (*masherytypes.Endpoint, *masherytypes.ServiceEndpointIdentifier, error) {
			if createdEndpoint, err := client.CreateEndpoint(ctx, serviceIdent, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := createdEndpoint.Identifier()
				return createdEndpoint, &rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier, m masherytypes.Endpoint) (*masherytypes.Endpoint, error) {
			m.Id = identifier.ServiceId

			if updatedService, err := client.UpdateEndpoint(ctx, m); err != nil {
				return nil, err
			} else {
				return updatedService, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier) error {
			return client.DeleteEndpoint(ctx, identifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier) (int64, error) {
			return client.CountEndpointsMethodsOf(ctx, identifier)
		},
	}
}
