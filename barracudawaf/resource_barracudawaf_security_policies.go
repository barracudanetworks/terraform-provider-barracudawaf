package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceSecurityPoliciesParams = map[string][]string{
		"request_limits": {
			"max_cookie_name_length",
			"max_number_of_cookies",
			"max_header_name_length",
			"max_request_line_length",
			"max_number_of_headers",
			"max_query_length",
			"max_cookie_value_length",
			"max_header_value_length",
			"max_request_length",
			"max_url_length",
			"enable",
		},
		"url_normalization": {
			"default_charset",
			"detect_response_charset",
			"normalize_special_chars",
			"apply_double_decoding",
			"parameter_separators",
		},
		"url_protection": {
			"allowed_methods",
			"allowed_content_types",
			"custom_blocked_attack_types",
			"exception_patterns",
			"blocked_attack_types",
			"max_content_length",
			"maximum_parameter_name_length",
			"max_parameters",
			"maximum_upload_files",
			"csrf_prevention",
			"enable",
		},
		"parameter_protection": {
			"blocked_attack_types",
			"custom_blocked_attack_types",
			"base64_decode_parameter_value",
			"allowed_file_upload_type",
			"denied_metacharacters",
			"exception_patterns",
			"file_upload_extensions",
			"file_upload_mime_types",
			"maximum_instances",
			"maximum_parameter_value_length",
			"maximum_upload_file_size",
			"enable",
			"validate_parameter_name",
			"ignore_parameters",
		},
		"cloaking": {
			"return_codes_to_exempt",
			"headers_to_filter",
			"filter_response_header",
			"suppress_return_code",
		},
		"cookie_security": {
			"allow_unrecognized_cookies",
			"days_allowed",
			"cookies_exempted",
			"http_only",
			"cookie_max_age",
			"tamper_proof_mode",
			"secure_cookie",
			"cookie_replay_protection_type",
			"custom_headers",
		},
		"client_profile": {"medium_risk_score", "high_risk_score", "exception_client_fingerprints", "client_profile"},
		"tarpit_profile": {"backlog_requests_limit", "tarpit_inactivity_timeout"},
	}
)

func resourceCudaWAFSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFSecurityPoliciesCreate,
		Read:   resourceCudaWAFSecurityPoliciesRead,
		Update: resourceCudaWAFSecurityPoliciesUpdate,
		Delete: resourceCudaWAFSecurityPoliciesDelete,

		Schema: map[string]*schema.Schema{
			"based_on": {Type: schema.TypeString, Optional: true},
			"name":     {Type: schema.TypeString, Required: true, Description: "Policy Name"},
			"request_limits": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_cookie_name_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Cookie Name Length",
						},
						"max_number_of_cookies": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Number of Cookies",
						},
						"max_header_name_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Header Name Length",
						},
						"max_request_line_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Request Line Length",
						},
						"max_number_of_headers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Number of Headers",
						},
						"max_query_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Query Length",
						},
						"max_cookie_value_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Cookie Value Length",
						},
						"max_header_value_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Header Value Length",
						},
						"max_request_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Request Length",
						},
						"max_url_length": {Type: schema.TypeString, Optional: true, Description: "Max URL Length"},
						"enable": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Request Limits",
						},
					},
				},
			},
			"url_normalization": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_charset": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default Character Set",
						},
						"detect_response_charset": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detect Response Charset",
						},
						"normalize_special_chars": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detect Response Charset",
						},
						"apply_double_decoding": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Apply Double Decoding",
						},
						"parameter_separators": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter Separators",
						},
					},
				},
			},
			"url_protection": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_methods": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Allowed Methods",
						},
						"allowed_content_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Allowed Content Types",
						},
						"custom_blocked_attack_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom Blocked Attack Types",
						},
						"exception_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception Patterns",
						},
						"blocked_attack_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Blocked Attack Types",
						},
						"max_content_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Content Length",
						},
						"maximum_parameter_name_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum Parameter Name Length",
						},
						"max_parameters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Parameters",
						},
						"maximum_upload_files": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum Upload Files",
						},
						"csrf_prevention": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CSRF Prevention",
						},
						"enable": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable URL Protection",
						},
					},
				},
			},
			"parameter_protection": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blocked_attack_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Blocked Attack Types",
						},
						"custom_blocked_attack_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom Blocked Attack Types",
						},
						"base64_decode_parameter_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64 Decode Parameter Value",
						},
						"allowed_file_upload_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Allowed File Upload Type",
						},
						"denied_metacharacters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Denied Metacharacters",
						},
						"exception_patterns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exception Patterns",
						},
						"file_upload_extensions": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File Upload Extensions",
						},
						"file_upload_mime_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File Upload Mime Types",
						},
						"maximum_instances": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum Instances",
						},
						"maximum_parameter_value_length": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum Parameter Value Length",
						},
						"maximum_upload_file_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum Upload File Size",
						},
						"enable": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Parameter Protection",
						},
						"validate_parameter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Validate Parameter Name",
						},
						"ignore_parameters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ignore Parameters",
						},
					},
				},
			},
			"cloaking": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"return_codes_to_exempt": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Return Codes to Exempt",
						},
						"headers_to_filter": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Headers to Filter",
						},
						"filter_response_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter Response Header",
						},
						"suppress_return_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Suppress Return Code",
						},
					},
				},
			},
			"cookie_security": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_unrecognized_cookies": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Allow Unrecognized Cookies",
						},
						"days_allowed": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Days Allowed",
						},
						"cookies_exempted": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cookies Exempted",
						},
						"http_only": {Type: schema.TypeString, Optional: true, Description: "HTTP Only"},
						"cookie_max_age": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cookie Max Age",
						},
						"tamper_proof_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tamper Proof Mode",
						},
						"secure_cookie": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Secure Cookie",
						},
						"cookie_replay_protection_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cookie Replay Protection Type",
						},
						"custom_headers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom Headers",
						},
					},
				},
			},
			"client_profile": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"medium_risk_score": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Suspicious Clients",
						},
						"high_risk_score": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bad Clients",
						},
						"exception_client_fingerprints": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Exceptions",
						},
						"client_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Client Profile Validation",
						},
					},
				},
			},
			"tarpit_profile": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backlog_requests_limit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Backlog Requests Limit",
						},
						"tarpit_inactivity_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tarpit Inactivity Timeout",
						},
					},
				},
			},
		},

		Description: "`barracudawaf_security_policies` manages `Security Policies` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFSecurityPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFSecurityPoliciesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSecurityPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFSecurityPoliciesRead(d, m)
}

func resourceCudaWAFSecurityPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
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

func resourceCudaWAFSecurityPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFSecurityPoliciesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFSecurityPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFSecurityPoliciesRead(d, m)
}

func resourceCudaWAFSecurityPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/security-policies"
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

func hydrateBarracudaWAFSecurityPoliciesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{"based-on": d.Get("based_on").(string), "name": d.Get("name").(string)}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"based-on"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFSecurityPoliciesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceSecurityPoliciesParams {
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
