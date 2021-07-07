//go:generate packer-sdc mapstructure-to-hcl2 -type Config,imageFilter

package kamatera

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	ApiUrl      string `mapstructure:"api_url"`
	ApiClientID string `mapstructure:"api_client_id"`
	ApiSecret   string `mapstructure:"api_secret"`

	PollInterval time.Duration `mapstructure:"poll_interval"`

	ServerName string `mapstructure:"server_name"`
	Datacenter string `mapstructure:"datacenter"`
	CPU        string `mapstructure:"cpu"`
	RAM        string `mapstructure:"ram"`
	Image      string `mapstructure:"image"`
	Password   string `mapstructure:"password"`
	Disk       string `mapstructure:"disk"`

	SnapshotName string `mapstructure:"snapshot_name"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {
	var md mapstructure.Metadata
	err := config.Decode(c, &config.DecodeOpts{
		Metadata:           &md,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"run_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, err
	}

	// Defaults
	if c.ApiUrl == "" {
		if os.Getenv("KAMATERA_API_URL") != "" {
			c.ApiUrl = os.Getenv("KAMATERA_API_URL")
		} else {
			c.ApiUrl = "https://cloudcli.cloudwm.com"
		}
	}
	if c.ApiSecret == "" {
		c.ApiSecret = os.Getenv("KAMATERA_API_SECRET")
	}
	if c.ApiClientID == "" {
		c.ApiClientID = os.Getenv("KAMATERA_API_CLIENT_ID")
	}

	if c.PollInterval == 0 {
		c.PollInterval = 500 * time.Millisecond
	}

	if c.SnapshotName == "" {
		def, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			panic(err)
		}
		// Default to packer-{{ unix timestamp (utc) }}
		c.SnapshotName = def
	}

	if c.ServerName == "" {
		// Default to packer-[time-ordered-uuid]
		c.ServerName = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
		if len(c.ServerName) > 40 {
			c.ServerName = c.ServerName[:40]
		}
	}

	var errs *packersdk.MultiError
	if es := c.Comm.Prepare(&c.ctx); len(es) > 0 {
		errs = packersdk.MultiErrorAppend(errs, es...)
	}
	if c.ApiUrl == "" {
		errs = packersdk.MultiErrorAppend(errs, errors.New("api url must be specified"))
	}
	if c.ApiSecret == "" {
		errs = packersdk.MultiErrorAppend(errs, errors.New("api secret must be specified"))
	}
	if c.ApiClientID == "" {
		errs = packersdk.MultiErrorAppend(errs, errors.New("api client id must be specified"))
	}

	if c.Datacenter == "" {
		c.Datacenter = defaultServerOption.Datacenter
	}
	if c.CPU == "" {
		c.CPU = defaultServerOption.CPU
	}
	if c.RAM == "" {
		c.RAM = defaultServerOption.RAM
	}
	if c.Image == "" {
		c.Image = defaultServerOption.Image
	}
	if c.Password == "" {
		c.Password = "__generate__"
	}
	if c.Disk == "" {
		c.Disk = defaultServerOption.Disk
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	packersdk.LogSecretFilter.Set(c.ApiSecret)

	return nil, nil
}

func getServerIP(state multistep.StateBag) (string, error) {
	return state.Get("server_ip").(string), nil
}
