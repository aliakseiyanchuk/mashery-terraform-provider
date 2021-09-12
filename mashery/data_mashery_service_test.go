package mashery_test

import (
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"terraform-provider-mashery/mashery"
	"testing"
)

func sampleMasheryServiceData() []v3client.MasheryService {
	return []v3client.MasheryService{
		{
			AddressableV3Object: v3client.AddressableV3Object{
				Id:   "A",
				Name: "A1",
			},
		},
		{
			AddressableV3Object: v3client.AddressableV3Object{
				Id:   "B",
				Name: "B1",
			},
		},
		{
			AddressableV3Object: v3client.AddressableV3Object{
				Id:   "C",
				Name: "C1",
			},
		},
	}
}

func TestDataSourceServiceProcess_DefaultNoMatch(t *testing.T) {
	var data []v3client.MasheryService
	d := NewResourceData(&mashery.DataSourceMashSvcSchema)

	diags := mashery.DataSourceSvcReadReturned(data, d)
	assertHasDiagnostic(t, "Required V3 service not found", &diags)
}

func TestDataSourceServiceProcess_WithMultipleResults(t *testing.T) {
	d := NewResourceData(&mashery.DataSourceMashSvcSchema)
	_ = d.Set(mashery.MashDataSourceUnique, true)

	diags := mashery.DataSourceSvcReadReturned(sampleMasheryServiceData(), d)
	assertHasDiagnostic(t, "Multiple matches", &diags)
}

func TestDataSourceServiceProcess_MultipleResults(t *testing.T) {
	d := NewResourceData(&mashery.DataSourceMashSvcSchema)
	diags := mashery.DataSourceSvcReadReturned(sampleMasheryServiceData(), d)
	LogErrorDiagnostics(t, "multiple service data source results", &diags)

	str := mashery.ExtractStringArray(d, mashery.MashSvcMultiRef, &mashery.EmptyStringArray)
	if len(str) != 3 {
		t.Errorf("incorrect length: expected 3, got %d", len(str))
	}
}

func TestFilterMasherySvcName(t *testing.T) {
	rv := sampleMasheryServiceData()

	filtered := mashery.FilterMasherySvcName(&rv, []string{"A.*", "B.*"})
	if len(filtered) != 2 {
		t.Errorf("Filtering by name should have located 2 elements")
	}
}
