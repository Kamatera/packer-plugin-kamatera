# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Kamatera"
  description = "Build Kamatera cloud images and add them to the account's disk library."
  identifier = "packer/Kamatera/kamatera"
  component {
    type = "builder"
    name = "Kamatera"
    slug = "kamatera"
  }
}
