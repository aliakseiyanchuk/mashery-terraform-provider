package mashschema_test

import (
	"fmt"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestV3DomainsToResourceData(t *testing.T) {
	d := mashschema.DomainsMapper.TestResourceData()

	diags := mashschema.DomainsMapper.PersistTyped([]string{"a", "b", "c"}, d)
	LogErrorDiagnostics(t, "Parsing domains", &diags)

	fmt.Println(d.Get(mashschema.MashDomains))
}
