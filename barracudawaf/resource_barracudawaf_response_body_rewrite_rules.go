package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceResponseBodyRewriteRulesParams = map[string][]string{}
)

func resourceCudaWAFResponseBodyRewriteRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFResponseBodyRewriteRulesCreate,
		Read:   resourceCudaWAFResponseBodyRewriteRulesRead,
		Update: resourceCudaWAFResponseBodyRewriteRulesUpdate,
		Delete: resourceCudaWAFResponseBodyRewriteRulesDelete,

		Schema: map[string]*schema.Schema{
			"comments":        {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"host":            {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"name":            {Type: schema.TypeString, Required: true, Description: "Rule Name"},
			"replace_string":  {Type: schema.TypeString, Optional: true, Description: "Replace String"},
			"search_string":   {Type: schema.TypeString, Required: true, Description: "Search String"},
			"sequence_number": {Type: schema.TypeString, Required: true, Description: "Sequence number"},
			"url":             {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"parent":          {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_response_body_rewrite_rules` manages `Response Body Rewrite Rules` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFResponseBodyRewriteRulesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/response-body-rewrite-rules"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFResponseBodyRewriteRulesResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFResponseBodyRewriteRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFResponseBodyRewriteRulesRead(d, m)
}

func resourceCudaWAFResponseBodyRewriteRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/response-body-rewrite-rules"
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

func resourceCudaWAFResponseBodyRewriteRulesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/response-body-rewrite-rules"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFResponseBodyRewriteRulesResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFResponseBodyRewriteRulesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFResponseBodyRewriteRulesRead(d, m)
}

func resourceCudaWAFResponseBodyRewriteRulesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/response-body-rewrite-rules"
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

func hydrateBarracudaWAFResponseBodyRewriteRulesResource(
	d *schema.ResourceData,
	method string,
	endpoint string,
) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"comments":        d.Get("comments").(string),
		"host":            d.Get("host").(string),
		"name":            d.Get("name").(string),
		"replace-string":  d.Get("replace_string").(string),
		"search-string":   d.Get("search_string").(string),
		"sequence-number": d.Get("sequence_number").(string),
		"url":             d.Get("url").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFResponseBodyRewriteRulesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceResponseBodyRewriteRulesParams {
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
