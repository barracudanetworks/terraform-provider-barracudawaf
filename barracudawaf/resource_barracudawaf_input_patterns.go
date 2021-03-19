package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceInputPatternsParams = map[string][]string{}
)

func resourceCudaWAFInputPatterns() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFInputPatternsCreate,
		Read:   resourceCudaWAFInputPatternsRead,
		Update: resourceCudaWAFInputPatternsUpdate,
		Delete: resourceCudaWAFInputPatternsDelete,

		Schema: map[string]*schema.Schema{
			"algorithm":      {Type: schema.TypeString, Optional: true, Description: "Pattern Algorithm"},
			"case_sensitive": {Type: schema.TypeString, Optional: true, Description: "Case Sensitivity"},
			"description":    {Type: schema.TypeString, Optional: true, Description: "Pattern Description"},
			"name":           {Type: schema.TypeString, Required: true, Description: "Pattern Name"},
			"regex":          {Type: schema.TypeString, Required: true, Description: "Pattern Regex"},
			"status":         {Type: schema.TypeString, Optional: true, Description: "Status"},
			"parent":         {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_input_patterns` manages `Input Patterns` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFInputPatternsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/input-types/" + d.Get("parent.0").(string) + "/input-patterns"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFInputPatternsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFInputPatternsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFInputPatternsRead(d, m)
}

func resourceCudaWAFInputPatternsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/input-types/" + d.Get("parent.0").(string) + "/input-patterns"
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

func resourceCudaWAFInputPatternsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/input-types/" + d.Get("parent.0").(string) + "/input-patterns"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFInputPatternsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFInputPatternsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFInputPatternsRead(d, m)
}

func resourceCudaWAFInputPatternsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/input-types/" + d.Get("parent.0").(string) + "/input-patterns"
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

func hydrateBarracudaWAFInputPatternsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"algorithm":      d.Get("algorithm").(string),
		"case-sensitive": d.Get("case_sensitive").(string),
		"description":    d.Get("description").(string),
		"name":           d.Get("name").(string),
		"regex":          d.Get("regex").(string),
		"status":         d.Get("status").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFInputPatternsSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceInputPatternsParams {
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
