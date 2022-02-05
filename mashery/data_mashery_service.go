package mashery

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/masherytypes"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Mashery data implementation

func dataSourceMasheryServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	check := validateRegularExpressionSet(d.Get(MashDataSourceNameFilterRegexp), cty.GetAttrPath(MashDataSourceNameFilterRegexp))
	if len(check) > 0 {
		return check
	}

	doLogf("== Pre-conditions to read mashery services has been met ==")
	cl := m.(v3client.Client)

	query := extractStringMap(d, MashDataSourceSearch)
	doLogJson("Query parameters", query)

	doLogf("Submitting query:->")
	if srv, err := cl.ListServicesFiltered(ctx, query, v3client.MasheryServiceFieldsWithEndpoinds); err != nil {
		doLogf("<- Found service id=%s", srv[0].Id)
		return diag.FromErr(err)
	} else {
		doLogf("<- Found %d service matching the query, before filtering is applied", len(srv))
		return DataSourceSvcReadReturned(srv, d)
	}
}

func DataSourceSvcReadReturned(srv []masherytypes.MasheryService, d *schema.ResourceData) diag.Diagnostics {
	if len(srv) > 0 {
		nameExpr := ExtractStringArray(d, MashDataSourceNameFilterRegexp, &EmptyStringArray)
		if len(nameExpr) > 0 {
			srv = FilterMasherySvcName(&srv, nameExpr)
			doLogf("~ After name regexp filtering, there are %d services left", len(srv))
		}
	}

	rv := diag.Diagnostics{}

	// Save received service Ids
	if len(srv) > 0 {
		// The Ids will be set to the output field.
		ids := make([]string, len(srv))
		expl := map[string]interface{}{}

		for idx, v := range srv {
			ids[idx] = v.Id
			expl[v.Id] = v.Name
		}

		dataSet := map[string]interface{}{
			MashSvcMultiRef:  ids,
			MashSvcExplained: expl,
		}
		rv = append(rv, SetResourceFields(dataSet, d)...)
	}

	// Match is not unique
	if len(srv) > 1 {
		// If  the configuration expects a unique match, then an error will be returned.
		if extractBool(d, MashDataSourceUnique, false) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Multiple matches",
				Detail:        "Multiple services (%d) found matching search and filter criteria, where the configuration requires unique",
				AttributePath: nil,
			}}
		} else {
			d.SetId("multiple")
			return rv
		}
	} else if len(srv) == 1 {
		doLogf("<= Selected service id=%s to match the search criteria", srv[0].Id)
		d.SetId(srv[0].Id)
		rv = append(rv, V3ServiceToTerraform(&srv[0], d)...)

		return rv
	} else {
		doLogf("No service found for the specified filter on resource %s", d.Id())
		d.SetId("")

		if extractBool(d, MashDataSourceRequired, true) {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Required V3 service not found",
				Detail:   "No Mashery service matching this query exists in this area",
			}}
		}

		// Returning any conflicts that were found.
		return rv
	}
}

func FilterMasherySvcName(rv *[]masherytypes.MasheryService, exp []string) []masherytypes.MasheryService {
	rgxArr := toRegexpArray(exp)

	filtered := []masherytypes.MasheryService{}
	for _, v := range *rv {
		for _, rgx := range rgxArr {
			if rgx.MatchString(v.Name) {
				filtered = append(filtered, v)
				break
			}
		}
	}

	return filtered
}

func dataSourceMasheryService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMasheryServiceRead,
		Schema:      DataSourceMashSvcSchema,
	}
}
