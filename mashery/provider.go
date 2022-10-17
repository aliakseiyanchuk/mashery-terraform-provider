package mashery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"strconv"
	"time"
)

// Provider mashschema configuration settings
//
// The provider accepts three types of access credentials:
// - the token value passed to either from Vault secret engine, or by reading  <code>TF_MASHERY_V3_ACCESS_TOKEN</code>
//   environment variable;
// - by specifying a path to the token file created previously with the <code>mash-connect</code> command. The token
//   file should be valid and allow at least 5 minutes of validity.
// - by specifying the path to the file where the long-term credentials are saved.
//
// Due to the nature of the credentials required for Mashery V3, the long-term credentials cannot appear in the
// Terraform state in the clear. V3 credentials allow full access to the API configuration as well as it allow
// retrieving all package key/secret combination.

const (
	envVaultToken = "TF_MASHERY_VAULT_TOKEN"
	envV3Token    = "TF_MASHERY_V3_ACCESS_TOKEN"
	envV3QPS      = "TF_MASHERY_V3_QPS"
	envV3Latency  = "TF_MASHERY_V3_NETWORK_LATENCY"

	vaultAddrField      = "vault_addr"
	vaultMountPathField = "vault_mount"
	vaultRoleField      = "vault_role"
	vaultTokenField     = "vault_token"

	providerQPSField            = "qps"
	providerNetworkLatencyField = "network_latency"
	providerV3Token             = "v3_token"
)

var providerConfigSchema = map[string]*schema.Schema{
	"log_file": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Log file where detailed Mashery session information will be saved",
	},
	vaultAddrField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Vault server address that will function as a V3 proxy",
		Default:     "",
	},
	vaultMountPathField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Vault server mount path",
		Default:     "mash-creds",
	},
	vaultRoleField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Role name to use with this provider",
		Default:     "",
	},
	vaultTokenField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Vault token that the provider should use to authenticate to Vault",
		DefaultFunc: schema.EnvDefaultFunc(envVaultToken, ""),
	},
	providerQPSField: {
		Type:        schema.TypeInt,
		Optional:    true,
		DefaultFunc: intValueFromVariable(envV3QPS, 2),
		Description: "Queries per second to observe. Default to 2 queries per second",
	},
	providerNetworkLatencyField: {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envV3Latency, "173ms"),
		Description: "Mean travel time between machine where the Terraform is running and Mashery API. Defaults to 173 (milliseconds).",
	},
	"token": {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envV3Token, ""),
		Description: "Actual access token to be used. For best use, obtain this from Vault",
	},
}

// ----------------------------------------------------------------------------------------------
// Vault proxy mode configuration

type VaultProxyModeConfiguration struct {
	addr     string
	mount    string
	roleName string
	token    string
}

func (vpm *VaultProxyModeConfiguration) isComplete() bool {
	return len(vpm.addr) > 0 && len(vpm.mount) > 0 && len(vpm.roleName) > 0 && len(vpm.token) > 0
}

func (vpm *VaultProxyModeConfiguration) fullAddress() string {
	return fmt.Sprintf("%s/v1/%s/roles/%s/proxy/v3", vpm.addr, vpm.mount, vpm.roleName)
}

// --------------------------------------------------------------------------------------------
// Vault authorizer

type VaultAuthorizer struct {
	transport.Authorizer

	vaultAuth map[string]string
}

func (va VaultAuthorizer) HeaderAuthorization() (map[string]string, error) {
	return va.vaultAuth, nil
}
func (va VaultAuthorizer) QueryStringAuthorization() (map[string]string, error) {
	return nil, nil
}

func (va VaultAuthorizer) Close() {
	// Do nothing
}

// -----------------------------------------------------------------------------
// Implementation

func intValueFromVariable(envVar string, defaultVal int) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(envVar); v != "" {
			return strconv.ParseInt(v, 10, 0)
		} else {
			return defaultVal, nil
		}
	}
}

var logger *log.Logger
var encoder *json.Encoder

// Send a message to the log file if it exists
func doLogf(format string, params ...interface{}) {
	if logger != nil {
		logger.Printf(format, params...)
	}
}

func doLogJson(msg string, obj interface{}) {
	if logger != nil {
		logger.Println(msg)
		if obj != nil {
			if b, err := json.Marshal(obj); err != nil {
				logger.Println(err.Error())
			} else {
				logger.Println(string(b))
			}
		} else {
			logger.Print("NULL JSON")
		}
	}
}

func vaultProxyConfiguration(d *schema.ResourceData) VaultProxyModeConfiguration {
	rv := VaultProxyModeConfiguration{
		addr:     d.Get(vaultAddrField).(string),
		mount:    d.Get(vaultMountPathField).(string),
		roleName: d.Get(vaultRoleField).(string),
		token:    d.Get(vaultTokenField).(string),
	}

	return rv
}

func transportLogging(ctx context.Context, wrq *transport.WrappedRequest, wrs *transport.WrappedResponse, err error) {
	doLogf("-> %s %s", wrq.Request.Method, wrq.Request.URL)
	for k, v := range wrq.Request.Header {
		doLogf("H> %s = %s", k, v)
	}
	if wrq.Body != nil {
		bodyOut := wrq.Body

		if str, err := json.MarshalIndent(wrq.Body, "|>", "  "); err == nil {
			bodyOut = str
		}

		doLogf("B>\n%s", bodyOut)
	}

	if wrs != nil {
		doLogf("<- %d", wrs.StatusCode)
		for k, v := range wrs.Header {
			doLogf("<H %s = %s", k, v)
		}

		if body, err := wrs.Body(); err != nil {
			doLogf("<H Can't read body: %s", err.Error())
		} else if len(body) > 0 {
			doLogf("<H Response body:%s\n", string(body))
		}

	}

	if err != nil {
		doLogf("Error: %s", err.Error())
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics
	logFile := d.Get("log_file").(string)
	if len(logFile) > 0 {
		encoder = new(json.Encoder)
		encoder.SetIndent("", "  ")

		now := time.Now()
		f, _ := os.Create(fmt.Sprintf("%s_%d%d%d_%d%d%d.log", logFile, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))
		logger = log.New(f, "TF_MASHERY :", log.LstdFlags)
	}

	var tokenProvider v3client.V3AccessTokenProvider
	qps := d.Get(providerQPSField).(int)

	// Prefer to use Vault proxy mode, if sufficiently configured.
	if vaultProxyMode := vaultProxyConfiguration(d); vaultProxyMode.isComplete() {
		doLogf("Provider was initialized with the Vault proxy mode, proxy=%s", vaultProxyMode.fullAddress())

		clParams := v3client.Params{
			HTTPClientParams: transport.HTTPClientParams{
				// Since the connection is made to Vault, the client will trust whatever the system
				// is trusting.
				TLSConfigDelegateSystem: true,
				ExchangeListener:        transportLogging,
			},
			Authorizer: VaultAuthorizer{
				vaultAuth: map[string]string{
					"X-Vault-Token": vaultProxyMode.token,
				},
			},
			QPS:          int64(qps),
			MashEndpoint: vaultProxyMode.fullAddress(),
		}

		cl := v3client.NewHttpClient(clParams)
		return cl, diags
	} else {
		doLogf("Provider configuration does not meet Vault proxy mode requirements")
	}

	// If Vault proxy mode is not configured, then Mashery V3 token should be supplied

	cfgNetLatency := d.Get(providerNetworkLatencyField).(string)

	travelComp, latErr := time.ParseDuration(cfgNetLatency)

	if latErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Not a valid network latency: %s (%s)", cfgNetLatency, latErr),
		})
	}

	if tknRaw, ok := d.GetOk(providerV3Token); ok {
		tkn := tknRaw.(string)
		// Mashery tokens is 24 characters long. This condition is a built-in protection for the
		// developer passing invalid tokens
		if len(tkn) > 20 {
			tokenProvider = v3client.NewFixedTokenProvider(tkn)
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("insufficient V3 token for connecting directly (%d chars supplied)", len(tkn)),
			})
		}
	}

	var cl v3client.Client

	if len(diags) == 0 {
		clParams := v3client.Params{
			Authorizer:    tokenProvider,
			QPS:           int64(qps),
			AvgNetLatency: travelComp,
		}

		cl = v3client.NewHttpClient(clParams)
		doLogf("Provider is initialized with explicitly supplied token")
	} else {
		doLogf("WARN: no suitable provider authentication methods exist")
	}

	return cl, diags
}

// Provider Mashery Terraform Provider mashschema definition
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: providerConfigSchema,
		ResourcesMap: map[string]*schema.Resource{
			"mashery_service":                      resourceMasheryService(),
			"mashery_service_error_set":            resourceMasheryErrorSet(),
			"mashery_processor_chain":              resourceMasheryProcessorChain(),
			"mashery_endpoint":                     EndpointResource.TFDataSourceSchema(),
			"mashery_endpoint_method":              EndpointMethodResource.TFDataSourceSchema(),
			"mashery_endpoint_method_filter":       EndpointMethodFilterResponse.TFDataSourceSchema(),
			"mashery_package":                      PackageResource.TFDataSourceSchema(),
			"mashery_package_plan":                 PackagePlanResource.TFDataSourceSchema(),
			"mashery_package_plan_service":         PackagePlanServiceResource.TFDataSourceSchema(),
			"mashery_package_plan_endpoint":        PackagePlanServiceEndpointResource.TFDataSourceSchema(),
			"mashery_package_plan_endpoint_method": resourceMasheryPlanMethod(),
			"mashery_member":                       MemberResource.TFDataSourceSchema(),
			"mashery_application":                  ApplicationResource.TFDataSourceSchema(),
			"mashery_package_key":                  PackageKeyResource.TFDataSourceSchema(),
			"mashery_unique_path":                  resourceMasheryUniquePath(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mashery_system_domains":     systemDomainsDataSource.TFDataSourceSchema(),
			"mashery_public_domains":     publicDomainsDataSource.TFDataSourceSchema(),
			"mashery_email_template_set": emailTemplateSet.TFDataSourceSchema(),
			"mashery_role":               roleDataSource.TFDataSourceSchema(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
