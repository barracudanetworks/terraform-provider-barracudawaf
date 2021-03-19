package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceTrustedHostsParams = map[string][]string{}
)

func resourceCudaWAFTrustedHosts() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFTrustedHostsCreate,
		Read:   resourceCudaWAFTrustedHostsRead,
		Update: resourceCudaWAFTrustedHostsUpdate,
		Delete: resourceCudaWAFTrustedHostsDelete,

		Schema: map[string]*schema.Schema{
			"ip_address":   {Type: schema.TypeString, Optional: true, Description: "IP Address"},
			"ipv6_address": {Type: schema.TypeString, Optional: true, Description: "IPv6 Address"},
			"ipv6_mask":    {Type: schema.TypeString, Optional: true, Description: "Mask"},
			"mask":         {Type: schema.TypeString, Optional: true, Description: "Mask"},
			"comments":     {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"name":         {Type: schema.TypeString, Required: true, Description: "Trusted Host Name"},
			"version":      {Type: schema.TypeString, Optional: true, Description: "Version"},
			"parent":       {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_trusted_hosts` manages `Trusted Hosts` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFTrustedHostsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-host-groups/" + d.Get("parent.0").(string) + "/trusted-hosts"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFTrustedHostsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFTrustedHostsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFTrustedHostsRead(d, m)
}

func resourceCudaWAFTrustedHostsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-host-groups/" + d.Get("parent.0").(string) + "/trusted-hosts"
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

func resourceCudaWAFTrustedHostsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-host-groups/" + d.Get("parent.0").(string) + "/trusted-hosts"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFTrustedHostsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFTrustedHostsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFTrustedHostsRead(d, m)
}

func resourceCudaWAFTrustedHostsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/trusted-host-groups/" + d.Get("parent.0").(string) + "/trusted-hosts"
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

func hydrateBarracudaWAFTrustedHostsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"ip-address":   d.Get("ip_address").(string),
		"ipv6-address": d.Get("ipv6_address").(string),
		"ipv6-mask":    d.Get("ipv6_mask").(string),
		"mask":         d.Get("mask").(string),
		"comments":     d.Get("comments").(string),
		"name":         d.Get("name").(string),
		"version":      d.Get("version").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFTrustedHostsSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceTrustedHostsParams {
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
