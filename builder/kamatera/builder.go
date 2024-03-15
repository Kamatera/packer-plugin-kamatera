package kamatera

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/Kamatera/packer-plugin-kamatera/httpclient"
)

const BuilderId = "kamatera.builder"

type Builder struct {
	config         Config
	runner         multistep.Runner
	kamateraClient *httpclient.Kamatera
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	warnings, errs := b.config.Prepare(raws...)
	if errs != nil {
		return nil, warnings, errs
	}

	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packersdk.Ui, hook packersdk.Hook) (packersdk.Artifact, error) {
	b.kamateraClient = httpclient.NewKamateraClient(b.config.ApiUrl, b.config.ApiClientID, b.config.ApiSecret, ui)
	// Set up the state
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("kamateraClient", b.kamateraClient)
	state.Put("hook", hook)
	state.Put("ui", ui)

	if b.config.DisableSsh {
		ui.Say("SSH has been disabled. Provisioning steps will not run.")
	}

	var steps []multistep.Step
	if !b.config.DisableSsh {
		steps = append(steps,
			&stepCreateSSHKey{
				Debug:        b.config.PackerDebug,
				DebugKeyPath: fmt.Sprintf("ssh_key_%s.pem", b.config.PackerBuildName),
			},
		)
	}
	steps = append(steps, &stepCreateServer{})
	if !b.config.DisableSsh {
		steps = append(steps,
			&communicator.StepConnect{
				Config:    &b.config.Comm,
				Host:      getServerIP,
				SSHConfig: b.config.Comm.SSHConfigFunc(),
			},
			&commonsteps.StepProvision{},
			&commonsteps.StepCleanupTempKeys{
				Comm: &b.config.Comm,
			},
		)
	}
	steps = append(steps,
		&stepGetHduuid{},
		&stepPowerOffServer{},
		&stepCreateImage{},
	)
	// Run the steps
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)
	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk("image_name"); !ok {
		return nil, nil
	}

	artifact := &Artifact{
		imageName: state.Get("image_name").(string),
		StateData: map[string]interface{}{"generated_data": state.Get("generated_data")},
	}

	return artifact, nil
}
