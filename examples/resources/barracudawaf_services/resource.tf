resource "barracudawaf_services" "application_1" {
    name            = "DemoApp1"
    ip_address      = "x.x.x.x"
    port            = "80"
    type            = "HTTP"
    vsite           = "default"
    address_version = "IPv4"
    status          = "On"
    group           = "default"
    comments        = "Demo App for Terraform"
}