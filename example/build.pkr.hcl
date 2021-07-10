//packer {
//  required_plugins {
//    kamatera = {
//      version = ">=v0.1.0"
//      source  = "github.com/Kamatera/kamatera"
//    }
//  }
//}

source "kamatera" "test" {
  # api_client_id = ""
  # api_secret = ""

  ssh_username = "root"

  datacenter = "IL"
  //  cpu = "1A"
  //  ram = "1024"
  //  image = "ubuntu_server_18.04_64-bit"
  //  image_name =
}

build {
  sources = ["source.kamatera.test"]
}
