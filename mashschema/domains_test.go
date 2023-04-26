package mashschema_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"terraform-provider-mashery/mashschema"
	"testing"
)

func TestV3DomainsToResourceData(t *testing.T) {
	d := mashschema.DomainsMapper.TestResourceData()

	diags := mashschema.DomainsMapper.PersistTyped([]string{"a", "b", "c"}, d)
	LogErrorDiagnostics(t, "Parsing domains", &diags)

	fmt.Println(d.Get(mashschema.MashDomains))
}

func TestV3DomainsMapptingTerraformSchema(t *testing.T) {
	fmt.Println(mashschema.DomainsMapper.V3ObjectName())
	schema := mashschema.DomainsMapper.TerraformSchema()
	assert.NotNil(t, schema)
}
