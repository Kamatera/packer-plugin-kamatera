packer {
  required_plugins {
    kamatera = {
      version = ">=v0.1.0"
      source  = "github.com/Kamatera/kamatera"
    }
  }
}

source "kamatera-kamatera-builder" "test" {
  api_client_id = "3ee67c7e4287f33fed57667b8c03e5cf"
  api_secret = "3f611d8360b39303258c9402607f2c66"

  ssh_username = "root"

  server_name = "packer_test"
  //  datacenter = ""
}

build {
  sources = ["source.kamatera-kamatera-builder.test"]
}
