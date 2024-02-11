package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ApplicationPackageKeyResource *ResourceTemplate[masherytypes.ApplicationIdentifier, mashschemag.ApplicationPackageKeyIdentifier, masherytypes.PackageKey]

func init() {
	ApplicationPackageKeyResource = &ResourceTemplate[masherytypes.ApplicationIdentifier, mashschemag.ApplicationPackageKeyIdentifier, masherytypes.PackageKey]{
		Schema: mashschemag.ApplicationPackageKeyResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ApplicationPackageKeyResourceSchemaBuilder.Mapper(),

		UpsertableFunc: func() masherytypes.PackageKey { return masherytypes.PackageKey{} },

		DoRead: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationPackageKeyIdentifier) (masherytypes.PackageKey, bool, error) {
			return client.GetPackageKey(ctx, identifier.PackageKeyIdentifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, parent masherytypes.ApplicationIdentifier, application masherytypes.PackageKey) (masherytypes.PackageKey, mashschemag.ApplicationPackageKeyIdentifier, error) {
			if createApp, err := client.CreatePackageKey(ctx, parent, application); err != nil {
				return masherytypes.PackageKey{}, mashschemag.ApplicationPackageKeyIdentifier{}, err
			} else {
				rvIdent := mashschemag.ApplicationPackageKeyIdentifier{
					ApplicationIdentifier: parent,
					PackageKeyIdentifier:  createApp.Identifier(),
				}
				return createApp, rvIdent, err
			}
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationPackageKeyIdentifier, pacakgeKey masherytypes.PackageKey) (masherytypes.PackageKey, error) {
			pacakgeKey.Id = identifier.PackageKeyId
			return client.UpdatePackageKey(ctx, pacakgeKey)
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier mashschemag.ApplicationPackageKeyIdentifier) error {
			return client.DeletePackageKey(ctx, identifier.PackageKeyIdentifier)
		},

		// There are no offending objects
	}
}
