package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceHttpRequestRewriteRulesParams = map[string][]string{}
)

func resourceCudaWAFHttpRequestRewriteRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFHttpRequestRewriteRulesCreate,
		Read:   resourceCudaWAFHttpRequestRewriteRulesRead,
		Update: resourceCudaWAFHttpRequestRewriteRulesUpdate,
		Delete: resourceCudaWAFHttpRequestRewriteRulesDelete,

		Schema: map[string]*schema.Schema{
			"action":              {Type: schema.TypeString, Optional: true, Description: "Action"},
			"comments":            {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"condition":           {Type: schema.TypeString, Optional: true, Description: "Rewrite Condition"},
			"continue_processing": {Type: schema.TypeString, Optional: true, Description: "Continue Processing"},
			"header":              {Type: schema.TypeString, Optional: true, Description: "Header Name"},
			"old_value":           {Type: schema.TypeString, Optional: true, Description: "Old Value"},
			"name":                {Type: schema.TypeString, Required: true, Description: "Rule Name"},
			"sequence_number":     {Type: schema.TypeString, Required: true, Description: "Sequence Number"},
			"rewrite_value":       {Type: schema.TypeString, Optional: true, Description: "Rewrite Value"},
			"parent":              {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_http_request_rewrite_rules` manages `Http Request Rewrite Rules` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFHttpRequestRewriteRulesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/http-request-rewrite-rules"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFHttpRequestRewriteRulesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFHttpRequestRewriteRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFHttpRequestRewriteRulesRead(d, m)
}

func resourceCudaWAFHttpRequestRewriteRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/http-request-rewrite-rules"
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

func resourceCudaWAFHttpRequestRewriteRulesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/http-request-rewrite-rules"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFHttpRequestRewriteRulesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFHttpRequestRewriteRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFHttpRequestRewriteRulesRead(d, m)
}

func resourceCudaWAFHttpRequestRewriteRulesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/http-request-rewrite-rules"
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

func hydrateBarracudaWAFHttpRequestRewriteRulesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"action":              d.Get("action").(string),
		"comments":            d.Get("comments").(string),
		"condition":           d.Get("condition").(string),
		"continue-processing": d.Get("continue_processing").(string),
		"header":              d.Get("header").(string),
		"old-value":           d.Get("old_value").(string),
		"name":                d.Get("name").(string),
		"sequence-number":     d.Get("sequence_number").(string),
		"rewrite-value":       d.Get("rewrite_value").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFHttpRequestRewriteRulesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceHttpRequestRewriteRulesParams {
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
