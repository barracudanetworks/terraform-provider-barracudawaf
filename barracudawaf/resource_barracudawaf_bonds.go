package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceBondsParams = map[string][]string{}
)

func resourceCudaWAFBonds() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFBondsCreate,
		Read:   resourceCudaWAFBondsRead,
		Update: resourceCudaWAFBondsUpdate,
		Delete: resourceCudaWAFBondsDelete,

		Schema: map[string]*schema.Schema{
			"duplexity":  {Type: schema.TypeString, Optional: true, Description: "None"},
			"name":       {Type: schema.TypeString, Required: true, Description: "Name"},
			"speed":      {Type: schema.TypeString, Optional: true, Description: "None"},
			"bond_ports": {Type: schema.TypeString, Required: true, Description: "Ports"},
			"min_link":   {Type: schema.TypeString, Optional: true, Description: "Minimum Links"},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cannot have VLAN(s) on LAN interface when enabling Network Port Configuration.",
			},
			"mtu": {Type: schema.TypeString, Optional: true, Description: "Maximum Transmission Unit (MTU)"},
		},

		Description: "`barracudawaf_bonds` manages `Bonds` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFBondsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFBondsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFBondsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFBondsRead(d, m)
}

func resourceCudaWAFBondsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
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

func resourceCudaWAFBondsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFBondsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFBondsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFBondsRead(d, m)
}

func resourceCudaWAFBondsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/bonds"
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

func hydrateBarracudaWAFBondsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"duplexity":  d.Get("duplexity").(string),
		"name":       d.Get("name").(string),
		"speed":      d.Get("speed").(string),
		"bond-ports": d.Get("bond_ports").(string),
		"min-link":   d.Get("min_link").(string),
		"mode":       d.Get("mode").(string),
		"mtu":        d.Get("mtu").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"name"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFBondsSubResource(d *schema.ResourceData, name string, endpoint string) error {

	for subResource, subResourceParams := range subResourceBondsParams {
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
