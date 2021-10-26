resource "barracudawaf_content_rule_servers" "demo_rule_group_server_1" {
    name        = "DemoRuleGroupServer1"
    identifier  = "Hostname"
    hostname    = "barracuda.com"
    parent      = [ barracudawaf_services.demo_app_1.name, barracudawaf_content_rules.demo_rule_group_1.name ]

    depends_on = [ barracudawaf_content_rules.demo_rule_group_1 ]
}