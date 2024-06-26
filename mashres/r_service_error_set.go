package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceErrorSetResource *ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]

func init() {
	ServiceErrorSetResource = &ResourceTemplate[masherytypes.ServiceIdentifier, masherytypes.ErrorSetIdentifier, masherytypes.ErrorSet]{
		Schema: mashschemag.ServiceErrorSetResourceSchemaBuilder.ResourceSchema(),
		Mapper: mashschemag.ServiceErrorSetResourceSchemaBuilder.Mapper(),

		ImportIdentityParser: regexpImportIdentityParser("/services/(.+)/errorSets/(.+)", masherytypes.ErrorSetIdentifier{}, func(items []string) masherytypes.ErrorSetIdentifier {
			rv := masherytypes.ErrorSetIdentifier{}
			rv.ServiceId = items[1]
			rv.ErrorSetId = items[2]

			return rv
		}),

		UpsertableFunc: func() masherytypes.ErrorSet {
			return masherytypes.ErrorSet{}
		},

		DoRead: func(ctx context.Context, client v3client.Client, identifier masherytypes.ErrorSetIdentifier) (masherytypes.ErrorSet, bool, error) {
			return client.GetErrorSet(ctx, identifier)
		},

		DoCreate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ServiceIdentifier, set masherytypes.ErrorSet) (masherytypes.ErrorSet, masherytypes.ErrorSetIdentifier, error) {
			readBack, err := client.CreateErrorSet(ctx, identifier, set)

			var rvIdent masherytypes.ErrorSetIdentifier
			if err == nil {
				rvIdent = readBack.Identifier()
			}

			return readBack, rvIdent, err
		},

		DoUpdate: func(ctx context.Context, client v3client.Client, identifier masherytypes.ErrorSetIdentifier, set masherytypes.ErrorSet) (masherytypes.ErrorSet, error) {
			set.ParentServiceId = identifier.ServiceIdentifier

			return client.UpdateErrorSet(ctx, set)
		},

		DoDelete: func(ctx context.Context, client v3client.Client, identifier masherytypes.ErrorSetIdentifier) error {
			return client.DeleteErrorSet(ctx, identifier)
		},
	}
}
