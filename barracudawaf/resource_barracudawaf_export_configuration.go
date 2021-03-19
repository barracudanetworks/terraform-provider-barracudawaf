package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceExportConfigurationParams = map[string][]string{}
)

func resourceCudaWAFExportConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFExportConfigurationCreate,
		Read:   resourceCudaWAFExportConfigurationRead,
		Update: resourceCudaWAFExportConfigurationUpdate,
		Delete: resourceCudaWAFExportConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"backup_type":    {Type: schema.TypeString, Optional: true, Description: "Select the type of JSON export."},
			"name":           {Type: schema.TypeString, Optional: true, Description: "Changes After"},
			"destination":    {Type: schema.TypeString, Optional: true, Description: "Amazon S3 Bucket Name"},
			"day_of_week":    {Type: schema.TypeString, Optional: true, Description: "days"},
			"hour_of_day":    {Type: schema.TypeString, Optional: true, Description: "hour"},
			"minute_of_hour": {Type: schema.TypeString, Optional: true, Description: "minute"},
		},

		Description: "`barracudawaf_export_configuration` manages `Export Configuration` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFExportConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/export-configuration"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFExportConfigurationResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFExportConfigurationSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFExportConfigurationRead(d, m)
}

func resourceCudaWAFExportConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/export-configuration"
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

func resourceCudaWAFExportConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/export-configuration"
	err := client.UpdateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFExportConfigurationResource(d, "put", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFExportConfigurationSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFExportConfigurationRead(d, m)
}

func resourceCudaWAFExportConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/export-configuration"
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

func hydrateBarracudaWAFExportConfigurationResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"backup-type":    d.Get("backup_type").(string),
		"name":           d.Get("name").(string),
		"destination":    d.Get("destination").(string),
		"day-of-week":    d.Get("day_of_week").(string),
		"hour-of-day":    d.Get("hour_of_day").(string),
		"minute-of-hour": d.Get("minute_of_hour").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFExportConfigurationSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceExportConfigurationParams {
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
