data "ovh_dedicated_server" "server" {
  service_name = "nsxxxxxxx.ip-xx-xx-xx.eu"
}

data "ovh_dedicated_installation_template" "template" {
  template_name = "debian12_64"
}

resource "ovh_dedicated_server_reinstall_task" "server_install" {
  service_name = data.ovh_dedicated_server.server.service_name
  os           = "win2022core-std_64"
  customizations {
    hostname                 = "ma-fenetre"
    language                 = "fr-fr"
    post_installation_script = "ImNvdWNvdSBwb3N0SW5zdGFsbGF0aW9uU2NyaXB0UG93ZXJTaGVsbCIgfCBPdXQtRmlsZSAtRmlsZVBhdGggImM6XG92aHVwZFxzY3JpcHRcY291Y291LnR4dCIKKEdldC1JdGVtUHJvcGVydHkgLUxpdGVyYWxQYXRoICJSZWdpc3RyeTo6SEtMTVxTT0ZUV0FSRVxNaWNyb3NvZnRcQ3J5cHRvZ3JhcGh5IiAtTmFtZSAiTWFjaGluZUd1aWQiKS5NYWNoaW5lR3VpZCB8IE91dC1GaWxlIC1GaWxlUGF0aCAiYzpcb3ZodXBkXHNjcmlwdFxjb3Vjb3UudHh0IiAtQXBwZW5kCihHZXQtRGF0ZSkuVG9Vbml2ZXJzYWxUaW1lKCkuVG9TdHJpbmcoInl5eXktTU0tZGQgSEg6bW06c3MiKSB8IE91dC1GaWxlIC1GaWxlUGF0aCAiYzpcb3ZodXBkXHNjcmlwdFxjb3Vjb3UudHh0IiAtQXBwZW5kCg=="
  }
}
