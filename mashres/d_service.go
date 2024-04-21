package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ServiceDataSource = CreateSingularDataSource(
	"service",
	mashschemag.ServiceResourceSchemaBuilder,
	queryService,
)

func queryService(ctx context.Context, client v3client.Client, m map[string]string) (masherytypes.ServiceIdentifier, *masherytypes.Service, error) {
	if sets, err := client.ListServicesFiltered(ctx, m); err != nil {
		return masherytypes.ServiceIdentifier{}, nil, err
	} else if len(sets) == 1 {
		return sets[0].Identifier(), &sets[0], nil
	} else {
		return masherytypes.ServiceIdentifier{}, nil, nil
	}
}
