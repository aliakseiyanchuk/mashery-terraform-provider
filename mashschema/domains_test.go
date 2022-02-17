package mashschema_test

import (
	"context"
	"fmt"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestV3DomainsToResourceData(t *testing.T) {
	d := mashschema.DomainsMapper.NewResourceData()

	diags := mashschema.DomainsMapper.PersistTyped(context.TODO(), []string{"a", "b", "c"}, d)
	LogErrorDiagnostics(t, "Parsing domains", &diags)

	fmt.Println(d.Get(mashschema.MashDomains))
}
