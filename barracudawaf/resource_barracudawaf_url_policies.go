package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceUrlPoliciesParams = map[string][]string{}
)

func resourceCudaWAFUrlPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFUrlPoliciesCreate,
		Read:   resourceCudaWAFUrlPoliciesRead,
		Update: resourceCudaWAFUrlPoliciesUpdate,
		Delete: resourceCudaWAFUrlPoliciesDelete,

		Schema: map[string]*schema.Schema{
			"enable_data_theft_protection": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable Data Theft Protection",
			},
			"enable_batd_scan": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable BATP Scan",
			},
			"comments": {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"host":     {Type: schema.TypeString, Optional: true, Description: "Host Match"},
			"extended_match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended Match",
			},
			"extended_match_sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended Match Sequence",
			},
			"mode": {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL Policy Name",
			},
			"parse_urls_in_scripts": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parse URLs in Scripts",
			},
			"rate_control_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rate Control Pool",
			},
			"response_charset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Response Charset",
			},
			"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
			"url":    {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"enable_virus_scan": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable Virus Scan",
			},
			"web_scraping_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Web Scraping Policy",
			},
			"counting_criterion": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Counting Criterion",
			},
			"enable_count_auth_resp_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Count Auth Response Codes",
			},
			"exception_clients": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exception Clients",
			},
			"exception_fingerprints": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exception Fingerprints",
			},
			"max_allowed_accesses_per_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Allowed Accesses Per IP",
			},
			"max_allowed_accesses_per_fingerprint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Allowed Accesses Per Client Fingerprint",
			},
			"max_allowed_accesses_from_all_sources": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Allowed Accesses From All Sources",
			},
			"max_bandwidth_per_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Max Bandwidth Per IP",
			},
			"max_bandwidth_per_fingerprint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Max Bandwidth Per Client Fingerprint",
			},
			"max_bandwidth_from_all_sources": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Max Bandwidth From All Sources",
			},
			"max_failed_accesses_per_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Failed Accesses Per IP",
			},
			"max_failed_accesses_per_fingerprint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Failed Accesses Per Client Fingerprint",
			},
			"max_failed_accesses_from_all_sources": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Failed Accesses From All Sources",
			},
			"count_window": {Type: schema.TypeString, Optional: true, Description: "Count Window"},
			"enable_bruteforce_prevention": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable Bruteforce Prevention",
			},
			"credential_stuffing_username_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username Parameter",
			},
			"credential_stuffing_password_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Password Parameter",
			},
			"credential_spraying_blocking_threshold": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block Threshold",
			},
			"credential_protection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protection Type",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_url_policies` manages `Url Policies` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFUrlPoliciesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-policies"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFUrlPoliciesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFUrlPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFUrlPoliciesRead(d, m)
}

func resourceCudaWAFUrlPoliciesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-policies"
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

func resourceCudaWAFUrlPoliciesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-policies"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFUrlPoliciesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFUrlPoliciesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFUrlPoliciesRead(d, m)
}

func resourceCudaWAFUrlPoliciesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-policies"
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

func hydrateBarracudaWAFUrlPoliciesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"enable-data-theft-protection":           d.Get("enable_data_theft_protection").(string),
		"enable-batd-scan":                       d.Get("enable_batd_scan").(string),
		"comments":                               d.Get("comments").(string),
		"host":                                   d.Get("host").(string),
		"extended-match":                         d.Get("extended_match").(string),
		"extended-match-sequence":                d.Get("extended_match_sequence").(string),
		"mode":                                   d.Get("mode").(string),
		"name":                                   d.Get("name").(string),
		"parse-urls-in-scripts":                  d.Get("parse_urls_in_scripts").(string),
		"rate-control-pool":                      d.Get("rate_control_pool").(string),
		"response-charset":                       d.Get("response_charset").(string),
		"status":                                 d.Get("status").(string),
		"url":                                    d.Get("url").(string),
		"enable-virus-scan":                      d.Get("enable_virus_scan").(string),
		"web-scraping-policy":                    d.Get("web_scraping_policy").(string),
		"counting-criterion":                     d.Get("counting_criterion").(string),
		"enable-count-auth-resp-code":            d.Get("enable_count_auth_resp_code").(string),
		"exception-clients":                      d.Get("exception_clients").(string),
		"exception-fingerprints":                 d.Get("exception_fingerprints").(string),
		"max-allowed-accesses-per-ip":            d.Get("max_allowed_accesses_per_ip").(string),
		"max-allowed-accesses-per-fingerprint":   d.Get("max_allowed_accesses_per_fingerprint").(string),
		"max-allowed-accesses-from-all-sources":  d.Get("max_allowed_accesses_from_all_sources").(string),
		"max-bandwidth-per-ip":                   d.Get("max_bandwidth_per_ip").(string),
		"max-bandwidth-per-fingerprint":          d.Get("max_bandwidth_per_fingerprint").(string),
		"max-bandwidth-from-all-sources":         d.Get("max_bandwidth_from_all_sources").(string),
		"max-failed-accesses-per-ip":             d.Get("max_failed_accesses_per_ip").(string),
		"max-failed-accesses-per-fingerprint":    d.Get("max_failed_accesses_per_fingerprint").(string),
		"max-failed-accesses-from-all-sources":   d.Get("max_failed_accesses_from_all_sources").(string),
		"count-window":                           d.Get("count_window").(string),
		"enable-bruteforce-prevention":           d.Get("enable_bruteforce_prevention").(string),
		"credential-stuffing-username-field":     d.Get("credential_stuffing_username_field").(string),
		"credential-stuffing-password-field":     d.Get("credential_stuffing_password_field").(string),
		"credential-spraying-blocking-threshold": d.Get("credential_spraying_blocking_threshold").(string),
		"credential-protection-type":             d.Get("credential_protection_type").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{}
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

func (b *BarracudaWAF) hydrateBarracudaWAFUrlPoliciesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceUrlPoliciesParams {
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
