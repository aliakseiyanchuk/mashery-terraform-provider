package mashery_test

import (
	"fmt"
	"terraform-provider-mashery/mashery"
	"testing"
)

func TestV3DomainsToResourceData(t *testing.T) {
	d := NewResourceData(&mashery.DomainsSchema)

	diags := mashery.V3DomainsToResourceData([]string{"a", "b", "c"}, d)
	LogErrorDiagnostics(t, "Parsing domains", &diags)

	fmt.Println(d.Get(mashery.MashDomains))
}
