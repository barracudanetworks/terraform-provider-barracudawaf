package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceServicesParams = map[string][]string{
		"authentication": {
			"send_domain_name",
			"dual_authentication",
			"secondary_authentication_service",
			"authentication_service",
			"status",
			"access_denied_page",
			"session_timeout_for_activesync",
			"cookie_domain",
			"cookie_path",
			"session_replay_protection_status",
			"creation_timeout",
			"dual_login_page",
			"idle_timeout",
			"login_challenge_page",
			"login_failed_page",
			"login_page",
			"login_processor_path",
			"openidc_redirect_url",
			"openidc_attribute_name",
			"openidc_local_id",
			"login_successful_page",
			"challenge_prompt_field",
			"challenge_user_field",
			"login_failure_url",
			"login_success_url",
			"logout_page",
			"logout_processor_path",
			"logout_successful_page",
			"logout_success_url",
			"password_expired_url",
			"post_processor_path",
			"saml_logout_url",
			"action",
			"groups",
			"sso_cookie_update_interval",
			"max_failed_attempts",
			"count_window",
			"enable_bruteforce_prevention",
			"kerberos_debug_status",
			"kerberos_enable_delegation",
			"kerberos_ldap_authorization",
			"krb_authorization_policy",
			"master_service",
			"master_service_url",
			"service_provider_display_name",
			"service_provider_entity_id",
			"service_provider_org_name",
			"service_provider_org_url",
			"attribute_format",
			"attribute_id",
			"attribute_name",
			"attribute_type",
			"encryption_certificate",
			"signing_certificate",
		},
		"caching": {
			"expiry_age",
			"file_extensions",
			"max_size",
			"min_size",
			"cache_negative_response",
			"ignore_request_headers",
			"ignore_response_headers",
			"status",
		},
		"clickjacking": {"allowed_origin", "options", "status"},
		"compression":  {"content_types", "min_size", "status", "unknown_content_types"},
		"exception_profiling": {
			"exception_profiling_trusted_host_group",
			"exception_profiling_learn_from_trusted_host",
			"exception_profiling_level",
		},
		"ftp_security": {
			"attack_prevention_status",
			"allowed_verbs",
			"allowed_verb_status",
			"pasv_ip_address",
			"pasv_ports",
		},
		"instant_ssl": {"status", "sharepoint_rewrite_support", "secure_site_domain"},
		"ip_reputation": {
			"anonymous_proxy",
			"barracuda_reputation_blocklist",
			"custom_blacklisted_ip_status",
			"datacenter_ip",
			"fake_crawler",
			"check_registered_country",
			"block_unclassified_ips",
			"apply_policy_at",
			"geo_pool",
			"geoip_action",
			"enable_ip_reputation_filter",
			"geoip_enable_logging",
			"known_http_attack_sources",
			"public_proxy",
			"satellite_provider",
			"known_ssh_attack_sources",
			"tor_nodes",
		},
		"adaptive_profiling": {
			"content_types",
			"ignore_parameters",
			"navigation_parameters",
			"request_learning",
			"response_learning",
			"status",
			"trusted_host_group",
		},
		"sensitive_parameter_names": {"sensitive_parameter_names"},
		"session_tracking":          {"identifiers", "exception_clients", "max_interval", "max_sessions_per_ip", "status"},
		"slow_client_attack": {
			"data_transfer_rate",
			"exception_clients",
			"incremental_request_timeout",
			"incremental_response_timeout",
			"max_request_timeout",
			"max_response_timeout",
			"status",
		},
		"website_profile": {
			"strict_profile_check",
			"allowed_domains",
			"exclude_url_patterns",
			"include_url_patterns",
			"mode",
			"use_profile",
		},
		"advanced_configuration": {
			"enable_web_application_firewall",
			"accept_list",
			"accept_list_status",
			"proxy_list",
			"proxy_list_status",
			"ddos_exception_list",
			"enable_fingerprint",
			"enable_http2",
			"keepalive_requests",
			"ntlm_ignore_extra_data",
			"enable_proxy_protocol",
			"enable_vdi",
			"enable_websocket",
		},
		"basic_security": {
			"web_firewall_log_level",
			"mode",
			"trusted_hosts_action",
			"trusted_hosts_group",
			"ignore_case",
			"client_ip_addr_header",
			"rate_control_pool",
			"rate_control_status",
			"web_firewall_policy",
		},
		"load_balancing": {
			"algorithm",
			"persistence_cookie_domain",
			"cookie_age",
			"persistence_cookie_name",
			"persistence_cookie_path",
			"failover_method",
			"header_name",
			"persistence_idle_timeout",
			"persistence_method",
			"source_ip_netmask",
			"parameter_name",
		},
		"ssl_client_authentication": {
			"client_certificate_for_rule",
			"client_authentication_rule_count",
			"client_authentication",
			"enforce_client_certificate",
			"trusted_certificates",
		},
		"ssl_security": {
			"certificate",
			"ciphers",
			"ecdsa_certificate",
			"include_hsts_sub_domains",
			"hsts_max_age",
			"create_hsts_redirect_service",
			"selected_ciphers",
			"override_ciphers_ssl3",
			"override_ciphers_tls_1_1",
			"override_ciphers_tls_1_2",
			"override_ciphers_tls_1_3",
			"override_ciphers_tls_1",
			"enable_pfs",
			"enable_ssl_3",
			"enable_tls_1",
			"enable_tls_1_1",
			"enable_tls_1_2",
			"enable_tls_1_3",
			"enable_hsts",
			"enable_ocsp_stapling",
			"sni_certificate",
			"domain",
			"sni_ecdsa_certificate",
			"enable_sni",
			"enable_strict_sni_check",
			"status",
			"ssl_tls_presets",
		},
		"captcha_settings":         {"recaptcha_type", "recaptcha_domain", "recaptcha_site_key", "recaptcha_site_secret"},
		"ssl_ocsp":                 {"enable", "responder_url", "certificate"},
		"url_encryption":           {"status"},
		"referer_spam":             {"exception_patterns", "custom_blocked_patterns", "status"},
		"comment_spam":             {"exception_patterns", "parameter"},
		"advanced_client_analysis": {"advanced_analysis", "exclude_url_patterns"},
		"form_spam":                {"status", "honeypot_status", "autoconfigure_status"},
		"waas_account":             {"waas_account_id", "waas_account_serial"},
	}
)

func resourceCudaWAFServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFServicesCreate,
		Read:   resourceCudaWAFServicesRead,
		Update: resourceCudaWAFServicesUpdate,
		Delete: resourceCudaWAFServicesDelete,

		Schema: map[string]*schema.Schema{
			"address_version": {Type: schema.TypeString, Optional: true, Description: "Version"},
			"dps_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable Active DDoS Prevention",
			},
			"mask":                {Type: schema.TypeString, Optional: true, Description: "Mask"},
			"session_timeout":     {Type: schema.TypeString, Optional: true, Description: "Session Timeout"},
			"linked_service_name": {Type: schema.TypeString, Optional: true},
			"enable_access_logs":  {Type: schema.TypeString, Optional: true, Description: "Enable Access Logs"},
			"app_id":              {Type: schema.TypeString, Optional: true, Description: "Service App Id"},
			"comments":            {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"group":               {Type: schema.TypeString, Optional: true, Description: "Service Group"},
			"service_id":          {Type: schema.TypeString, Optional: true},
			"ip_address":          {Type: schema.TypeString, Optional: true, Description: "VIP"},
			"cloud_ip_select":     {Type: schema.TypeString, Optional: true},
			"name":                {Type: schema.TypeString, Required: true, Description: "Web Application Name"},
			"port":                {Type: schema.TypeString, Optional: true, Description: "Port"},
			"status":              {Type: schema.TypeString, Optional: true, Description: "Status"},
			"type":                {Type: schema.TypeString, Optional: true, Description: "Type"},
			"certificate":         {Type: schema.TypeString, Optional: true},
			"service_hostname":    {Type: schema.TypeString, Optional: true},
			"vsite":               {Type: schema.TypeString, Optional: true, Description: "Vsite"},
			"authentication": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_domain_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Send Domain Name to RADIUS/RSA Server",
						},
						"dual_authentication": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: " Dual Authentication Required",
						},
						"secondary_authentication_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Secondary Authentication Service",
						},
						"authentication_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication Service",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
						"access_denied_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access Denied Page",
						},
						"session_timeout_for_activesync": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Session Timeout for ActiveSync",
						},
						"cookie_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Session Cookie Domain",
						},
						"cookie_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Session Cookie Path",
						},
						"session_replay_protection_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enforce strict session controls",
						},
						"creation_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creation Timeout",
						},
						"dual_login_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dual Login Page",
						},
						"idle_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Idle Timeout",
						},
						"login_challenge_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login Challenge Page",
						},
						"login_failed_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login Failed Page",
						},
						"login_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login Page",
						},
						"login_processor_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login Processor Path",
						},
						"openidc_redirect_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OpenID connect redirect url",
						},
						"openidc_attribute_name": {Type: schema.TypeString, Optional: true, Description: "None"},
						"openidc_local_id":       {Type: schema.TypeString, Optional: true, Description: "None"},
						"login_successful_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login Successful Page",
						},
						"challenge_prompt_field": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Challenge Prompt Field",
						},
						"challenge_user_field": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Challenge User Field",
						},
						"login_failure_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auth Failure URL",
						},
						"login_success_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auth Success URL",
						},
						"logout_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logout Page",
						},
						"logout_processor_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logout Processor Path",
						},
						"logout_successful_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logout Successful Page",
						},
						"logout_success_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auth Logout Success URL",
						},
						"password_expired_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auth Password Expired URL",
						},
						"post_processor_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Post Processor Path",
						},
						"saml_logout_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Logout Processor Path",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Hosts Action",
						},
						"groups": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Hosts Group",
						},
						"sso_cookie_update_interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSO Cookie Update Interval",
						},
						"max_failed_attempts": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Failed Attempts Allowed Per IP",
						},
						"count_window": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Count Window",
						},
						"enable_bruteforce_prevention": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Bruteforce Prevention",
						},
						"kerberos_debug_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Debug Logs",
						},
						"kerberos_enable_delegation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kerberos Delegation",
						},
						"kerberos_ldap_authorization": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authorize using Ldap Groups",
						},
						"krb_authorization_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Select Ldap Policy",
						},
						"master_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Master Service",
						},
						"master_service_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Master Service URL",
						},
						"service_provider_display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Organization Display Name",
						},
						"service_provider_entity_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SP Entity ID",
						},
						"service_provider_org_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Organization Name",
						},
						"service_provider_org_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Organization URL",
						},
						"attribute_format": {Type: schema.TypeString, Optional: true, Description: "Format"},
						"attribute_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Attribute",
						},
						"attribute_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Attribute",
						},
						"attribute_type": {Type: schema.TypeString, Optional: true, Description: "Type"},
						"encryption_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Encryption Certificate",
						},
						"signing_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Signing Certificate",
						},
					},
				},
			},
			"caching": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expiry_age": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expiry Age (minutes)",
						},
						"file_extensions": {Type: schema.TypeString, Optional: true, Description: "File Extensions"},
						"max_size":        {Type: schema.TypeString, Optional: true, Description: "Max Size (KB)"},
						"min_size":        {Type: schema.TypeString, Optional: true, Description: "Min Size (B)"},
						"cache_negative_response": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cache Negative Responses",
						},
						"ignore_request_headers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore Request Headers",
						},
						"ignore_response_headers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore Response Headers",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
					},
				},
			},
			"clickjacking": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origin": {Type: schema.TypeString, Optional: true, Description: "Allowed Origin URI"},
						"options": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Render Page Inside Iframe",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
					},
				},
			},
			"compression": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_types": {Type: schema.TypeString, Optional: true, Description: "Content Types"},
						"min_size":      {Type: schema.TypeString, Optional: true, Description: "Min Size (B)"},
						"status":        {Type: schema.TypeString, Optional: true, Description: "Status"},
						"unknown_content_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Compress Unknown Content Types",
						},
					},
				},
			},
			"exception_profiling": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exception_profiling_trusted_host_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Hosts Group",
						},
						"exception_profiling_learn_from_trusted_host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Learn From Trusted Host Group",
						},
						"exception_profiling_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception Profiling Level",
						},
					},
				},
			},
			"ftp_security": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attack_prevention_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "FTP Attack Prevention",
						},
						"allowed_verbs": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "FTP Allowed Verbs",
						},
						"allowed_verb_status": {Type: schema.TypeString, Optional: true, Description: "Status"},
						"pasv_ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PASV IP Address",
						},
						"pasv_ports": {Type: schema.TypeString, Optional: true, Description: "PASV Ports"},
					},
				},
			},
			"instant_ssl": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
						"sharepoint_rewrite_support": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SharePoint Rewrite Support",
						},
						"secure_site_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Secure Site Domain",
						},
					},
				},
			},
			"ip_reputation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"anonymous_proxy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Anonymous Proxy",
						},
						"barracuda_reputation_blocklist": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Barracuda Reputation Blocklist",
						},
						"custom_blacklisted_ip_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom IP List",
						},
						"datacenter_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "DataCenter IP",
						},
						"fake_crawler": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Fake Crawler",
						},
						"check_registered_country": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Check Registered Country",
						},
						"block_unclassified_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unrecognized IP",
						},
						"apply_policy_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Apply Policy at",
						},
						"geo_pool":     {Type: schema.TypeString, Optional: true, Description: "Geo Pool"},
						"geoip_action": {Type: schema.TypeString, Optional: true, Description: "Action"},
						"enable_ip_reputation_filter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable IP Reputation Filter",
						},
						"geoip_enable_logging": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Logging",
						},
						"known_http_attack_sources": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Known HTTP Attack Sources",
						},
						"public_proxy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Public Proxy",
						},
						"satellite_provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Satellite Provider",
						},
						"known_ssh_attack_sources": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Known SSH Attack Sources",
						},
						"tor_nodes": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TOR Nodes",
						},
					},
				},
			},
			"adaptive_profiling": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_types":         {Type: schema.TypeString, Optional: true, Description: "Content Types"},
						"ignore_parameters":     {Type: schema.TypeString, Optional: true, Description: "Ignore Parameters"},
						"navigation_parameters": {Type: schema.TypeString, Optional: true, Description: "Navigation Params"},
						"request_learning":      {Type: schema.TypeString, Optional: true, Description: "Request Learning"},
						"response_learning":     {Type: schema.TypeString, Optional: true, Description: "Response Learning"},
						"status":                {Type: schema.TypeString, Optional: true, Description: "Status"},
						"trusted_host_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Hosts Group",
						},
					},
				},
			},
			"sensitive_parameter_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sensitive_parameter_names": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sensitive Parameter Names",
						},
					},
				},
			},
			"session_tracking": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifiers":         {Type: schema.TypeString, Optional: true, Description: "Session Identifiers"},
						"exception_clients":   {Type: schema.TypeString, Optional: true, Description: "Exception Clients"},
						"max_interval":        {Type: schema.TypeString, Optional: true, Description: "Counting Criterion"},
						"max_sessions_per_ip": {Type: schema.TypeString, Optional: true, Description: "New Session Count"},
						"status":              {Type: schema.TypeString, Optional: true, Description: "Status"},
					},
				},
			},
			"slow_client_attack": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_transfer_rate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data Transfer Rate in KB/Sec",
						},
						"exception_clients": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception Clients",
						},
						"incremental_request_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Incremental Request Timeout",
						},
						"incremental_response_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Incremental Response Timeout",
						},
						"max_request_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Request Timeout",
						},
						"max_response_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Response Timeout",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
					},
				},
			},
			"website_profile": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strict_profile_check": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Strict Profile Check",
						},
						"allowed_domains": {Type: schema.TypeString, Optional: true, Description: "Allowed Domains"},
						"exclude_url_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exclude URL Patterns",
						},
						"include_url_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Include URL Patterns",
						},
						"mode":        {Type: schema.TypeString, Optional: true, Description: "Mode"},
						"use_profile": {Type: schema.TypeString, Optional: true, Description: "Use Profile"},
					},
				},
			},
			"advanced_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_web_application_firewall": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Web Application Firewall",
						},
						"accept_list": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Accept List",
						},
						"accept_list_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Accept List",
						},
						"proxy_list": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy list",
						},
						"proxy_list_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Proxy List",
						},
						"ddos_exception_list": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "DDoS Exception List",
						},
						"enable_fingerprint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Client Fingerprinting",
						},
						"enable_http2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable HTTP2",
						},
						"keepalive_requests": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keepalive Requests",
						},
						"ntlm_ignore_extra_data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "NTLM Ignore Extra Data",
						},
						"enable_proxy_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Proxy Protocol",
						},
						"enable_vdi": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable VDI",
						},
						"enable_websocket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable WebSocket",
						},
					},
				},
			},
			"basic_security": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web_firewall_log_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Web Firewall Log Level",
						},
						"mode": {Type: schema.TypeString, Optional: true, Description: "Mode"},
						"trusted_hosts_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Hosts Action",
						},
						"trusted_hosts_group": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Hosts Group",
						},
						"ignore_case": {Type: schema.TypeString, Optional: true, Description: "Ignore case"},
						"client_ip_addr_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Header for Client IP Address",
						},
						"rate_control_pool": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rate Control Pool",
						},
						"rate_control_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rate Control Status",
						},
						"web_firewall_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Web Firewall Policy",
						},
					},
				},
			},
			"load_balancing": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {Type: schema.TypeString, Optional: true, Description: "Algorithm"},
						"persistence_cookie_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Persistence Cookie Domain",
						},
						"cookie_age": {Type: schema.TypeString, Optional: true, Description: "Cookie Age"},
						"persistence_cookie_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Persistence Cookie Name",
						},
						"persistence_cookie_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Persistence Cookie Path",
						},
						"failover_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Failover Method",
						},
						"header_name":              {Type: schema.TypeString, Optional: true, Description: "Header Name"},
						"persistence_idle_timeout": {Type: schema.TypeString, Optional: true, Description: "Idle Timeout"},
						"persistence_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Persistence Method",
						},
						"source_ip_netmask": {Type: schema.TypeString, Optional: true, Description: "Source IP"},
						"parameter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter Name",
						},
					},
				},
			},
			"ssl_client_authentication": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_certificate_for_rule": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule Group Name",
						},
						"client_authentication_rule_count": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule Group Name",
						},
						"client_authentication": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client Authentication",
						},
						"enforce_client_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enforce Client Certificate",
						},
						"trusted_certificates": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trusted Certificates",
						},
					},
				},
			},
			"ssl_security": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Certificate",
						},
						"ciphers": {Type: schema.TypeString, Optional: true, Description: "Ciphers"},
						"ecdsa_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ECDSA Certificate",
						},
						"include_hsts_sub_domains": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Include HSTS Sub-Domains",
						},
						"hsts_max_age": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HSTS Max-Age",
						},
						"create_hsts_redirect_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specify whether the redirect service should be created when HSTS is enabled.",
						},
						"selected_ciphers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Selected Ciphers",
						},
						"override_ciphers_ssl3": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Override ciphers for SSL 3.0",
						},
						"override_ciphers_tls_1_1": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Override ciphers for TLS 1.1",
						},
						"override_ciphers_tls_1_2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Override ciphers for TLS 1.2",
						},
						"override_ciphers_tls_1_3": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Override ciphers for TLS 1.3",
						},
						"override_ciphers_tls_1": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Override ciphers for TLS 1.0",
						},
						"enable_pfs": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Perfect Forward Secrecy",
						},
						"enable_ssl_3": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL 3.0 (Insecure)",
						},
						"enable_tls_1": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TLS 1.0 (Insecure)",
						},
						"enable_tls_1_1": {Type: schema.TypeString, Optional: true, Description: "TLS 1.1"},
						"enable_tls_1_2": {Type: schema.TypeString, Optional: true, Description: "TLS 1.2"},
						"enable_tls_1_3": {Type: schema.TypeString, Optional: true, Description: "TLS 1.3"},
						"enable_hsts": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable HSTS",
						},
						"enable_ocsp_stapling": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable OCSP Stapling",
						},
						"sni_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain Certificate",
						},
						"domain": {Type: schema.TypeString, Optional: true, Description: "Domain"},
						"sni_ecdsa_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain ECDSA Certificate",
						},
						"enable_sni": {Type: schema.TypeString, Optional: true, Description: "Enable SNI"},
						"enable_strict_sni_check": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Strict SNI Check",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
						"ssl_tls_presets": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL/TLS Quick Settings",
						},
					},
				},
			},
			"captcha_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recaptcha_type":        {Type: schema.TypeString, Optional: true, Description: "Captcha Method"},
						"recaptcha_domain":      {Type: schema.TypeString, Optional: true, Description: "Domain"},
						"recaptcha_site_key":    {Type: schema.TypeString, Optional: true, Description: "Site key"},
						"recaptcha_site_secret": {Type: schema.TypeString, Optional: true, Description: "Site Secret"},
					},
				},
			},
			"ssl_ocsp": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable":        {Type: schema.TypeString, Optional: true, Description: "Enabled"},
						"responder_url": {Type: schema.TypeString, Optional: true, Description: "OCSP Responder URL"},
						"certificate":   {Type: schema.TypeString, Optional: true, Description: "OCSP Issuer Cetificate"},
					},
				},
			},
			"url_encryption": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {Type: schema.TypeString, Optional: true, Description: "URL Encryption"},
					},
				},
			},
			"referer_spam": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exception_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception Patterns",
						},
						"custom_blocked_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom Blocked Attack Types",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
					},
				},
			},
			"comment_spam": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exception_patterns": {Type: schema.TypeString, Optional: true, Description: "Exception Patterns"},
						"parameter":          {Type: schema.TypeString, Optional: true, Description: "Parameter"},
					},
				},
			},
			"advanced_client_analysis": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_analysis": {Type: schema.TypeString, Optional: true, Description: "Advanced Analysis"},
						"exclude_url_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exclude URL Patterns",
						},
					},
				},
			},
			"form_spam": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {Type: schema.TypeString, Optional: true, Description: "Form Spam Status"},
						"honeypot_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Insert Honeypot Field",
						},
						"autoconfigure_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auto Configure Status",
						},
					},
				},
			},
			"waas_account": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"waas_account_id": {Type: schema.TypeString, Required: true, Description: "WaaS Account ID"},
						"waas_account_serial": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "WaaS subscription serial",
						},
					},
				},
			},
		},

		Description: "`barracudawaf_services` manages `Services` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFServicesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFServicesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFServicesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFServicesRead(d, m)
}

func resourceCudaWAFServicesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	var dataItems map[string]interface{}
	resources, err := client.GetBarracudaWAFResource(name, request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	for _, dataItems = range resources.Data {
		if dataItems["name"] == name {
			break
		}
	}

	if dataItems["name"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)
	return nil
}

func resourceCudaWAFServicesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFServicesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFServicesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFServicesRead(d, m)
}

func resourceCudaWAFServicesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services"
	request := &APIRequest{
		Method: "delete",
		URL:    resourceEndpoint,
	}

	err := client.DeleteBarracudaWAFResource(name, request)

	if err != nil {
		return fmt.Errorf("Unable to delete the Barracuda WAF resource (%s) (%v)", name, err)
	}

	return nil
}

func hydrateBarracudaWAFServicesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"address-version":     d.Get("address_version").(string),
		"dps-enabled":         d.Get("dps_enabled").(string),
		"mask":                d.Get("mask").(string),
		"session-timeout":     d.Get("session_timeout").(string),
		"linked-service-name": d.Get("linked_service_name").(string),
		"enable-access-logs":  d.Get("enable_access_logs").(string),
		"app-id":              d.Get("app_id").(string),
		"comments":            d.Get("comments").(string),
		"group":               d.Get("group").(string),
		"service-id":          d.Get("service_id").(string),
		"ip-address":          d.Get("ip_address").(string),
		"cloud-ip-select":     d.Get("cloud_ip_select").(string),
		"name":                d.Get("name").(string),
		"port":                d.Get("port").(string),
		"status":              d.Get("status").(string),
		"type":                d.Get("type").(string),
		"certificate":         d.Get("certificate").(string),
		"service-hostname":    d.Get("service_hostname").(string),
		"vsite":               d.Get("vsite").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"address-version", "group", "vsite"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if len(val) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}

func (b *BarracudaWAF) hydrateBarracudaWAFServicesSubResource(d *schema.ResourceData, name string, endpoint string) error {

	for subResource, subResourceParams := range subResourceServicesParams {
		subResourceParamsLength := d.Get(subResource + ".#").(int)

		log.Printf("[INFO] Updating Barracuda WAF sub resource (%s) (%s)", name, subResource)

		for i := 0; i < subResourceParamsLength; i++ {
			subResourcePayload := map[string]string{}
			suffix := fmt.Sprintf(".%d", i)

			for _, param := range subResourceParams {
				paramSuffix := fmt.Sprintf(".%s", param)
				paramVaule := d.Get(subResource + suffix + paramSuffix).(string)

				if len(paramVaule) > 0 {
					param = strings.Replace(param, "_", "-", -1)
					subResourcePayload[param] = paramVaule
				}
			}

			err := b.UpdateBarracudaWAFSubResource(name, endpoint, &APIRequest{
				URL:  strings.Replace(subResource, "_", "-", -1),
				Body: subResourcePayload,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
