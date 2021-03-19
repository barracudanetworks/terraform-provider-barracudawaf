package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceParameterProfilesParams = map[string][]string{}
)

func resourceCudaWAFParameterProfiles() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFParameterProfilesCreate,
		Read:   resourceCudaWAFParameterProfilesRead,
		Update: resourceCudaWAFParameterProfilesUpdate,
		Delete: resourceCudaWAFParameterProfilesDelete,

		Schema: map[string]*schema.Schema{
			"allowed_metachars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allowed Metacharacters",
			},
			"base64_decode_parameter_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base64 Decode Parameter Value",
			},
			"allowed_file_upload_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allowed File Upload Type",
			},
			"comments": {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"custom_parameter_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Parameter Class",
			},
			"exception_patterns": {Type: schema.TypeString, Optional: true, Description: "Exception Patterns"},
			"file_upload_extensions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File Upload Extensions",
			},
			"file_upload_mime_types": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File Upload Mime Types",
			},
			"ignore":            {Type: schema.TypeString, Optional: true, Description: "Ignore"},
			"maximum_instances": {Type: schema.TypeString, Optional: true, Description: "Maximum Instances"},
			"max_value_length":  {Type: schema.TypeString, Optional: true, Description: "Max Value Length"},
			"parameter":         {Type: schema.TypeString, Required: true, Description: "Parameter"},
			"parameter_class":   {Type: schema.TypeString, Required: true, Description: "Parameter Class"},
			"required":          {Type: schema.TypeString, Optional: true, Description: "Required"},
			"status":            {Type: schema.TypeString, Optional: true, Description: "Status"},
			"type":              {Type: schema.TypeString, Optional: true, Description: "Type"},
			"validate_parameter_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Validate Parameter Name",
			},
			"values": {Type: schema.TypeString, Optional: true, Description: "Values"},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Parameter Profile Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_parameter_profiles` manages `Parameter Profiles` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFParameterProfilesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles/" + d.Get("parent.1").(string) + "/parameter-profiles"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFParameterProfilesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFParameterProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFParameterProfilesRead(d, m)
}

func resourceCudaWAFParameterProfilesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles/" + d.Get("parent.1").(string) + "/parameter-profiles"
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

func resourceCudaWAFParameterProfilesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles/" + d.Get("parent.1").(string) + "/parameter-profiles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFParameterProfilesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFParameterProfilesSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFParameterProfilesRead(d, m)
}

func resourceCudaWAFParameterProfilesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/services/" + d.Get("parent.0").(string) + "/url-profiles/" + d.Get("parent.1").(string) + "/parameter-profiles"
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

func hydrateBarracudaWAFParameterProfilesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"allowed-metachars":             d.Get("allowed_metachars").(string),
		"base64-decode-parameter-value": d.Get("base64_decode_parameter_value").(string),
		"allowed-file-upload-type":      d.Get("allowed_file_upload_type").(string),
		"comments":                      d.Get("comments").(string),
		"custom-parameter-class":        d.Get("custom_parameter_class").(string),
		"exception-patterns":            d.Get("exception_patterns").(string),
		"file-upload-extensions":        d.Get("file_upload_extensions").(string),
		"file-upload-mime-types":        d.Get("file_upload_mime_types").(string),
		"ignore":                        d.Get("ignore").(string),
		"maximum-instances":             d.Get("maximum_instances").(string),
		"max-value-length":              d.Get("max_value_length").(string),
		"parameter":                     d.Get("parameter").(string),
		"parameter-class":               d.Get("parameter_class").(string),
		"required":                      d.Get("required").(string),
		"status":                        d.Get("status").(string),
		"type":                          d.Get("type").(string),
		"validate-parameter-name":       d.Get("validate_parameter_name").(string),
		"values":                        d.Get("values").(string),
		"name":                          d.Get("name").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFParameterProfilesSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceParameterProfilesParams {
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
