resource "barracudawaf_content_rules" "demo_rule_group_1" {
    name                = "DemoRuleGroup1"
    url_match           = "/index.html"
    host_match          = "www.example.com"
    web_firewall_policy = "DemoPolicy1"
    mode                = "Active"
    parent              = [ barracudawaf_services.demo_app_1.name ]
    
    depends_on          = [ barracudawaf_security_policies.demo_security_policy_1 ]
}