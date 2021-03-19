package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceCustomParameterClassesParams = map[string][]string{}
)

func resourceCudaWAFCustomParameterClasses() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFCustomParameterClassesCreate,
		Read:   resourceCudaWAFCustomParameterClassesRead,
		Update: resourceCudaWAFCustomParameterClassesUpdate,
		Delete: resourceCudaWAFCustomParameterClassesDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Parameter Class Name",
			},
			"custom_blocked_attack_types": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Blocked Attack Types",
			},
			"custom_input_type_validation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Input Type Validation",
			},
			"denied_metacharacters": {Type: schema.TypeString, Optional: true, Description: "Denied Metacharacters"},
			"input_type_validation": {Type: schema.TypeString, Optional: true, Description: "Input Type Validation"},
			"blocked_attack_types":  {Type: schema.TypeString, Optional: true, Description: "Blocked Attack Types"},
		},

		Description: "`barracudawaf_custom_parameter_classes` manages `Custom Parameter Classes` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFCustomParameterClassesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/custom-parameter-classes"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFCustomParameterClassesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFCustomParameterClassesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFCustomParameterClassesRead(d, m)
}

func resourceCudaWAFCustomParameterClassesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/custom-parameter-classes"
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

func resourceCudaWAFCustomParameterClassesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/custom-parameter-classes"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFCustomParameterClassesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFCustomParameterClassesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFCustomParameterClassesRead(d, m)
}

func resourceCudaWAFCustomParameterClassesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/custom-parameter-classes"
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

func hydrateBarracudaWAFCustomParameterClassesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                         d.Get("name").(string),
		"custom-blocked-attack-types":  d.Get("custom_blocked_attack_types").(string),
		"custom-input-type-validation": d.Get("custom_input_type_validation").(string),
		"denied-metacharacters":        d.Get("denied_metacharacters").(string),
		"input-type-validation":        d.Get("input_type_validation").(string),
		"blocked-attack-types":         d.Get("blocked_attack_types").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFCustomParameterClassesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceCustomParameterClassesParams {
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
