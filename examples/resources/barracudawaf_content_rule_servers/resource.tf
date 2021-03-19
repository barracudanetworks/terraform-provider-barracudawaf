resource "barracudawaf_content_rule_servers" "rule_group_server_1" {
    name       = "DemoRgServer1"
    identifier = "Hostname"
    hostname   = "www.example.com"
    parent     = [ "DemoApp1", "DemoRuleGroup1" ]
    depends_on = [barracudawaf_content_rules.rule_group_1]
}