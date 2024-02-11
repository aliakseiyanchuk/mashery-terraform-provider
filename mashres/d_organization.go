package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var OrganizationDataSource = CreateSingularDataSource(mashschemag.OrganizationResourceSchemaBuilder, queryOrganization)

func queryOrganization(ctx context.Context, client v3client.Client, m map[string]string) (string, *masherytypes.Organization, error) {
	if sets, err := client.ListOrganizationsFiltered(ctx, m); err != nil {
		return "", nil, err
	} else if len(sets) == 1 {
		return sets[0].Id, &sets[0], nil
	} else {
		return "", nil, nil
	}
}
