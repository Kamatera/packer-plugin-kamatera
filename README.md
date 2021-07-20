# packer-plugin-kamatera

The `Kamatera` multi-component plugin can be used with HashiCorp [Packer](https://www.packer.io)
to create custom images.

## Installation

Assuming you have [installed Packer v1.7 or above](https://www.packer.io/downloads) you can install the plugin using the Packer init command:

Add the following code into your Packer configuration file:

```
packer {
  required_plugins {
    kamatera = {
      version = ">= 0.3.3"
      source  = "github.com/kamatera/kamatera"
    }
  }
}
```

Run `packer init` to install the plugin.

