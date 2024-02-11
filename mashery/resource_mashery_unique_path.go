package mashery

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"time"
)

const (
	MashUniquePathPrefix = "prefix"
	MashUniquePath       = "path"
)

func ResourceMasheryUniquePath() *schema.Resource {
	return &schema.Resource{
		CreateContext: createUniquePath,
		ReadContext:   schema.NoopContext,
		DeleteContext: deleteUniquePath,
		Schema: map[string]*schema.Schema{
			MashUniquePathPrefix: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Path prefix",
			},
			MashUniquePath: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

var valueRage []rune

func init() {
	var i rune
	for i = 'a'; i <= 'z'; i++ {
		valueRage = append(valueRage, i)
	}

	for i := 'A'; i <= 'Z'; i++ {
		valueRage = append(valueRage, i)
	}

	for i := '0'; i <= '9'; i++ {
		valueRage = append(valueRage, i)
	}
}

func CreateUniquePath(prefix string, refTime int64) string {
	val := refTime
	rng := int64(len(valueRage))

	rv := strings.Builder{}
	rv.WriteString("/")

	if len(prefix) > 0 {
		rv.WriteString(prefix)
		rv.WriteString("_")
	}

	for val > rng {
		rv.WriteRune(valueRage[(val % rng)])
		val = val / rng
	}

	rv.WriteRune(valueRage[val])

	return rv.String()
}

func createUniquePath(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	if _, exists := data.GetOk(MashUniquePath); !exists {
		key := time.Now().Unix()

		v := CreateUniquePath(data.Get(MashUniquePathPrefix).(string), key)
		_ = data.Set(MashUniquePath, v)

		data.SetId(fmt.Sprintf("path-%d", key))
	}

	return diag.Diagnostics{}
}

func deleteUniquePath(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	data.SetId("")
	return diag.Diagnostics{}
}
