resource "barracudawaf_security_policies" "security_policy_1" {
    name        = "DemoSecurityPolicy1"
    based_on    = "Create New"
    depends_on  = [barracudawaf_servers.web_server_1]
}