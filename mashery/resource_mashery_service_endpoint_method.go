package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschema"
)

var EndpointMethodResource *ResourceTemplate

func init() {
	EndpointMethodResource = &ResourceTemplate{
		Mapper: mashschema.ServiceEndpointMethodMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetEndpointMethod(ctx, identifier.(masherytypes.EndpointMethodIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreateEndpointMethod(ctx, identifier.(masherytypes.EndpointIdentifier), upsertable.(masherytypes.Method))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdateEndpointMethod(ctx, upsertable.(masherytypes.Method))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeleteEndpointMethod(ctx, identifier.(masherytypes.EndpointMethodIdentifier))
		},
	}
}
