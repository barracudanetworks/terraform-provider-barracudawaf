package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceServersParams = map[string][]string{
		"ssl_policy": {
			"client_certificate",
			"enable_ssl_compatibility_mode",
			"validate_certificate",
			"enable_https",
			"enable_sni",
			"enable_ssl_3",
			"enable_tls_1",
			"enable_tls_1_1",
			"enable_tls_1_2",
			"enable_tls_1_3",
		},
		"connection_pooling": {
			"keepalive_timeout",
			"enable_connection_pooling",
		},
	}
)

func resourceCudaWAFServers() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFServersCreate,
		Read:   resourceCudaWAFServersRead,
		Update: resourceCudaWAFServersUpdate,
		Delete: resourceCudaWAFServersDelete,

		Schema: map[string]*schema.Schema{
			"address_version": {Type: schema.TypeString, Optional: true, Description: "Version"},
			"comments":        {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"name":            {Type: schema.TypeString, Optional: true, Description: "Server Name"},
			"hostname":        {Type: schema.TypeString, Optional: true, Description: "Hostname"},
			"identifier":      {Type: schema.TypeString, Optional: true, Description: "Identifier"},
			"ip_address":      {Type: schema.TypeString, Optional: true, Description: "Server IP"},
			"port":            {Type: schema.TypeString, Optional: true, Description: "Server Port"},
			"status":          {Type: schema.TypeString, Optional: true, Description: "Status"},
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
						"validate_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Validate Server Certificate",
						},
						"enable_https": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Server uses SSL",
						},
						"enable_sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable SNI",
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
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_servers` manages `Servers` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFServersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/servers"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFServersResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFServersRead(d, m)
}

func resourceCudaWAFServersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/servers"
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

func resourceCudaWAFServersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/servers"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFServersResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFServersRead(d, m)
}

func resourceCudaWAFServersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/servers"
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

func hydrateBarracudaWAFServersResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"address-version": d.Get("address_version").(string),
		"comments":        d.Get("comments").(string),
		"name":            d.Get("name").(string),
		"hostname":        d.Get("hostname").(string),
		"identifier":      d.Get("identifier").(string),
		"ip-address":      d.Get("ip_address").(string),
		"port":            d.Get("port").(string),
		"status":          d.Get("status").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"address-version"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFServersSubResource(d *schema.ResourceData, name string, endpoint string) error {

	for subResource, subResourceParams := range subResourceServersParams {
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
