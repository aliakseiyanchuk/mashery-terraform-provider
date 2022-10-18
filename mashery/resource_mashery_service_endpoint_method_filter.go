package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	mashschema "terraform-provider-mashery/mashschema"
)

var EndpointMethodFilterResponse *ResourceTemplate

func init() {
	EndpointMethodFilterResponse = &ResourceTemplate{
		Mapper: mashschema.ServiceEndpointMethodFilterMapper,
		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.GetEndpointMethodFilter(ctx, identifier.(masherytypes.ServiceEndpointMethodFilterIdentifier))
		},
		DoCreate: func(ctx context.Context, client v3client.Client, obj mashschema.Upsertable, ident mashschema.V3ObjectIdentifier) (mashschema.Upsertable, error) {
			return client.CreateEndpointMethodFilter(ctx,
				ident.(masherytypes.ServiceEndpointMethodIdentifier),
				obj.(masherytypes.ServiceEndpointMethodFilter))
		},
		DoUpdate: func(ctx context.Context, client v3client.Client, upsertable mashschema.Upsertable) (mashschema.Upsertable, error) {
			return client.UpdateEndpointMethodFilter(ctx, upsertable.(masherytypes.ServiceEndpointMethodFilter))
		},
		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschema.V3ObjectIdentifier) error {
			return client.DeleteEndpointMethodFilter(ctx, identifier.(masherytypes.ServiceEndpointMethodFilterIdentifier))
		},
	}
}
