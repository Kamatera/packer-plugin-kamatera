The Kamatera plugin is intended to be used for managing [Kamatera cloud](https://www.kamatera.com/express/compute/) images through Packer

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    kamatera = {
      version = ">= 0.5.0"
      source  = "github.com/kamatera/kamatera"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/kamatera/kamatera
```


### Components

#### Builders

- [kamatera](/packer/integrations/Kamatera/kamatera/latest/components/builder/kamatera) - The builder builds Kamatera images.
