package mashres

import (
	"context"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var PackagePlanServiceEndpointResource *ResourceTemplate[masherytypes.PackagePlanServiceIdentifier, masherytypes.PackagePlanServiceEndpointIdentifier, mashschemag.PackagePlanServiceEndpointParam]

func init() {
	PackagePlanServiceEndpointResource = &ResourceTemplate[masherytypes.PackagePlanServiceIdentifier, masherytypes.PackagePlanServiceEndpointIdentifier, mashschemag.PackagePlanServiceEndpointParam]{
		Schema: mashschemag.PackagePlanServiceEndpointResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.PackagePlanServiceEndpointResourceSchemaBuilder.Mapper(),

		UpsertableFunc: func() mashschemag.PackagePlanServiceEndpointParam {
			return mashschemag.PackagePlanServiceEndpointParam{}
		},

		ValidateFunc: func(parent masherytypes.PackagePlanServiceIdentifier, upsertable mashschemag.PackagePlanServiceEndpointParam) string {
			// The endpoint should be included in the same service. Mixing parameters is not allowed
			if parent.ServiceId != upsertable.ServiceEndpointIdentifier.ServiceId {
				return fmt.Sprintf("parameter conflict: endpoint %s belongs to service %s while only endpoints of service %s are expected",
					upsertable.ServiceEndpointIdentifier.EndpointId,
					upsertable.ServiceEndpointIdentifier.ServiceId,
					parent.ServiceId,
				)
			}

			return ""
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.PackagePlanServiceEndpointIdentifier) (*mashschemag.PackagePlanServiceEndpointParam, error) {
			ppsEndpointExists, err := client.CheckPlanEndpointExists(ctx, identifier)
			if ppsEndpointExists {
				// If endpoint exists; then the reference to the same identifier will be returned. The value
				// in the state file will not change as this is an identity operation
				return &mashschemag.PackagePlanServiceEndpointParam{
					ServiceEndpointIdentifier: identifier.ServiceEndpointIdentifier,
				}, err
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
