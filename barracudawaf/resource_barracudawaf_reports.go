package barracudawaf

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceReportsParams = map[string][]string{}
)

func resourceCudaWAFReports() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFReportsCreate,
		Read:   resourceCudaWAFReportsRead,
		Update: resourceCudaWAFReportsUpdate,
		Delete: resourceCudaWAFReportsDelete,

		Schema: map[string]*schema.Schema{
			"report_format":    {Type: schema.TypeString, Optional: true, Description: "None"},
			"frequency":        {Type: schema.TypeString, Optional: true, Description: "None"},
			"ftp_directory":    {Type: schema.TypeString, Optional: true, Description: "Folder/Path"},
			"ftp_ip_address":   {Type: schema.TypeString, Optional: true, Description: "Server Name/IP"},
			"ftp_password":     {Type: schema.TypeString, Optional: true, Description: "Password"},
			"ftp_port":         {Type: schema.TypeString, Optional: true, Description: "Port"},
			"ftp_username":     {Type: schema.TypeString, Optional: true, Description: "Username"},
			"email_id":         {Type: schema.TypeString, Optional: true, Description: "Email Report to:"},
			"name":             {Type: schema.TypeString, Required: true, Description: "None"},
			"report_types":     {Type: schema.TypeString, Required: true, Description: "Report Type"},
			"delivery_options": {Type: schema.TypeString, Optional: true, Description: "Chart Type"},
		},

		Description: "`barracudawaf_reports` manages `Reports` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFReportsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/reports"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFReportsResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFReportsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFReportsRead(d, m)
}

func resourceCudaWAFReportsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/reports"
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

func resourceCudaWAFReportsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/reports"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFReportsResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFReportsSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFReportsRead(d, m)
}

func resourceCudaWAFReportsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/reports"
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

func hydrateBarracudaWAFReportsResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"report-format":    d.Get("report_format").(string),
		"frequency":        d.Get("frequency").(string),
		"ftp-directory":    d.Get("ftp_directory").(string),
		"ftp-ip-address":   d.Get("ftp_ip_address").(string),
		"ftp-password":     d.Get("ftp_password").(string),
		"ftp-port":         d.Get("ftp_port").(string),
		"ftp-username":     d.Get("ftp_username").(string),
		"email-id":         d.Get("email_id").(string),
		"name":             d.Get("name").(string),
		"report-types":     d.Get("report_types").(string),
		"delivery-options": d.Get("delivery_options").(string),
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

func (b *BarracudaWAF) hydrateBarracudaWAFReportsSubResource(d *schema.ResourceData, name string, endpoint string) error {

	for subResource, subResourceParams := range subResourceReportsParams {
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
