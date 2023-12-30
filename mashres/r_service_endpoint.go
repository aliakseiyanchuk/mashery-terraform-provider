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

		UpsertableFunc: func() masherytypes.Endpoint {
			return masherytypes.Endpoint{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier) (masherytypes.Endpoint, bool, error) {
			return client.GetEndpoint(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, serviceIdent masherytypes.ServiceIdentifier, m masherytypes.Endpoint) (masherytypes.Endpoint, masherytypes.ServiceEndpointIdentifier, error) {
			if createdEndpoint, err := client.CreateEndpoint(ctx, serviceIdent, m); err != nil {
				return masherytypes.Endpoint{}, masherytypes.ServiceEndpointIdentifier{}, err
			} else {
				return createdEndpoint, createdEndpoint.Identifier(), nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier, m masherytypes.Endpoint) (masherytypes.Endpoint, error) {
			m.Id = identifier.EndpointId
			m.ParentServiceId = identifier.ServiceIdentifier

			if updatedService, err := client.UpdateEndpoint(ctx, m); err != nil {
				return masherytypes.Endpoint{}, err
			} else {
				return updatedService, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointIdentifier) error {
			return client.DeleteEndpoint(ctx, identifier)
		},
	}
}
