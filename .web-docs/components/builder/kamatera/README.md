Type: `kamatera`

The `kamatera` Packer builder builds [Kamatera cloud](https://www.kamatera.com/express/compute/)
images and adds them to the account's disk library.

The builder process is as follows:

1. Create a temporary Kamatera server according to the builder configuration
2. Establish SSH connection to the server
3. Run provisioners
4. Create a Kamatera hard disk image from the server
5. Poweroff and terminate the server

For details on how to use the Kamatera builder please see the [Tutorial](#tutorial) below.

## Configuration Reference

### Required Builder Configuration options:

- `api_secret` (string) - The Kamatera API secret to use. It can also be specified via environment variable `KAMATERA_API_SECRET`.
- `api_client_id` (string) - The Kamatera API client ID to use. It can also be specified via environment variable `KAMATERA_API_CLIENT_ID`.
- `datacenter` (string) - The Kamatera datacenter to which the server is deployed.
- `ssh_username` (string) - Must be set with a value of `root`

### Optional:

- `image_name` (string) - The name of the resulting image that will appear in your Kamatera hard disk library. Defaults to `packer-{{timestamp}}`.
- `cpu` (string) - The cpu assigned to the temporary server. Default to `1A`.
- `ram` (string) - The ram assigned to the temporary server. Default to `1024`.
- `image` (string) - The image used by the temporary server. Default to `ubuntu_server_18.04_64-bit`.
- `disk` (string) - The disk size in GB assigned to the temporary server. Default to `10`.
- `script` (string) - Script to run for provisioning the server - Bash for linux images, Powershell for windows images.
- `disable_ssh` (bool) - Disable SSH access to the server, this is required when provisioning a Windows image.

## Tutorial

This tutorial creates a Linux machine image, for Windows images see below.

Set environment variables with your Kamatera API key (you can generate an API key from the [Kamatera console](https://console.kamatera.com) -> API -> Keys):

```shell
export KAMATERA_API_CLIENT_ID=
export KAMATERA_API_SECRET=
```

Create a new directory named `packer_tutorial` and paste the following configuration into a file named `kamatera-ubuntu.pkr.hcl`:

```hcl
packer {
  required_plugins {
    kamatera = {
      version = ">= 0.5.0"
      source  = "github.com/kamatera/kamatera"
    }
  }
}

source "kamatera" "ubuntu" {
  datacenter = "EU"
  ssh_username = "root"
  image_name = "ubuntu-example-packer-image"
  script = <<-EOF
    echo Creating example.txt &&\
    echo "Hello World" > example.txt
  EOF
}

build {
  sources = [
    "source.kamatera.ubuntu"
  ]
}
```

Change to the `packer_tutorial` directory and run the following commands to download the Kamatera plugin:

```shell
packer init .
```

Build the image:

```shell
packer build .
```

Please wait, this will take a while...

When done, log-in to [Kamatera console](https://console.kamatera.com) and navigate to My Cloud -> Hard Disk Library.
Choose zone `EU` and click on MY PRIVATE IMAGES - you should see the created image there.

You can now create a new server from this image, SSH into that server and run `cat example.txt`.

## Creating Windows machine images

Windows machine images are supported but have to use the inline script argument instead of custom provisioners.

Following is an example of a windows image source configuration:

```
source "kamatera" "windows" {
  datacenter = "EU"
  image = "windows_server_2022_datacenter_64-bit"
  // windows images do not support SSH, this also prevents running any custom provisioners
  disable_ssh = true
  // ssh_user is still required even though it is not used
  ssh_username = "root"
  image_name = "windows-example-packer-image"
  // powershell script can be used to provision the server
  script = <<-EOF
    Write-Host "Hello World!"
    "Hello World" | Out-File -FilePath test.txt
  EOF
}
```

If you create a server using this configuration you should be able to find the created `test.txt` file on the server.
