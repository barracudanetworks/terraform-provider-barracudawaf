package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceCredentialServersParams = map[string][]string{}
)

func resourceCudaWAFCredentialServers() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFCredentialServersCreate,
		Read:   resourceCudaWAFCredentialServersRead,
		Update: resourceCudaWAFCredentialServersUpdate,
		Delete: resourceCudaWAFCredentialServersDelete,

		Schema: map[string]*schema.Schema{
			"cache_expiry":         {Type: schema.TypeString, Optional: true, Description: "Cache Expiry (Seconds)"},
			"cache_valid_sessions": {Type: schema.TypeString, Optional: true, Description: "Cache Valid Sessions"},
			"redirect_url":         {Type: schema.TypeString, Optional: true, Description: "Redirect URL"},
			"ip_address":           {Type: schema.TypeString, Required: true, Description: "Server Name/IP Address"},
			"policy_name":          {Type: schema.TypeString, Required: true, Description: "Policy Name"},
			"port":                 {Type: schema.TypeString, Optional: true, Description: "Server Port"},
			"armored_browser_type": {Type: schema.TypeString, Required: true, Description: "Armored Browser Type"},
			"name":                 {Type: schema.TypeString, Required: true, Description: "Name"},
		},

		Description: "`barracudawaf_credential_servers` manages `Credential Servers` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFCredentialServersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/credential-servers"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFCredentialServersResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFCredentialServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFCredentialServersRead(d, m)
}

func resourceCudaWAFCredentialServersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/credential-servers"
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

func resourceCudaWAFCredentialServersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/credential-servers"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFCredentialServersResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFCredentialServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFCredentialServersRead(d, m)
}

func resourceCudaWAFCredentialServersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/credential-servers"
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

func hydrateBarracudaWAFCredentialServersResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"cache-expiry":         d.Get("cache_expiry").(string),
		"cache-valid-sessions": d.Get("cache_valid_sessions").(string),
		"redirect-url":         d.Get("redirect_url").(string),
		"ip-address":           d.Get("ip_address").(string),
		"policy-name":          d.Get("policy_name").(string),
		"port":                 d.Get("port").(string),
		"armored-browser-type": d.Get("armored_browser_type").(string),
		"name":                 d.Get("name").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFCredentialServersSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceCredentialServersParams {
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
