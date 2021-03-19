package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceJsonProfilesParams = map[string][]string{}
)

func resourceCudaWAFJsonProfiles() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFJsonProfilesCreate,
		Read:   resourceCudaWAFJsonProfilesRead,
		Update: resourceCudaWAFJsonProfilesUpdate,
		Delete: resourceCudaWAFJsonProfilesDelete,

		Schema: map[string]*schema.Schema{
			"allowed_content_types": {Type: schema.TypeString, Optional: true, Description: "Inspect Mime Types"},
			"host_match":            {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"ignore_keys":           {Type: schema.TypeString, Optional: true, Description: "Ignore Keys"},
			"method":                {Type: schema.TypeString, Required: true, Description: "Methods"},
			"comment":               {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"name":                  {Type: schema.TypeString, Required: true, Description: "JSON Profile Name"},
			"mode":                  {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"status":                {Type: schema.TypeString, Optional: true, Description: "Status"},
			"url_match":             {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"exception_patterns":    {Type: schema.TypeString, Optional: true, Description: "None"},
			"json_policy":           {Type: schema.TypeString, Optional: true, Description: "JSON Policy"},
			"validate_key":          {Type: schema.TypeString, Optional: true, Description: "Validate Key"},
			"parent":                {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}, Required: true},
		},

		Description: "`barracudawaf_json_profiles` manages `Json Profiles` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFJsonProfilesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFJsonProfilesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFJsonProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFJsonProfilesRead(d, m)
}

func resourceCudaWAFJsonProfilesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles"
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

func resourceCudaWAFJsonProfilesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFJsonProfilesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFJsonProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFJsonProfilesRead(d, m)
}

func resourceCudaWAFJsonProfilesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/json-profiles"
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

func hydrateBarracudaWAFJsonProfilesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"allowed-content-types": d.Get("allowed_content_types").(string),
		"host_match":            d.Get("host_match").(string),
		"ignore-keys":           d.Get("ignore_keys").(string),
		"method":                d.Get("method").(string),
		"comment":               d.Get("comment").(string),
		"name":                  d.Get("name").(string),
		"mode":                  d.Get("mode").(string),
		"status":                d.Get("status").(string),
		"url_match":             d.Get("url_match").(string),
		"exception-patterns":    d.Get("exception_patterns").(string),
		"json_policy":           d.Get("json_policy").(string),
		"validate_key":          d.Get("validate_key").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFJsonProfilesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceJsonProfilesParams {
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
