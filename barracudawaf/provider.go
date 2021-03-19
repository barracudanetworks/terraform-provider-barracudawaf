package barracudawaf

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Provider : Schema definition for barracudawaf provider
func Provider() *schema.Provider {

	// The actual provider
	provider := &schema.Provider{

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP Address of the WAF to be configured",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Admin port on the WAF to be configured",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of the WAF to be configured",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password of the WAF to be configured",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"barracudawaf_form_spam_forms":             resourceCudaWAFFormSpamForms(),
			"barracudawaf_trusted_ca_certificate":      resourceCudaWAFTrustedCaCertificate(),
			"barracudawaf_json_key_profiles":           resourceCudaWAFJsonKeyProfiles(),
			"barracudawaf_users":                       resourceCudaWAFUsers(),
			"barracudawaf_http_response_rewrite_rules": resourceCudaWAFHttpResponseRewriteRules(),
			"barracudawaf_local_groups":                resourceCudaWAFLocalGroups(),
			"barracudawaf_protected_data_types":        resourceCudaWAFProtectedDataTypes(),
			"barracudawaf_local_users":                 resourceCudaWAFLocalUsers(),
			"barracudawaf_content_rules":               resourceCudaWAFContentRules(),
			"barracudawaf_saml_services":               resourceCudaWAFSamlServices(),
			"barracudawaf_response_body_rewrite_rules": resourceCudaWAFResponseBodyRewriteRules(),
			"barracudawaf_network_acls":                resourceCudaWAFNetworkAcls(),
			"barracudawaf_attack_patterns":             resourceCudaWAFAttackPatterns(),
			"barracudawaf_trusted_server_certificate":  resourceCudaWAFTrustedServerCertificate(),
			"barracudawaf_secure_browsing_policies":    resourceCudaWAFSecureBrowsingPolicies(),
			"barracudawaf_external_radius_services":    resourceCudaWAFExternalRadiusServices(),
			"barracudawaf_url_profiles":                resourceCudaWAFUrlProfiles(),
			"barracudawaf_syslog_servers":              resourceCudaWAFSyslogServers(),
			"barracudawaf_services":                    resourceCudaWAFServices(),
			"barracudawaf_url_acls":                    resourceCudaWAFUrlAcls(),
			"barracudawaf_http_request_rewrite_rules":  resourceCudaWAFHttpRequestRewriteRules(),
			"barracudawaf_bonds":                       resourceCudaWAFBonds(),
			"barracudawaf_input_patterns":              resourceCudaWAFInputPatterns(),
			"barracudawaf_radius_services":             resourceCudaWAFRadiusServices(),
			"barracudawaf_web_scraping_policies":       resourceCudaWAFWebScrapingPolicies(),
			"barracudawaf_action_policies":             resourceCudaWAFActionPolicies(),
			"barracudawaf_external_ldap_services":      resourceCudaWAFExternalLdapServices(),
			"barracudawaf_export_configuration":        resourceCudaWAFExportConfiguration(),
			"barracudawaf_kerberos_services":           resourceCudaWAFKerberosServices(),
			"barracudawaf_header_acls":                 resourceCudaWAFHeaderAcls(),
			"barracudawaf_adaptive_profiling_rules":    resourceCudaWAFAdaptiveProfilingRules(),
			"barracudawaf_authorization_policies":      resourceCudaWAFAuthorizationPolicies(),
			"barracudawaf_global_acls":                 resourceCudaWAFGlobalAcls(),
			"barracudawaf_parameter_profiles":          resourceCudaWAFParameterProfiles(),
			"barracudawaf_configuration_checkpoints":   resourceCudaWAFConfigurationCheckpoints(),
			"barracudawaf_identity_types":              resourceCudaWAFIdentityTypes(),
			"barracudawaf_ldap_services":               resourceCudaWAFLdapServices(),
			"barracudawaf_service_groups":              resourceCudaWAFServiceGroups(),
			"barracudawaf_rsa_securid_services":        resourceCudaWAFRsaSecuridServices(),
			"barracudawaf_vlans":                       resourceCudaWAFVlans(),
			"barracudawaf_bot_spam_patterns":           resourceCudaWAFBotSpamPatterns(),
			"barracudawaf_openidc_identity_providers":  resourceCudaWAFOpenidcIdentityProviders(),
			"barracudawaf_bot_spam_types":              resourceCudaWAFBotSpamTypes(),
			"barracudawaf_administrator_roles":         resourceCudaWAFAdministratorRoles(),
			"barracudawaf_content_rule_servers":        resourceCudaWAFContentRuleServers(),
			"barracudawaf_parameter_optimizers":        resourceCudaWAFParameterOptimizers(),
			"barracudawaf_url_optimizers":              resourceCudaWAFUrlOptimizers(),
			"barracudawaf_rate_control_pools":          resourceCudaWAFRateControlPools(),
			"barracudawaf_whitelisted_bots":            resourceCudaWAFWhitelistedBots(),
			"barracudawaf_ddos_policies":               resourceCudaWAFDdosPolicies(),
			"barracudawaf_identity_theft_patterns":     resourceCudaWAFIdentityTheftPatterns(),
			"barracudawaf_preferred_clients":           resourceCudaWAFPreferredClients(),
			"barracudawaf_network_interfaces":          resourceCudaWAFNetworkInterfaces(),
			"barracudawaf_trusted_host_groups":         resourceCudaWAFTrustedHostGroups(),
			"barracudawaf_reports":                     resourceCudaWAFReports(),
			"barracudawaf_input_types":                 resourceCudaWAFInputTypes(),
			"barracudawaf_security_policies":           resourceCudaWAFSecurityPolicies(),
			"barracudawaf_auto_system_acls":            resourceCudaWAFAutoSystemAcls(),
			"barracudawaf_geo_pools":                   resourceCudaWAFGeoPools(),
			"barracudawaf_client_certificate_crls":     resourceCudaWAFClientCertificateCrls(),
			"barracudawaf_signed_certificate":          resourceCudaWAFSignedCertificate(),
			"barracudawaf_trusted_hosts":               resourceCudaWAFTrustedHosts(),
			"barracudawaf_openidc_services":            resourceCudaWAFOpenidcServices(),
			"barracudawaf_session_identifiers":         resourceCudaWAFSessionIdentifiers(),
			"barracudawaf_credential_servers":          resourceCudaWAFCredentialServers(),
			"barracudawaf_json_profiles":               resourceCudaWAFJsonProfiles(),
			"barracudawaf_access_rules":                resourceCudaWAFAccessRules(),
			"barracudawaf_url_encryption_rules":        resourceCudaWAFUrlEncryptionRules(),
			"barracudawaf_attack_types":                resourceCudaWAFAttackTypes(),
			"barracudawaf_url_policies":                resourceCudaWAFUrlPolicies(),
			"barracudawaf_self_signed_certificate":     resourceCudaWAFSelfSignedCertificate(),
			"barracudawaf_json_security_policies":      resourceCudaWAFJsonSecurityPolicies(),
			"barracudawaf_custom_parameter_classes":    resourceCudaWAFCustomParameterClasses(),
			"barracudawaf_servers":                     resourceCudaWAFServers(),
			"barracudawaf_response_pages":              resourceCudaWAFResponsePages(),
			"barracudawaf_saml_identity_providers":     resourceCudaWAFSamlIdentityProviders(),
			"barracudawaf_url_translations":            resourceCudaWAFUrlTranslations(),
			"barracudawaf_allow_deny_clients":          resourceCudaWAFAllowDenyClients(),
			"barracudawaf_vsites":                      resourceCudaWAFVsites(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		IPAddress: d.Get("address").(string),
		AdminPort: d.Get("port").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}
	cfg, err := config.Client()
	if err != nil {
		return cfg, err
	}
	cfg.UserAgent = fmt.Sprintf("Terraform/%s", terraformVersion)
	return cfg, err
}
