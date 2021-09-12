package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"terraform-provider-mashery/mashery"
	"testing"
)

func assertNotNil(t *testing.T, msg string, obj interface{}) {
	if obj == nil {
		t.Errorf("Nil object received: %s", msg)
	}
}

func assertSchemaType(t *testing.T, msg string, expType schema.ValueType, actualType schema.ValueType) {
	if expType != actualType {
		t.Errorf("Unexpected type for %s: expected %s, got %s", msg, expType, actualType)
	}
}

func assertComputed(t *testing.T, msg string, sch *schema.Schema) {
	if !sch.Computed {
		t.Errorf("%s is not computed", msg)
	}
}

func assertResource(t *testing.T, elem string, inp interface{}) *schema.Resource {
	if rv, ok := inp.(*schema.Resource); ok {
		return rv
	} else {
		t.Errorf("Not a resource: %s", elem)
		return nil
	}
}

func assertSameString(t *testing.T, key string, orig, ref *string) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if *orig != *ref {
		t.Errorf("%s != %s for %s", *orig, *ref, key)
	}
}

func assertSameStringArray(t *testing.T, key string, orig, ref []string) {
	if !reflect.DeepEqual(orig, ref) {
		t.Errorf("Srings arrays differ for %s", key)
	}
}

func assertSameEAV(t *testing.T, key string, orig, ref *v3client.EAV) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if !reflect.DeepEqual(orig, ref) {
		t.Errorf("eavs are not the same for %s", key)
	}
}

func assertDeepEqualLimits(t *testing.T, key string, orig, ref *[]v3client.Limit) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if !reflect.DeepEqual(orig, ref) {
		t.Errorf("eavs are not the same for %s", key)
	}
}

func assertSameSet(t *testing.T, key string, orig, ref *[]string) {
	if len(*orig) != len(*ref) {
		t.Errorf("Aray length is not the same: %d != %d", len(*orig), len(*ref))
	}

	for _, v := range *orig {
		if mashery.FindInArray(v, ref) < 0 {
			t.Errorf("Missing in ref: %s", v)
		}
	}

	for _, v := range *ref {
		if mashery.FindInArray(v, orig) < 0 {
			t.Errorf("Missing in origin: %s", v)
		}
	}
}

func assertDeepEqual(t *testing.T, key string, orig, ref interface{}) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if !reflect.DeepEqual(orig, ref) {
		t.Errorf("arbitrary objects are not the same for %s", key)
	}
}

func assertSameBool(t *testing.T, key string, orig, ref *bool) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if *orig != *ref {
		t.Errorf("%t != %t for %s", *orig, *ref, key)
	}
}

func assertSameInt64(t *testing.T, key string, orig, ref *int64) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if *orig != *ref {
		t.Errorf("%d != %d for %s", *orig, *ref, key)
	}
}

func assertSameInt(t *testing.T, key string, orig, ref *int) {
	if orig == nil && ref == nil {
		return
	} else if (orig != nil && ref == nil) || (orig == nil && ref != nil) {
		t.Errorf("%s nil mismatcch", key)
	} else if *orig != *ref {
		t.Errorf("%d != %d for %s", *orig, *ref, key)
	}
}

func assertResourceHasStringKey(t *testing.T, d *schema.ResourceData, key, value string) {
	extr := mashery.ExtractStringPointer(d, key)
	if extr == nil {
		t.Errorf("No value found for key %s", key)
	} else if *extr != value {
		t.Errorf("%s != %s for key %s", *extr, value, key)
	}
}

func assertResourceDoesNotHaveKey(t *testing.T, d *schema.ResourceData, key string) {
	if _, ok := d.GetOk(key); ok {
		t.Errorf("Unexpected presence of key %s", key)
	}
}

func assertResourceContainsKey(t *testing.T, d *schema.ResourceData, key string) {
	if _, ok := d.GetOk(key); !ok {
		t.Errorf("Unexpected absence of key %s", key)
	}
}

func assertOk(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assertHasDiagnostic(t *testing.T, summary string, dgs *diag.Diagnostics) {
	for _, v := range *dgs {
		if summary == v.Summary {
			return
		}
	}

	t.Errorf("no diagnostic with summary %s found", summary)
}
