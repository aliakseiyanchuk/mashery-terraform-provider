package mashery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"strconv"
	"strings"
	"terraform-provider-mashery/mashschema"
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
	envVaultAddress = "VAULT_ADDR" // Re-use vault integration

	envVaultToken    = "TF_MASHERY_VAULT_TOKEN"
	envVaultMount    = "TF_MASHERY_VAULT_MOUNT"
	envVaultRole     = "TF_MASHERY_VAULT_ROLE"
	envV3Token       = "TF_MASHERY_V3_ACCESS_TOKEN"
	envV3QPS         = "TF_MASHERY_QPS"
	envV3Latency     = "TF_MASHERY_NETWORK_LATENCY"
	envCacheServer   = "TF_MASHERY_REDIS_SERVER"
	envCacheDuration = "TF_MASHERY_CACHE_DURATION"

	vaultAddrField      = "vault_addr"
	vaultMountPathField = "vault_mount"
	engineRoleField     = "role"
	vaultTokenField     = "vault_token"
	vaultProxyMode      = "vault_proxy_mode"

	providerQPSField            = "qps"
	providerNetworkLatencyField = "network_latency"
	providerV3Token             = "v3_token"

	ProviderRedisCacheField = "redis_url"
	ProviderCacheDuration   = "cache_duration"
)

var ProviderConfigSchema = map[string]*schema.Schema{
	vaultAddrField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Vault server address that will function as a V3 proxy",
		DefaultFunc: schema.EnvDefaultFunc(envVaultAddress, ""),
	},
	vaultMountPathField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Vault server mount path",
		DefaultFunc: schema.EnvDefaultFunc(envVaultMount, "mash-auth"),
	},
	engineRoleField: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Role name to use with this provider",
		Default:     schema.EnvDefaultFunc(envVaultRole, "mash-auth"),
	},
	vaultProxyMode: {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Whether the provider should operate in Vault proxy mode (i.e. delegate all operations to Vault)",
		Default:     false,
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
		DefaultFunc: initValueFromVariable(envV3QPS, 1),
		Description: "Queries per second to observe. Default to 1 queries per second",
	},
	providerNetworkLatencyField: {
		Type:             schema.TypeString,
		Optional:         true,
		DefaultFunc:      schema.EnvDefaultFunc(envV3Latency, "173ms"),
		ValidateDiagFunc: mashschema.ValidateDuration,
		Description:      "Mean travel time between machine where the Terraform is running and Mashery API. Defaults to 173 (milliseconds).",
	},
	providerV3Token: {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envV3Token, ""),
		Description: "Actual access token to be used. For best use, obtain this from Vault",
	},
	ProviderRedisCacheField: {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envCacheServer, ""),
		Description: "URL to the Redis cache server and database to use",
	},
	ProviderCacheDuration: {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envCacheDuration, "24h"),
		Description: "Duration of data source caches; defaults to 24h if not specified",
	},
}

// ----------------------------------------------------------------------------------------------
// Vault proxy mode configuration

type VaultPairingConfiguration struct {
	addr      string
	mount     string
	roleName  string
	proxyMode bool
	token     string
}

func (vpm *VaultPairingConfiguration) isCompleteForProxyMode() bool {
	return len(vpm.addr) > 0 && len(vpm.mount) > 0 && len(vpm.roleName) > 0 && len(vpm.token) > 0 && vpm.proxyMode
}

func (vpm *VaultPairingConfiguration) isCompleteForResourceFetch() bool {
	return len(vpm.addr) > 0 && len(vpm.mount) > 0 && len(vpm.roleName) > 0 && len(vpm.token) > 0 && !vpm.proxyMode
}

func (vpm *VaultPairingConfiguration) fullAddress() string {
	return fmt.Sprintf("%s/v1/%s/roles/%s/proxy/v3", vpm.addr, vpm.mount, vpm.roleName)
}

func (vpm *VaultPairingConfiguration) tokenAddress() string {
	return fmt.Sprintf("%s/v1/%s/roles/%s/token", vpm.addr, vpm.mount, vpm.roleName)
}

func (vpm *VaultPairingConfiguration) vaultToken() transport.VaultToken {
	return transport.VaultToken(vpm.token)
}

// -----------------------------------------------------------------------------
// Implementation

func initValueFromVariable(envVar string, defaultVal int) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(envVar); v != "" {
			return strconv.ParseInt(v, 10, 0)
		} else {
			return defaultVal, nil
		}
	}
}

// Send a message to the log file if it exists

func vaultPairingConfiguration(d *schema.ResourceData) VaultPairingConfiguration {
	rv := VaultPairingConfiguration{
		addr:      d.Get(vaultAddrField).(string),
		mount:     d.Get(vaultMountPathField).(string),
		roleName:  d.Get(engineRoleField).(string),
		proxyMode: d.Get(vaultProxyMode).(bool),
		token:     d.Get(vaultTokenField).(string),
	}

	if tknFromEnv := os.Getenv(envVaultToken); len(tknFromEnv) > 0 {
		rv.token = tknFromEnv
	}

	return rv
}

func transportLogging(ctx context.Context, wrq *transport.WrappedRequest, wrs *transport.WrappedResponse, err error) {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("-> %s %s\n", wrq.Request.Method, wrq.Request.URL))
	for k, v := range wrq.Request.Header {
		// The Authorization and X-Vault-Token header value is never written into the logs.
		// All other headers are written into the logs with their value as-is.
		switch strings.ToLower(k) {
		case "authorization":
		case "x-vault-token":
			b.WriteString(fmt.Sprintf("H> %s = %s\n", k, "****<REDACTED>****"))
		default:
			b.WriteString(fmt.Sprintf("H> %s = %s\n", k, v))
		}
	}
	if wrq.Body != nil {
		bodyOut := wrq.Body

		if str, err := json.MarshalIndent(wrq.Body, "|>", "  "); err == nil {
			bodyOut = str
		}

		b.WriteString(fmt.Sprintf("B>\n%s\n", bodyOut))
	}

	if wrs != nil {
		b.WriteString(fmt.Sprintf("<- %d\n", wrs.StatusCode))
		for k, v := range wrs.Header {
			b.WriteString(fmt.Sprintf("<H %s = %s\n", k, v))
		}

		if body, err := wrs.Body(); err != nil {
			b.WriteString(fmt.Sprintf("<H Can't read body: %s\n", err.Error()))
		} else if len(body) > 0 {
			b.WriteString(fmt.Sprintf("<H Response body:%s\n", string(body)))
		}

	}

	tflog.Debug(ctx, b.String())

	if err != nil {
		b.WriteString(fmt.Sprintf("[Request Error] %s\n", err.Error()))
		tflog.Warn(ctx, b.String())
	} else {
		tflog.Trace(ctx, b.String())
	}

}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	var tokenProvider v3client.V3AccessTokenProvider
	qps := d.Get(providerQPSField).(int)

	requestedLatencyCompensation := mashschema.ExtractString(d, providerNetworkLatencyField, "173ms")
	netLatency, err := time.ParseDuration(requestedLatencyCompensation)

	tflog.Info(ctx, fmt.Sprintf("Requested observed QPS: %d", qps))
	tflog.Info(ctx, fmt.Sprintf("Requested network latency compensation: %s", netLatency))

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Network latency compensation is not valid: %s", err.Error()))

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "invalid network latency format",
			Detail:   fmt.Sprintf("network compensation value must be a valid Go time format. Supplied value %s is not valid: %s", requestedLatencyCompensation, err.Error()),
		})
	}

	// Prefer to use Vault proxy mode, if sufficiently configured.
	if vaultPairingCfg := vaultPairingConfiguration(d); vaultPairingCfg.isCompleteForProxyMode() {
		clParams := v3client.Params{
			HTTPClientParams: transport.HTTPClientParams{
				// Since the connection is made to Vault, the client will trust whatever the system
				// is trusting.
				TLSConfigDelegateSystem: true,
				ExchangeListener:        transportLogging,
			},
			Authorizer:    transport.NewVaultAuthorizer(vaultPairingCfg.vaultToken()),
			QPS:           int64(qps),
			AvgNetLatency: netLatency,
			MashEndpoint:  vaultPairingCfg.fullAddress(),
		}

		cl := v3client.NewHttpClient(clParams)

		tflog.Info(ctx, fmt.Sprintf("Provider initialized with the Vault *proxy* mode, proxy=%s", vaultPairingCfg.fullAddress()))
		return cl, diags
	} else if vaultPairingCfg.isCompleteForResourceFetch() {
		clParams := v3client.Params{
			HTTPClientParams: transport.HTTPClientParams{
				TLSConfigDelegateSystem: true,
				ExchangeListener:        transportLogging,
			},
			Authorizer:    transport.NewVaultTokenResourceAuthorizer(vaultPairingCfg.tokenAddress(), vaultPairingCfg.vaultToken()),
			QPS:           int64(qps),
			AvgNetLatency: netLatency,
			// No Mashery endpoint override is necessary here: the connection is established
			// directly to the Mashery API
		}

		cl := v3client.NewHttpClient(clParams)

		tflog.Info(ctx, fmt.Sprintf("Provider initialized with the Vault *token fetch* mode, token endpoint=%s", vaultPairingCfg.tokenAddress()))
		return cl, diags
	} else {
		tflog.Info(ctx, "Provider configuration does not meet Vault proxy mode requirements")
	}

	if tknRaw, ok := d.GetOk(providerV3Token); ok {
		tkn := tknRaw.(string)
		// Mashery tokens is 24 characters long. This condition is a built-in protection for the
		// developer passing invalid tokens
		if len(tkn) > 20 {
			tokenProvider = v3client.NewFixedTokenProvider(tkn)
			tflog.Info(ctx, "Provider is initialized with explicitly supplied token.")
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("insufficient V3 token for connecting directly (%d chars supplied)", len(tkn)),
			})
		}
	}

	if tokenProvider == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("no suitable methods provided to authenticate to Mashery API"),
		})
	}

	var cl v3client.Client

	if len(diags) == 0 {
		clParams := v3client.Params{
			Authorizer:    tokenProvider,
			QPS:           int64(qps),
			AvgNetLatency: netLatency,
			HTTPClientParams: transport.HTTPClientParams{
				ExchangeListener: transportLogging,
			},
		}

		cl = v3client.NewHttpClientWithBadRequestAutoRetries(clParams)
	}

	tflog.Info(ctx, "Provider initialized completed", map[string]interface{}{
		"diagnostic_count": len(diags),
		"diagnostic_error": diags.HasError(),
	})
	return cl, diags
}
