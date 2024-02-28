package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ApplicationResource *ResourceTemplate[masherytypes.MemberIdentifier, mashschemag.ApplicationOfMemberIdentifier, masherytypes.Application]

func init() {
	ApplicationResource = &ResourceTemplate[masherytypes.MemberIdentifier, mashschemag.ApplicationOfMemberIdentifier, masherytypes.Application]{
		Schema: mashschemag.ApplicationResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ApplicationResourceSchemaBuilder.Mapper(),

		UpsertableFunc: func() masherytypes.Application { return masherytypes.Application{} },

		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier) (masherytypes.Application, bool, error) {
			if application, appExists, err := client.GetApplication(ctx, identifier.ApplicationIdentifier); err != nil {
				return application, appExists, err
			} else if appExists {
				// Read the EAV attributes of the application exists
				if eavs, eavReadErr := client.GetApplicationExtendedAttributes(ctx, identifier.ApplicationIdentifier); eavReadErr != nil {
					return application, appExists, eavReadErr
				} else {
					application.Eav = eavs
				}

				return application, appExists, err
			} else {
				// Default: return what the Mashery client has provided
				return application, appExists, err
			}
		},

		DoCreate: func(ctx context.Context, client v3client.Client, parent masherytypes.MemberIdentifier, application masherytypes.Application) (masherytypes.Application, mashschemag.ApplicationOfMemberIdentifier, error) {
			if createApp, err := client.CreateApplication(ctx, parent, application); err != nil {
				return masherytypes.Application{}, mashschemag.ApplicationOfMemberIdentifier{}, err
			} else {
				rvIdent := mashschemag.ApplicationOfMemberIdentifier{
					MemberIdentifier:      parent,
					ApplicationIdentifier: createApp.Identifier(),
				}

				// If application contains EAVs, then these will be set on application creation
				if len(application.Eav) > 0 {
					if updatedEav, eavErr := client.UpdateApplicationExtendedAttributes(ctx, rvIdent.ApplicationIdentifier, application.Eav); eavErr != nil {
						return createApp, rvIdent, eavErr
					} else {
						createApp.Eav = updatedEav
					}
				}

				return createApp, rvIdent, err
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier, application masherytypes.Application) (masherytypes.Application, error) {
			application.Id = identifier.ApplicationId

			if updatedApp, updateErr := client.UpdateApplication(ctx, application); updateErr != nil {
				return updatedApp, updateErr
			} else {
				if len(application.Eav) > 0 {
					if updatedEavs, eavUpdateErr := client.UpdateApplicationExtendedAttributes(ctx, identifier.ApplicationIdentifier, application.Eav); eavUpdateErr != nil {
						return updatedApp, eavUpdateErr
					} else {
						updatedApp.Eav = updatedEavs
					}
				}

				return updatedApp, updateErr
			}
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier) error {
			return client.DeleteApplication(ctx, identifier.ApplicationIdentifier)
		},

		DoCountOffending: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationOfMemberIdentifier) (int64, error) {
			return client.CountApplicationPackageKeys(ctx, identifier.ApplicationIdentifier)
		},
	}
}
