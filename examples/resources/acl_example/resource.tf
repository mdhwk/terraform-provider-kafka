resource "kafka_acls" "topic_knative_broker_prefixed" {
  resource_name                = "<topic_name>"
  resource_type                = "Topic"
  resource_pattern_type_filter = "Prefixed"
  acl_principal = "CN=example_user"
  acl_operation = "Create"
  acl_host = "*"
}