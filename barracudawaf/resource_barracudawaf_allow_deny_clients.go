package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceAllowDenyClientsParams = map[string][]string{}
)

func resourceCudaWAFAllowDenyClients() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAllowDenyClientsCreate,
		Read:   resourceCudaWAFAllowDenyClientsRead,
		Update: resourceCudaWAFAllowDenyClientsUpdate,
		Delete: resourceCudaWAFAllowDenyClientsDelete,

		Schema: map[string]*schema.Schema{
			"name":                {Type: schema.TypeString, Required: true, Description: "Rule Name"},
			"action":              {Type: schema.TypeString, Optional: true, Description: "Action"},
			"sequence":            {Type: schema.TypeString, Optional: true, Description: "Sequence"},
			"certificate_serial":  {Type: schema.TypeString, Optional: true, Description: "Certificate Serial Number"},
			"common_name":         {Type: schema.TypeString, Optional: true, Description: "Common Name"},
			"country":             {Type: schema.TypeString, Optional: true, Description: "Country"},
			"locality":            {Type: schema.TypeString, Optional: true, Description: "Locality"},
			"organization":        {Type: schema.TypeString, Optional: true, Description: "Organization"},
			"organizational_unit": {Type: schema.TypeString, Optional: true, Description: "Organizational Unit"},
			"state":               {Type: schema.TypeString, Optional: true, Description: "State"},
			"status":              {Type: schema.TypeString, Optional: true, Description: "Status"},
			"parent":              {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_allow_deny_clients` manages `Allow Deny Clients` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAllowDenyClientsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/allow-deny-clients"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFAllowDenyClientsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAllowDenyClientsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAllowDenyClientsRead(d, m)
}

func resourceCudaWAFAllowDenyClientsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/allow-deny-clients"
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

func resourceCudaWAFAllowDenyClientsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/allow-deny-clients"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFAllowDenyClientsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAllowDenyClientsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAllowDenyClientsRead(d, m)
}

func resourceCudaWAFAllowDenyClientsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/allow-deny-clients"
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

func hydrateBarracudaWAFAllowDenyClientsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                d.Get("name").(string),
		"action":              d.Get("action").(string),
		"sequence":            d.Get("sequence").(string),
		"certificate-serial":  d.Get("certificate_serial").(string),
		"common-name":         d.Get("common_name").(string),
		"country":             d.Get("country").(string),
		"locality":            d.Get("locality").(string),
		"organization":        d.Get("organization").(string),
		"organizational-unit": d.Get("organizational_unit").(string),
		"state":               d.Get("state").(string),
		"status":              d.Get("status").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFAllowDenyClientsSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceAllowDenyClientsParams {
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
