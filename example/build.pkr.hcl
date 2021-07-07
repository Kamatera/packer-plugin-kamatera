//packer {
//  required_plugins {
//    kamatera = {
//      version = ">=v0.1.0"
//      source  = "github.com/Kamatera/kamatera"
//    }
//  }
//}

source "kamatera" "test" {
  api_client_id = "d66f959d724b45bddad8750b5fd5e728"
  api_secret = "f0396ba3b682e764594668aae6cfc524"

  ssh_username = "root"

  //  server_name = "packer_test"
  //  datacenter = "IL"
  //  cpu = "1A"
  //  ram = "1024"
  //  image = "ubuntu_server_18.04_64-bit"
  //  password = "__generate__"
  //  snapshot_name =
}

build {
  sources = ["source.kamatera.test"]
}
