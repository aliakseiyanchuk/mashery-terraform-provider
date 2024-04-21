package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ApplicationPackageKeyResource *ResourceTemplate[masherytypes.ApplicationIdentifier, masherytypes.ApplicationPackageKeyIdentifier, masherytypes.ApplicationPackageKey]

func init() {
	ApplicationPackageKeyResource = &ResourceTemplate[masherytypes.ApplicationIdentifier, masherytypes.ApplicationPackageKeyIdentifier, masherytypes.ApplicationPackageKey]{
		Schema: mashschemag.ApplicationPackageKeyResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ApplicationPackageKeyResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/applications/(.+)/packageKeys/(.+)",
			masherytypes.ApplicationPackageKeyIdentifier{},
			func(items []string) masherytypes.ApplicationPackageKeyIdentifier {
				rv := masherytypes.ApplicationPackageKeyIdentifier{}
				rv.ApplicationId = items[1]
				rv.PackageKeyId = items[2]

				return rv
			}),

		UpsertableFunc: func() masherytypes.ApplicationPackageKey { return masherytypes.ApplicationPackageKey{} },

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ApplicationPackageKeyIdentifier) (masherytypes.ApplicationPackageKey, bool, error) {
			return client.GetApplicationPackageKey(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, parent masherytypes.ApplicationIdentifier, appKey masherytypes.ApplicationPackageKey) (masherytypes.ApplicationPackageKey, masherytypes.ApplicationPackageKeyIdentifier, error) {
			if createApp, err := client.CreateApplicationPackageKey(ctx, parent, appKey); err != nil {
				return masherytypes.ApplicationPackageKey{}, masherytypes.ApplicationPackageKeyIdentifier{}, err
			} else {
				return createApp, createApp.Identifier(), err
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ApplicationPackageKeyIdentifier, pacakgeKey masherytypes.ApplicationPackageKey) (masherytypes.ApplicationPackageKey, error) {
			pacakgeKey.Id = identifier.PackageKeyId
			pacakgeKey.ParentApplicationId = identifier.ApplicationIdentifier

			// Delete the reference to the package during update. The mapper is responsible for checking
			// that the package identifier has not changed.
			pacakgeKey.Package = nil
			return client.UpdateApplicationPackageKey(ctx, pacakgeKey)
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ApplicationPackageKeyIdentifier) error {
			return client.DeletePackageKey(ctx, identifier.PackageKeyIdentifier)
		},

		// There are no offending objects
	}
}
