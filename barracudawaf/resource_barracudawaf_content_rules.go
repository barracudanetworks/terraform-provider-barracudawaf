package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceContentRulesParams = map[string][]string{
		"caching": {
			"expiry_age",
			"file_extensions",
			"max_size",
			"min_size",
			"cache_negative_responses",
			"ignore_request_headers",
			"ignore_response_headers",
			"status",
		},
		"compression": {"content_types", "min_size", "status", "compress_unknown_content_types"},
		"load_balancing": {
			"lb_algorithm",
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
		"captcha_settings": {
			"recaptcha_type",
			"rg_recaptcha_domain",
			"rg_recaptcha_site_key",
			"rg_recaptcha_site_secret",
		},
		"advanced_client_analysis": {"advanced_analysis"},
	}
)

func resourceCudaWAFContentRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFContentRulesCreate,
		Read:   resourceCudaWAFContentRulesRead,
		Update: resourceCudaWAFContentRulesUpdate,
		Delete: resourceCudaWAFContentRulesDelete,

		Schema: map[string]*schema.Schema{
			"access_log":              {Type: schema.TypeString, Optional: true, Description: "Access Log"},
			"app_id":                  {Type: schema.TypeString, Optional: true, Description: "Rule App Id"},
			"comments":                {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"host_match":              {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"name":                    {Type: schema.TypeString, Required: true, Description: "Rule Group Name"},
			"status":                  {Type: schema.TypeString, Optional: true, Description: "Status"},
			"extended_match":          {Type: schema.TypeString, Optional: true, Description: "Extended Match"},
			"extended_match_sequence": {Type: schema.TypeString, Optional: true, Description: "Extended Match Sequence"},
			"mode":                    {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"url_match":               {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"web_firewall_policy":     {Type: schema.TypeString, Optional: true, Description: "Web Firewall Policy"},
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
						"file_extensions": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File Extensions",
						},
						"max_size": {Type: schema.TypeString, Optional: true, Description: "Max Size (KB)"},
						"min_size": {Type: schema.TypeString, Optional: true, Description: "Min Size (B)"},
						"cache_negative_responses": {
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
			"compression": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content Types",
						},
						"min_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Min Size (B)",
						},
						"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
						"compress_unknown_content_types": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Compress Others",
						},
					},
				},
			},
			"load_balancing": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lb_algorithm": {Type: schema.TypeString, Optional: true, Description: "Algorithm"},
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
			"captcha_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recaptcha_type":           {Type: schema.TypeString, Optional: true, Description: "Captcha Method"},
						"rg_recaptcha_domain":      {Type: schema.TypeString, Optional: true, Description: "Domain"},
						"rg_recaptcha_site_key":    {Type: schema.TypeString, Optional: true, Description: "Site key"},
						"rg_recaptcha_site_secret": {Type: schema.TypeString, Optional: true, Description: "Site Secret"},
					},
				},
			},
			"advanced_client_analysis": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_analysis": {Type: schema.TypeString, Optional: true, Description: "Advanced Analysis"},
					},
				},
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_content_rules` manages `Content Rules` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFContentRulesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFContentRulesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFContentRulesRead(d, m)
}

func resourceCudaWAFContentRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules"
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

func resourceCudaWAFContentRulesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFContentRulesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFContentRulesRead(d, m)
}

func resourceCudaWAFContentRulesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules"
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

func hydrateBarracudaWAFContentRulesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"access-log":              d.Get("access_log").(string),
		"app-id":                  d.Get("app_id").(string),
		"comments":                d.Get("comments").(string),
		"host-match":              d.Get("host_match").(string),
		"name":                    d.Get("name").(string),
		"status":                  d.Get("status").(string),
		"extended-match":          d.Get("extended_match").(string),
		"extended-match-sequence": d.Get("extended_match_sequence").(string),
		"mode":                    d.Get("mode").(string),
		"url-match":               d.Get("url_match").(string),
		"web-firewall-policy":     d.Get("web_firewall_policy").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFContentRulesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceContentRulesParams {
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
