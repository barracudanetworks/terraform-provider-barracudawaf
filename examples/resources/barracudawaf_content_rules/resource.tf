resource "barracudawaf_content_rules" "rule_group_1" {
    name                = "DemoRuleGroup1"
    url_match           = "/testing.html"
    host_match          = "www.example.com"
    web_firewall_policy = "DemoSecurityPolicy1"
    parent              = [ "DemoApp1" ]
    depends_on          = [barracudawaf_security_policies.security_policy_1]
}