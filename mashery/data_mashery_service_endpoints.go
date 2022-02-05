package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Query for Mashery service endpoints that help searching for the endpoints.
func dataSourceMasheryServiceEndpoints() *schema.Resource {
	return &schema.Resource{
		Schema:      DataSourceSvcEndpointsSchema,
		ReadContext: readDataSourceMasheryServiceEndpoints,
	}
}

func readDataSourceMasheryServiceEndpoints(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	check := validateRegularExpressionSet(d.Get(MashDataSourceNameFilterRegexp), cty.GetAttrPath(MashDataSourceNameFilterRegexp))
	if len(check) > 0 {
		return check
	}

	check = validateRegularExpressionSet(d.Get(DataSourceServiceEndpointPathRegexp), cty.GetAttrPath(DataSourceServiceEndpointPathRegexp))
	if len(check) > 0 {
		return check
	}

	cl := m.(v3client.Client)

	// Input may contain a comment
	srvId := IdWithoutComment(d.Get(MashSvcId).(string))

	if endpoints, err := cl.ListEndpointsWithFullInfo(ctx, srvId); err != nil {
		return diag.FromErr(err)
	} else {

		nameFilter := ExtractStringArray(d, MashDataSourceNameFilterRegexp, &EmptyStringArray)
		pathFilter := ExtractStringArray(d, MashDataSourceNameFilterRegexp, &EmptyStringArray)

		if len(nameFilter) > 0 || len(pathFilter) > 0 {
			endpoints = filterEndpoints(endpoints, nameFilter, pathFilter)
		}

		if len(endpoints) > 0 {
			rv := make([]interface{}, len(endpoints))
			expl := make(map[string]interface{}, len(endpoints))

			for idx, v := range endpoints {
				eid := ServiceEndpointIdentifier{
					ServiceId:  srvId,
					EndpointId: v.Id,
				}
				rv[idx] = eid.Id()
				expl[eid.Id()] = v.Name
			}

			d.SetId("multiple")
			setData := map[string]interface{}{
				MashEndpointMultiRef:   rv,
				MashEndpointsExplained: expl,
			}

			return SetResourceFields(setData, d)
		} else {
			_ = d.Set(MashEndpointMultiRef, []interface{}{})
		}

		return diag.Diagnostics{}
	}
}

func filterEndpoints(inp []masherytypes.MasheryEndpoint, name, path []string) []masherytypes.MasheryEndpoint {
	nameRx := toRegexpArray(name)
	pathRx := toRegexpArray(path)

	rv := []masherytypes.MasheryEndpoint{}

	for _, v := range inp {
		nameMatched := false
		pathMatched := false

		if len(nameRx) > 0 {
			for _, rgx := range nameRx {
				if rgx.MatchString(v.Name) {
					nameMatched = true
					break
				}
			}
		} else {
			nameMatched = true
		}

		if len(pathRx) > 0 {
			for _, rgx := range pathRx {
				if rgx.MatchString(v.RequestPathAlias) {
					pathMatched = true
					break
				}
			}
		} else {
			pathMatched = true
		}

		if nameMatched && pathMatched {
			rv = append(rv, v)
		}
	}

	return rv
}
