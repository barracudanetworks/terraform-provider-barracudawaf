resource "barracudawaf_security_policies" "demo_security_policy_1" {
    name       = "DemoPolicy1"
    based_on   = "Create New"
    
    depends_on = [ barracudawaf_servers.demo_server_1 ]
}