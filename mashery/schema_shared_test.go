package mashery_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-mashery/mashery"
	"testing"
)

func TestParseCompoundId(t *testing.T) {
	v := mashery.CreateCompoundId("a1", "b2")
	var a, b string

	mashery.ParseCompoundId(v, &a, &b)
	if a != "a1" {
		t.Errorf("Incorrect value for a")
	}

	if b != "b2" {
		t.Errorf("Incorrect value for b")
	}
}

func TestParseCompoundIdWithComment(t *testing.T) {
	v := fmt.Sprintf("%s # comment", mashery.CreateCompoundId("a1", "b2"))
	var a, b string

	mashery.ParseCompoundId(v, &a, &b)
	if a != "a1" {
		t.Errorf("Incorrect value for a")
	}

	if b != "b2" {
		t.Errorf("Incorrect value for b")
	}
}

func NewResourceData(sch *map[string]*schema.Schema) *schema.ResourceData {
	res := schema.Resource{
		Schema: *sch,
	}

	return res.TestResourceData()
}

// Log error diagnostic of the test, so that the test would fail if any errors were
// encountered.
func LogErrorDiagnostics(t *testing.T, ctx string, diagnostics *diag.Diagnostics) {
	cnt := 0

	for _, d := range *diagnostics {
		if d.Severity == diag.Error {
			cnt++
			break
		}
	}

	if cnt > 0 {
		t.Errorf("Diagnostic of %s indicates %d errors", ctx, cnt)
		for _, v := range *diagnostics {
			severityLabel := "ERROR"
			if v.Severity == diag.Warning {
				severityLabel = "WARNING"
			}

			t.Errorf("%s: %s", severityLabel, v.Detail)
		}
	}
}

// Exchange a struct via resource data
func ExchangeViaResourceData(t *testing.T, sch *map[string]*schema.Schema, id string, writer func(d *schema.ResourceData) diag.Diagnostics, reader func(d *schema.ResourceData) interface{}) interface{} {
	res := schema.Resource{
		Schema: *sch,
	}

	d := res.TestResourceData()
	w := writer(d)
	if len(w) > 0 {
		t.Errorf("Setting encountered %d errors", len(w))
		for _, v := range w {
			t.Errorf("%s", v.Detail)
		}
	}
	d.SetId(id)

	return reader(d)
}
