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


##  with Multiple Network Interfaces

resource "ovh_cloud_project_network_private" "net_front" {
  service_name = "XXX"
  name         = "public-front"
  regions      = ["GRA11"]
}

resource "ovh_cloud_project_subnet" "subnet_front" {
  service_name = "XXX"
  network_id   = ovh_cloud_project_network_private.net_front.id
  region       = "GRA11"
  network      = "10.0.1.0/24"
  dhcp         = true
}

resource "ovh_cloud_project_network_private" "net_back" {
  service_name = "XXX"
  name         = "private-back"
  regions      = ["GRA11"]
}

resource "ovh_cloud_project_subnet" "subnet_back" {
  service_name = "XXX"
  network_id   = private.net_back.id
  region       = "GRA11"
  network      = "10.0.2.0/24"
  dhcp         = true
}

resource "ovh_cloud_project_instance" "multinic_instance" {
  service_name   = "XXX"
  region         = "GRA11"
  name           = "multi-nic-instance"
  flavor_id      = "UUID"
  image_id       = "UUID"

  network {
    public = true

    private {
      network {
        id        = private.net_front.id
        subnet_id = subnet.subnet_front.id
      }
      ip = "10.0.1.10"
    }

    private {
      network {
        id        = private.net_back.id
        subnet_id = subnet.subnet_back.id
      }
    }
  }
}
