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
      version = ">= 0.1.6"
      source  = "github.com/kamatera/kamatera"
    }
  }
}
```

Run `packer init` to install the plugin.

## Configuration Reference

### Required Builder Configuration options:

- `api_secret` (string) - The Kamatera API secret to use. It can also be specified via environment variable `KAMATERA_API_SECRET`.
- `api_client_id` (string) - The Kamatera API client ID to use. It can also be specified via environment variable `KAMATERA_API_CLIENT_ID`.
- `datacenter` (string) - The Kamatera datacenter to which the server is deployed.
- `ssh_username` (string) - Must be set with a value of `root`

### Optional:

- `snapshot_name` (string) - The name of the resulting image that will appear in your Kamatera hard disk library. Defaults to `packer-{{timestamp}}`.
- `cpu` (string) - The cpu assigned to the server. Default to `1A`.
- `ram` (string) - The ram assigned to the server. Default to `1024`.
- `image` (string) - The image used by the server. Default to `ubuntu_server_18.04_64-bit`.
- `disk` (string) - The disk size in GB assigned to the server. Default to `10`.

## Tutorial

Set environment variables with your Kamatera API key (you can generate an API key from the [Kamatera console]() -> API -> Keys):

```
export KAMATERA_API_CLIENT_ID=
export KAMATERA_API_SECRET=
```

Create a new directory named `packer_tutorial` and paste the following configuration into a file named `kamatera-ubuntu.pkr.hcl`:

```
packer {
  required_plugins {
    kamatera = {
      version = ">= 0.1.6"
      source  = "github.com/kamatera/kamatera"
    }
  }
}

source "kamatera" "ubuntu" {
  datacenter = "EU"
  ssh_username = "root"
  snapshot_name = "ubuntu-example-packer-image"
}

build {
  sources = [
    "source.kamatera.ubuntu"
  ]
  provisioner "shell" {
    environment_vars = [
      "FOO=hello world",
    ]
    inline = [
      "echo Creating example.txt",
      "echo \"FOO is $FOO\" > example.txt",
    ]
  }
}
```

Change to the `packer_tutorial` directory and run the following commands to download the Kamatera plugin:

```
packer init .
```

Build the image:

```
packer build .
```

Please wait, this will take a while...

When done, log-in to [Kamatera console](https://console.kamatera.com) and navigate to My Cloud -> Hard Disk Library.
Choose zone `EU` and click on MY PRIVATE IMAGES - you should see the created image there.

You can now create a new server from this image, SSH into that server and run `cat example.txt`.
The output should be `FOO is hello world` - as it was generated in the provisioning script in your pkr.hcl file.
