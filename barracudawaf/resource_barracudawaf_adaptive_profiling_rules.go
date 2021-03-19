package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceAdaptiveProfilingRulesParams = map[string][]string{}
)

func resourceCudaWAFAdaptiveProfilingRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAdaptiveProfilingRulesCreate,
		Read:   resourceCudaWAFAdaptiveProfilingRulesRead,
		Update: resourceCudaWAFAdaptiveProfilingRulesUpdate,
		Delete: resourceCudaWAFAdaptiveProfilingRulesDelete,

		Schema: map[string]*schema.Schema{
			"host":                {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"learn_from_request":  {Type: schema.TypeString, Optional: true, Description: "Learn From Request"},
			"learn_from_response": {Type: schema.TypeString, Optional: true, Description: "Learn From Response"},
			"name":                {Type: schema.TypeString, Required: true, Description: "Learn Rule Name"},
			"status":              {Type: schema.TypeString, Optional: true, Description: "Status"},
			"url":                 {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"parent":              {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_adaptive_profiling_rules` manages `Adaptive Profiling Rules` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAdaptiveProfilingRulesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/adaptive-profiling-rules"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFAdaptiveProfilingRulesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAdaptiveProfilingRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAdaptiveProfilingRulesRead(d, m)
}

func resourceCudaWAFAdaptiveProfilingRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/adaptive-profiling-rules"
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

func resourceCudaWAFAdaptiveProfilingRulesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/adaptive-profiling-rules"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFAdaptiveProfilingRulesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFAdaptiveProfilingRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAdaptiveProfilingRulesRead(d, m)
}

func resourceCudaWAFAdaptiveProfilingRulesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/adaptive-profiling-rules"
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

func hydrateBarracudaWAFAdaptiveProfilingRulesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"host":                d.Get("host").(string),
		"learn-from-request":  d.Get("learn_from_request").(string),
		"learn-from-response": d.Get("learn_from_response").(string),
		"name":                d.Get("name").(string),
		"status":              d.Get("status").(string),
		"url":                 d.Get("url").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFAdaptiveProfilingRulesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceAdaptiveProfilingRulesParams {
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
