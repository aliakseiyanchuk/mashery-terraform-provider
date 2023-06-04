package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceEndpointMethodResource *ResourceTemplate[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod]

func init() {
	ServiceEndpointMethodResource = &ResourceTemplate[masherytypes.ServiceEndpointIdentifier, masherytypes.ServiceEndpointMethodIdentifier, masherytypes.ServiceEndpointMethod]{
		Schema: mashschemag.ServiceEndpointMethodResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceEndpointMethodResourceSchemaBuilder.Mapper(),

		UpsertableFunc: func() masherytypes.ServiceEndpointMethod {
			return masherytypes.ServiceEndpointMethod{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointMethodIdentifier) (*masherytypes.ServiceEndpointMethod, error) {
			return client.GetEndpointMethod(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, serviceEndpointIdent masherytypes.ServiceEndpointIdentifier, m masherytypes.ServiceEndpointMethod) (*masherytypes.ServiceEndpointMethod, *masherytypes.ServiceEndpointMethodIdentifier, error) {
			if createdMethod, err := client.CreateEndpointMethod(ctx, serviceEndpointIdent, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := createdMethod.Identifier()
				return createdMethod, &rvIdent, nil
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointMethodIdentifier, m masherytypes.ServiceEndpointMethod) (*masherytypes.ServiceEndpointMethod, error) {
			m.Id = identifier.MethodId
			m.ParentEndpointId = identifier.ServiceEndpointIdentifier

			if updatedService, err := client.UpdateEndpointMethod(ctx, m); err != nil {
				return nil, err
			} else {
				return updatedService, err
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceEndpointMethodIdentifier) error {
			return client.DeleteEndpointMethod(ctx, identifier)
		},
	}
}
