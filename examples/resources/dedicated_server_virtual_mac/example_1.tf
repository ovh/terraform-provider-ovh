resource "ovh_dedicated_server_virtual_mac" "vmac" {
  service_name         = "nsxxxxxxx.ip-xx-xx-xx.eu"
  ip_address           = "1.2.3.4"
  type                 = "ovh"
  virtual_machine_name = "my-vm"
}
