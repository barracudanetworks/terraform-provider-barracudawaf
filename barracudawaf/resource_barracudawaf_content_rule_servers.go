package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceContentRuleServersParams = map[string][]string{
		"advanced_configuration": {
			"client_impersonation",
			"source_ip_to_connect",
			"max_connections",
			"max_establishing_connections",
			"max_requests",
			"max_keepalive_requests",
			"max_spare_connections",
			"timeout",
		},
		"out_of_band_health_checks": {"interval", "enable_oob_health_checks"},
		"application_layer_health_checks": {
			"additional_headers",
			"match_content_string",
			"method",
			"status_code",
			"url",
			"domain",
		},
		"connection_pooling":    {"keepalive_timeout", "enable_connection_pooling"},
		"redirect":              {},
		"in_band_health_checks": {"max_other_failure", "max_refused", "max_timeout_failure", "max_http_errors"},
		"ssl_policy": {
			"client_certificate",
			"enable_ssl_compatibility_mode",
			"enable_ssl_3",
			"enable_tls_1",
			"enable_tls_1_1",
			"enable_tls_1_2",
			"enable_tls_1_3",
			"validate_certificate",
			"enable_https",
			"enable_sni",
		},
		"load_balancing": {"backup_server", "weight"},
	}
)

func resourceCudaWAFContentRuleServers() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFContentRuleServersCreate,
		Read:   resourceCudaWAFContentRuleServersRead,
		Update: resourceCudaWAFContentRuleServersUpdate,
		Delete: resourceCudaWAFContentRuleServersDelete,

		Schema: map[string]*schema.Schema{
			"comments":        {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"name":            {Type: schema.TypeString, Optional: true, Description: "Web Server Name"},
			"hostname":        {Type: schema.TypeString, Optional: true, Description: "Hostname"},
			"identifier":      {Type: schema.TypeString, Optional: true, Description: "Identifier:"},
			"ip_address":      {Type: schema.TypeString, Optional: true, Description: "IP Address"},
			"address_version": {Type: schema.TypeString, Optional: true, Description: "Version"},
			"port":            {Type: schema.TypeString, Optional: true, Description: "Port"},
			"status":          {Type: schema.TypeString, Optional: true, Description: "Status"},
			"resolved_ips":    {Type: schema.TypeString, Optional: true},
			"advanced_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_impersonation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client Impersonation",
						},
						"source_ip_to_connect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source IP to Connect",
						},
						"max_connections": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Connections",
						},
						"max_establishing_connections": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Establishing Connections",
						},
						"max_requests": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Requests",
						},
						"max_keepalive_requests": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Keepalive Requests",
						},
						"max_spare_connections": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Spare Connections",
						},
						"timeout": {Type: schema.TypeString, Optional: true, Description: "Timeout"},
					},
				},
			},
			"out_of_band_health_checks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval":                 {Type: schema.TypeString, Optional: true, Description: "Interval"},
						"enable_oob_health_checks": {Type: schema.TypeString, Optional: true, Description: "Status"},
					},
				},
			},
			"application_layer_health_checks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_headers": {Type: schema.TypeString, Optional: true, Description: "Additional Headers"},
						"match_content_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Match content String",
						},
						"method":      {Type: schema.TypeString, Optional: true, Description: "Method"},
						"status_code": {Type: schema.TypeString, Optional: true, Description: "Status Code"},
						"url":         {Type: schema.TypeString, Optional: true, Description: "URL"},
						"domain":      {Type: schema.TypeString, Optional: true, Description: "Domain"},
					},
				},
			},
			"connection_pooling": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keepalive_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keepalive Timeout",
						},
						"enable_connection_pooling": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Connection Pooling",
						},
					},
				},
			},
			"redirect": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Resource{Schema: map[string]*schema.Schema{}},
			},
			"in_band_health_checks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_other_failure": {Type: schema.TypeString, Optional: true, Description: "Max Other Failure"},
						"max_refused":       {Type: schema.TypeString, Optional: true, Description: "Max Refused"},
						"max_timeout_failure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Max Timeout Failures",
						},
						"max_http_errors": {Type: schema.TypeString, Optional: true, Description: "Max HTTP Errors"},
					},
				},
			},
			"ssl_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client Certificate",
						},
						"enable_ssl_compatibility_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable SSL Compatibility Mode",
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
						"validate_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Validate Server Certificate",
						},
						"enable_https": {Type: schema.TypeString, Optional: true, Description: "Status"},
						"enable_sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable SNI",
						},
					},
				},
			},
			"load_balancing": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_server": {Type: schema.TypeString, Optional: true, Description: "Backup Appliance"},
						"weight":        {Type: schema.TypeString, Optional: true, Description: "WRR Weight"},
					},
				},
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_content_rule_servers` manages `Content Rule Servers` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFContentRuleServersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules/" + d.Get("parent.1").(string) + "/content-rule-servers"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFContentRuleServersResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRuleServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFContentRuleServersRead(d, m)
}

func resourceCudaWAFContentRuleServersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules/" + d.Get("parent.1").(string) + "/content-rule-servers"
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

func resourceCudaWAFContentRuleServersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules/" + d.Get("parent.1").(string) + "/content-rule-servers"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFContentRuleServersResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRuleServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFContentRuleServersRead(d, m)
}

func resourceCudaWAFContentRuleServersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/content-rules/" + d.Get("parent.1").(string) + "/content-rule-servers"
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

func hydrateBarracudaWAFContentRuleServersResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"comments":        d.Get("comments").(string),
		"name":            d.Get("name").(string),
		"hostname":        d.Get("hostname").(string),
		"identifier":      d.Get("identifier").(string),
		"ip-address":      d.Get("ip_address").(string),
		"address-version": d.Get("address_version").(string),
		"port":            d.Get("port").(string),
		"status":          d.Get("status").(string),
		"resolved-ips":    d.Get("resolved_ips").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"address-version", "resolved-ips"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFContentRuleServersSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceContentRuleServersParams {
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
