package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanServiceEndpointResource *ResourceTemplate[masherytypes.PackagePlanServiceIdentifier, masherytypes.PackagePlanServiceEndpointIdentifier, mashschemag.PackagePlanServiceEndpointParam]

func init() {
	PackagePlanServiceEndpointResource = &ResourceTemplate[masherytypes.PackagePlanServiceIdentifier, masherytypes.PackagePlanServiceEndpointIdentifier, mashschemag.PackagePlanServiceEndpointParam]{
		Schema: mashschemag.PackagePlanServiceEndpointResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackagePlanServiceEndpointResourceSchemaBuilder.Mapper(),

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointIdentifier) (*mashschemag.PackagePlanServiceEndpointParam, error) {
			serviceExists, err := client.CheckPlanEndpointExists(ctx, identifier)
			if serviceExists {
				return &mashschemag.PackagePlanServiceEndpointParam{}, err
			} else {
				return nil, err
			}
		},

		DoCreate: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceIdentifier, m mashschemag.PackagePlanServiceEndpointParam) (*mashschemag.PackagePlanServiceEndpointParam, *masherytypes.PackagePlanServiceEndpointIdentifier, error) {

			ident := masherytypes.PackagePlanServiceEndpointIdentifier{
				ServiceEndpointIdentifier: m.ServiceEndpointIdentifier,
				PackagePlanIdentifier:     identifier.PackagePlanIdentifier,
			}

			if _, err := client.CreatePlanEndpoint(ctx, ident); err != nil {
				return nil, nil, err
			} else {
				return &m, &ident, nil
			}
		},

		// Update is not required: it will be delete-only

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointIdentifier) error {
			return client.DeletePlanEndpoint(ctx, identifier)
		},

		// Offending methods is not required: it's a sub-part of the endpoint, and it doesn't need to be specifically
		// tracked.
	}
}
