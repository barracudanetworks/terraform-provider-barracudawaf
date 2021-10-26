resource "barracudawaf_servers" "demo_server_1" {
    name            = "DemoServer1"
    identifier      = "IP Address"
    address_version = "IPv4"
    status          = "In Service"
    ip_address      = "x.x.x.x"
    port            = "80"
    comments        = "Creating the Demo Server"
    parent          = [ barracudawaf_services.demo_app_1.name ]

    depends_on      = [ barracudawaf_services.demo_app_2 ]
}