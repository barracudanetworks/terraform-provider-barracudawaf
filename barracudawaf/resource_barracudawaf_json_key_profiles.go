package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceJsonKeyProfilesParams = map[string][]string{}
)

func resourceCudaWAFJsonKeyProfiles() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFJsonKeyProfilesCreate,
		Read:   resourceCudaWAFJsonKeyProfilesRead,
		Update: resourceCudaWAFJsonKeyProfilesUpdate,
		Delete: resourceCudaWAFJsonKeyProfilesDelete,

		Schema: map[string]*schema.Schema{
			"allow_null": {Type: schema.TypeString, Optional: true, Description: "Allow NULL"},
			"allowed_metachars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allowed Metacharacters",
			},
			"base64_decode_parameter_value": {Type: schema.TypeString, Optional: true, Description: "Base64 Decode"},
			"comments":                      {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"value_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Parameter Class",
			},
			"exception_patterns": {Type: schema.TypeString, Optional: true, Description: "Exception Patterns"},
			"key":                {Type: schema.TypeString, Required: true, Description: "Key"},
			"max_array_elements": {Type: schema.TypeString, Optional: true, Description: "Max Array Elements"},
			"max_length":         {Type: schema.TypeString, Optional: true, Description: "Max Length"},
			"max_number_value":   {Type: schema.TypeString, Optional: true, Description: "Max Number Value"},
			"max_keys":           {Type: schema.TypeString, Optional: true, Description: "Max Number Of Keys"},
			"name":               {Type: schema.TypeString, Required: true, Description: "Profile Name"},
			"status":             {Type: schema.TypeString, Optional: true, Description: "Status"},
			"validate_key":       {Type: schema.TypeString, Optional: true, Description: "Validate Key"},
			"value_type":         {Type: schema.TypeString, Optional: true, Description: "Value Type"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_json_key_profiles` manages `Json Key Profiles` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFJsonKeyProfilesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles/" + d.Get("parent.1").(string) + "/json-key-profiles"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFJsonKeyProfilesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFJsonKeyProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFJsonKeyProfilesRead(d, m)
}

func resourceCudaWAFJsonKeyProfilesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles/" + d.Get("parent.1").(string) + "/json-key-profiles"
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

func resourceCudaWAFJsonKeyProfilesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles/" + d.Get("parent.1").(string) + "/json-key-profiles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFJsonKeyProfilesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFJsonKeyProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFJsonKeyProfilesRead(d, m)
}

func resourceCudaWAFJsonKeyProfilesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles/" + d.Get("parent.1").(string) + "/json-key-profiles"
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

func hydrateBarracudaWAFJsonKeyProfilesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"allow-null":                    d.Get("allow_null").(string),
		"allowed-metachars":             d.Get("allowed_metachars").(string),
		"base64-decode-parameter-value": d.Get("base64_decode_parameter_value").(string),
		"comments":                      d.Get("comments").(string),
		"value-class":                   d.Get("value_class").(string),
		"exception-patterns":            d.Get("exception_patterns").(string),
		"key":                           d.Get("key").(string),
		"max-array-elements":            d.Get("max_array_elements").(string),
		"max-length":                    d.Get("max_length").(string),
		"max-number-value":              d.Get("max_number_value").(string),
		"max-keys":                      d.Get("max_keys").(string),
		"name":                          d.Get("name").(string),
		"status":                        d.Get("status").(string),
		"validate-key":                  d.Get("validate_key").(string),
		"value-type":                    d.Get("value_type").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFJsonKeyProfilesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceJsonKeyProfilesParams {
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
