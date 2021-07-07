# packer-plugin-kamatera

The `Kamatera` multi-component plugin can be used with HashiCorp [Packer](https://www.packer.io)
to create custom images. For the full list of available features for this plugin see [docs](docs).

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    hcloud = {
      version = ">= 0.1.0"
      source  = "github.com/kamatera/kamatera"
    }
  }
}
```


#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/Kamatera/packer-plugin-kamatera/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### From Sources

If you prefer to build the plugin from sources, clone the GitHub repository
locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-kamatera` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


Type: `kamatera`
Artifact BuilderId: `kamatera.builder`

## Configuration Reference

There are many configuration options available for the builder. They are
segmented below into two categories: required and optional parameters. Within
each category, the available configuration keys are alphabetized.

### Required Builder Configuration options:

- `api_url` (string) - The server endpoint to use. It can also be specified via environment variable `KAMATERA_API_URL`, if set.
- `api_secret` (string) - The API secret to use. It can also be specified via environment variable `KAMATERA_API_SECRET`, if set.
- `api_client_id` (string) - The API client ID to use. It can also be specified via environment variable `KAMATERA_API_CLIENT_ID`, if set.

### Optional:

- `server_name` (string) - The name assigned to the server.
- `datacenter` (string) - The datacenter to which the server is deployed. Default to `IL`.
- `cpu` (string) - The cpu assigned to the server. Default to `1A`.
- `ram` (string) - The ram assigned to the server. Default to `1024`.
- `image` (string) - The image used by the server. Default to `ubuntu_server_18.04_64-bit`.
- `password` (string) - The password assigned to the server. Default to `__generate__`.
- `disk` (string) - The disk assigned to the server. Default to `10`.
- `snapshot_name` (string) - The name of the resulting snapshot that will appear in your account as image description. Defaults to `packer-{{timestamp}}`.


## Basic Example

Here is a basic example. It is completely valid as soon as you enter your own
access tokens:

```hcl
  api_client_id = ""
  api_secret = ""

  ssh_username = "root"

  server_name = "packer_test"
  datacenter = "IL"
  cpu = "1A"
  ram = "1024"
  image = "ubuntu_server_18.04_64-bit"
  password = "__generate__"
```

Please check `example` directory for the full example.

## Contributing

* If you think you've found a bug in the code, or you have a question regarding
  the usage of this software, please reach out to us by opening an issue in
  this GitHub repository.
* Contributions to this project are welcome: if you want to add a feature or a
  fix a bug, please do so by opening a Pull Request in this GitHub repository.
  In case of feature contribution, we kindly ask you to open an issue to
  discuss it beforehand.
  
