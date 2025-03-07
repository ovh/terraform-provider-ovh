resource "ovh_iam_resource_group" "my_resource_group" {
    name = "my_resource_group"
    resources = [
        "urn:v1:eu:resource:service1:service1-id",
        "urn:v1:eu:resource:service2:service2-id",
    ]
}
