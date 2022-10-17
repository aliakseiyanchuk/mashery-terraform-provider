package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var EndpointResource *ResourceTemplate

func init() {
	EndpointResource = &ResourceTemplate{
		Mapper: mashschema.ServiceEndpointMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetEndpoint(ctx, identifier.(masherytypes.EndpointIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreateEndpoint(ctx, identifier.(masherytypes.ServiceIdentifier), upsertable.(masherytypes.Endpoint))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdateEndpoint(ctx, upsertable.(masherytypes.Endpoint))
		},
		// The methods and method filters are deleted together with the endpoint.
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeleteEndpoint(ctx, identifier.(masherytypes.EndpointIdentifier))
		},
	}
}
