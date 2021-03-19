package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceUrlProfilesParams = map[string][]string{}
)

func resourceCudaWAFUrlProfiles() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFUrlProfilesCreate,
		Read:   resourceCudaWAFUrlProfilesRead,
		Update: resourceCudaWAFUrlProfilesUpdate,
		Delete: resourceCudaWAFUrlProfilesDelete,

		Schema: map[string]*schema.Schema{
			"allowed_content_types": {Type: schema.TypeString, Optional: true, Description: "Allowed Content Types"},
			"allowed_methods":       {Type: schema.TypeString, Optional: true, Description: "Allowed Methods"},
			"custom_blocked_attack_types": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Blocked Attack Types",
			},
			"comment":            {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"display_name":       {Type: schema.TypeString, Optional: true, Description: "Display Name"},
			"exception_patterns": {Type: schema.TypeString, Optional: true, Description: "Exception Patterns"},
			"extended_match":     {Type: schema.TypeString, Optional: true, Description: "Extended Match"},
			"extended_match_sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended Match Sequence",
			},
			"hidden_parameter_protection": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hidden Parameter Protection",
			},
			"blocked_attack_types": {Type: schema.TypeString, Optional: true, Description: "Blocked Attack Types"},
			"max_content_length":   {Type: schema.TypeString, Optional: true, Description: "Max Content Length"},
			"maximum_parameter_name_length": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Maximum Parameter Name Length",
			},
			"maximum_upload_files": {Type: schema.TypeString, Optional: true, Description: "Maximum Upload Files"},
			"minimum_form_fill_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Minimum Form Fill Time",
			},
			"name":               {Type: schema.TypeString, Required: true, Description: "URL Profile Name"},
			"csrf_prevention":    {Type: schema.TypeString, Optional: true, Description: "CSRF Prevention"},
			"allow_query_string": {Type: schema.TypeString, Optional: true, Description: "Allow Query String"},
			"referrers_for_the_url_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Referrers for the URL Profile",
			},
			"mode":   {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"status": {Type: schema.TypeString, Optional: true, Description: "Status"},
			"url":    {Type: schema.TypeString, Required: true, Description: "URL"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_url_profiles` manages `Url Profiles` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFUrlProfilesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFUrlProfilesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFUrlProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFUrlProfilesRead(d, m)
}

func resourceCudaWAFUrlProfilesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles"
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

func resourceCudaWAFUrlProfilesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFUrlProfilesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFUrlProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFUrlProfilesRead(d, m)
}

func resourceCudaWAFUrlProfilesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles"
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

func hydrateBarracudaWAFUrlProfilesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"allowed-content-types":         d.Get("allowed_content_types").(string),
		"allowed-methods":               d.Get("allowed_methods").(string),
		"custom-blocked-attack-types":   d.Get("custom_blocked_attack_types").(string),
		"comment":                       d.Get("comment").(string),
		"display-name":                  d.Get("display_name").(string),
		"exception-patterns":            d.Get("exception_patterns").(string),
		"extended-match":                d.Get("extended_match").(string),
		"extended-match-sequence":       d.Get("extended_match_sequence").(string),
		"hidden-parameter-protection":   d.Get("hidden_parameter_protection").(string),
		"blocked-attack-types":          d.Get("blocked_attack_types").(string),
		"max-content-length":            d.Get("max_content_length").(string),
		"maximum-parameter-name-length": d.Get("maximum_parameter_name_length").(string),
		"maximum-upload-files":          d.Get("maximum_upload_files").(string),
		"minimum-form-fill-time":        d.Get("minimum_form_fill_time").(string),
		"name":                          d.Get("name").(string),
		"csrf-prevention":               d.Get("csrf_prevention").(string),
		"allow-query-string":            d.Get("allow_query_string").(string),
		"referrers-for-the-url-profile": d.Get("referrers_for_the_url_profile").(string),
		"mode":                          d.Get("mode").(string),
		"status":                        d.Get("status").(string),
		"url":                           d.Get("url").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFUrlProfilesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceUrlProfilesParams {
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
