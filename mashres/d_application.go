package mashres

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashschemag"
)

var ApplicationDataSource = CreateSingularDataSource(
	"application",
	mashschemag.ApplicationResourceSchemaBuilder,
	queryApplication,
)

func queryApplication(ctx context.Context, client v3client.Client, m map[string]string) (mashschemag.ApplicationOfMemberIdentifier, *masherytypes.Application, error) {
	if sets, err := client.ListApplicationsFiltered(ctx, m); err != nil {
		return mashschemag.ApplicationOfMemberIdentifier{}, nil, err
	} else if len(sets) == 1 {
		rv := mashschemag.ApplicationOfMemberIdentifier{
			MemberIdentifier: masherytypes.MemberIdentifier{
				MemberId: "",
				Username: sets[0].Username,
			},
			ApplicationIdentifier: masherytypes.ApplicationIdentifier{
				ApplicationId: sets[0].Id,
			},
		}
		return rv, &sets[0], nil
	} else {
		return mashschemag.ApplicationOfMemberIdentifier{}, nil, nil
	}
}
