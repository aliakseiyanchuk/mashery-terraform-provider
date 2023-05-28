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
		UpsertableFunc: func() masherytypes.Service {
			return masherytypes.Service{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier) (*masherytypes.Service, error) {
			if service, err := client.GetService(ctx, identifier); err != nil {
				return service, err
			} else {
				if roles, rolesReadErr := client.GetServiceRoles(ctx, identifier); rolesReadErr != nil {
					return service, rolesReadErr
				} else if roles != nil {
					service.Roles = &roles
				}

				return service, nil
			}
		},

		DoCreate: func(ctx context.Context, client v3client.Client, orphan tfmapper.Orphan, m masherytypes.Service) (*masherytypes.Service, *masherytypes.ServiceIdentifier, error) {
			if cratedService, err := client.CreateService(ctx, m); err != nil {
				return nil, nil, err
			} else {
				rvIdent := cratedService.Identifier()
				var rvError error

				if m.Roles != nil {
					if readBackRoles, roleReadErr := client.GetServiceRoles(ctx, rvIdent); roleReadErr != nil {
						rvError = roleReadErr
					} else {
						cratedService.Roles = &readBackRoles
					}
				}

				return cratedService, &rvIdent, rvError
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier, m masherytypes.Service) (*masherytypes.Service, error) {
			m.Id = identifier.ServiceId

			if updatedService, err := client.UpdateService(ctx, m); err != nil {
				return nil, err
			} else {
				if m.Roles != nil {
					if readBackRoles, roleSetErr := client.GetServiceRoles(ctx, identifier); roleSetErr != nil {
						return updatedService, roleSetErr
					} else {
						updatedService.Roles = &readBackRoles
					}
				}

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
