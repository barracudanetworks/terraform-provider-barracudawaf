package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceFormSpamFormsParams = map[string][]string{}
)

func resourceCudaWAFFormSpamForms() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFFormSpamFormsCreate,
		Read:   resourceCudaWAFFormSpamFormsRead,
		Update: resourceCudaWAFFormSpamFormsUpdate,
		Delete: resourceCudaWAFFormSpamFormsDelete,

		Schema: map[string]*schema.Schema{
			"name":                   {Type: schema.TypeString, Required: true, Description: "Form Name"},
			"created_by":             {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"status":                 {Type: schema.TypeString, Optional: true, Description: "Status"},
			"mode":                   {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"action_url":             {Type: schema.TypeString, Required: true, Description: "Action URL"},
			"minimum_form_fill_time": {Type: schema.TypeString, Optional: true, Description: "Minimum Form Fill Time"},
			"parameter_name":         {Type: schema.TypeString, Optional: true, Description: "Parameter Name"},
			"parameter_class":        {Type: schema.TypeString, Optional: true, Description: "Parameter Class"},
			"parent":                 {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_form_spam_forms` manages `Form Spam Forms` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFFormSpamFormsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/form-spam-forms"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFFormSpamFormsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFFormSpamFormsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFFormSpamFormsRead(d, m)
}

func resourceCudaWAFFormSpamFormsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/form-spam-forms"
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

func resourceCudaWAFFormSpamFormsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/form-spam-forms"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFFormSpamFormsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFFormSpamFormsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFFormSpamFormsRead(d, m)
}

func resourceCudaWAFFormSpamFormsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/form-spam-forms"
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

func hydrateBarracudaWAFFormSpamFormsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                   d.Get("name").(string),
		"created-by":             d.Get("created_by").(string),
		"status":                 d.Get("status").(string),
		"mode":                   d.Get("mode").(string),
		"action-url":             d.Get("action_url").(string),
		"minimum-form-fill-time": d.Get("minimum_form_fill_time").(string),
		"parameter-name":         d.Get("parameter_name").(string),
		"parameter-class":        d.Get("parameter_class").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"created-by", "action-url"}
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

func (b *BarracudaWAF) hydrateBarracudaWAFFormSpamFormsSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceFormSpamFormsParams {
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
