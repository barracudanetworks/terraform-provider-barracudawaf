package barracudawaf

import (
	"fmt"
	"log"
)

//Config : container for Barracuda WAF session
type Config struct {
	IPAddress string
	Username  string
	Password  string
	AdminPort string
}

//Client : Barracuda WAF Client for REST API calls for resource crud
func (c *Config) Client() (*BarracudaWAF, error) {

	if c.IPAddress != "" && c.Username != "" && c.Password != "" && c.AdminPort != "" {
		log.Println("[INFO] Initializing Barracuda WAF connection")
		var client *BarracudaWAF

		client = NewSession(c.IPAddress, c.AdminPort, c.Username, c.Password)
		client, err := c.validateConnection(client)
		if err == nil {
			return client, nil
		}
		return nil, err
	}
	return nil, fmt.Errorf("Barracuda WAF provider requires IPAddress, Username, Password and AdminPort")
}

func (c *Config) validateConnection(client *BarracudaWAF) (*BarracudaWAF, error) {

	client, err := client.GetAuthToken()
	if err != nil {
		log.Printf("[ERROR] Connection to Barracuda WAF could not have been validated: %v ", err)
		return client, err
	}

	return client, nil
}
