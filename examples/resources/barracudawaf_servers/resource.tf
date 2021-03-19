resource "barracudawaf_servers" "web_server_1" {
    name            = "DemoWebServer1"
    identifier      = "IP Address"
    address_version = "IPv4"
    ip_address      = "x.x.x.x"
    port            = "80"
    comments        = "Demo web server behind DemoApp1"
    parent          = [ "DemoApp1" ]
    depends_on      = [barracudawaf_services.application_1]
}