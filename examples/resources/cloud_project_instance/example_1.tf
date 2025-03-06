resource "ovh_cloud_project_instance" "instance" {
  service_name  = "XXX"
    region = "RRRR"
    billing_period = "hourly"
    boot_from {
        image_id = "UUID"
    }
    flavor {
        flavor_id = "UUID"
    }
    name = "instance name"
    ssh_key {
        name = "sshname"
    }
    network {
        public = true
    }  
}
