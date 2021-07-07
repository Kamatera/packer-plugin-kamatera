# Hetzner Cloud Builder

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
- `disk` (string) - The disk assigned to the server. Default to `size=10`.
- `snapshot_name` (string) - The name of the resulting snapshot that will appear in your account as image description. Defaults to `packer-{{timestamp}}`. 


## Basic Example

Here is a basic example. It is completely valid as soon as you enter your own
access tokens:

```hcl
  api_client_id = "d66f959d724b45bddad8750b5fd5e728"
  api_secret = "f0396ba3b682e764594668aae6cfc524"

  ssh_username = "root"

  server_name = "packer_test"
  datacenter = "IL"
  cpu = "1A"
  ram = "1024"
  image = "ubuntu_server_18.04_64-bit"
  password = "__generate__"
```
