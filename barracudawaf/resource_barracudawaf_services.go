package barracudawaf

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceServicesParams = map[string][]string{
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
		"ssl_security": {
			"certificate",
			"ciphers",
			"ecdsa_certificate",
			"include_hsts_sub_domains",
			"hsts_max_age",
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
	}
)

func resourceCudaWAFServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFServicesCreate,
		Read:   resourceCudaWAFServicesRead,
		Update: resourceCudaWAFServicesUpdate,
		Delete: resourceCudaWAFServicesDelete,

		Schema: map[string]*schema.Schema{
			"address_version":    {Type: schema.TypeString, Optional: true, Description: "Version"},
			"mask":               {Type: schema.TypeString, Optional: true, Description: "Mask"},
			"session_timeout":    {Type: schema.TypeString, Optional: true, Description: "Session Timeout"},
			"enable_access_logs": {Type: schema.TypeString, Optional: true, Description: "Enable Access Logs"},
			"app_id":             {Type: schema.TypeString, Optional: true, Description: "Service App Id"},
			"comments":           {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"group":              {Type: schema.TypeString, Optional: true, Description: "Service Group"},
			"ip_address":         {Type: schema.TypeString, Optional: true, Description: "VIP"},
			"cloud_ip_select":    {Type: schema.TypeString, Optional: true},
			"name":               {Type: schema.TypeString, Required: true, Description: "Web Application Name"},
			"port":               {Type: schema.TypeString, Optional: true, Description: "Port"},
			"status":             {Type: schema.TypeString, Optional: true, Description: "Status"},
			"type":               {Type: schema.TypeString, Optional: true, Description: "Type"},
			"certificate":        {Type: schema.TypeString, Optional: true},
			"service_hostname":   {Type: schema.TypeString, Optional: true},
			"vsite":              {Type: schema.TypeString, Optional: true, Description: "Vsite"},
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
						"selected_ciphers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Selected Ciphers",
						},
						"override_ciphers_ssl3": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Override ciphers for SSL 3.0",
						},
						"override_ciphers_tls_1_1": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Override ciphers for TLS 1.1",
						},
						"override_ciphers_tls_1_2": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Override ciphers for TLS 1.2",
						},
						"override_ciphers_tls_1_3": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Override ciphers for TLS 1.3",
						},
						"override_ciphers_tls_1": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Domain Certificate",
						},
						"domain": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Domain",
						},
						"sni_ecdsa_certificate": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
			"secure_site_domain": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Secure Site Domain",
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
	resourcePayload := map[string]interface{}{
		"address-version":    d.Get("address_version").(string),
		"mask":               d.Get("mask").(string),
		"session-timeout":    d.Get("session_timeout").(string),
		"enable-access-logs": d.Get("enable_access_logs").(string),
		"app-id":             d.Get("app_id").(string),
		"comments":           d.Get("comments").(string),
		"group":              d.Get("group").(string),
		"ip-address":         d.Get("ip_address").(string),
		"cloud-ip-select":    d.Get("cloud_ip_select").(string),
		"name":               d.Get("name").(string),
		"port":               d.Get("port").(string),
		"status":             d.Get("status").(string),
		"type":               d.Get("type").(string),
		"certificate":        d.Get("certificate").(string),
		"service-hostname":   d.Get("service_hostname").(string),
		"vsite":              d.Get("vsite").(string),
		"secure-site-domain": d.Get("secure_site_domain"),
	}

	log.Println("[DEBUG] Resource payload for the service REST Call : ", resourcePayload)

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"address-version", "group", "vsite", "certificate", "secure-site-domain"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if reflect.ValueOf(val).Len() == 0 {
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
			subResourcePayload := make(map[string]interface{})
			suffix := fmt.Sprintf(".%d", i)

			for _, param := range subResourceParams {
				paramSuffix := fmt.Sprintf(".%s", param)
				paramVaule := d.Get(subResource + suffix + paramSuffix)

				if reflect.ValueOf(paramVaule).Kind() == reflect.String {
					paramVaule = paramVaule.(string)
				}

				if reflect.ValueOf(paramVaule).Len() > 0 {
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
