package kamatera

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/Kamatera/packer-plugin-kamatera/httpclient"
)

type stepCreateImage struct {
}

func (s *stepCreateImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	kamateraClient := state.Get("kamateraClient").(*httpclient.Kamatera)

	ui := state.Get("ui").(packersdk.Ui)
	c := state.Get("config").(*Config)

	ui.Say("Creating image ...")

	res, err := kamateraClient.Request("POST", "server/hdlib", struct {
		Name      string `json:"name"`
		Clone     string `json:"clone"`
		ImageName string `json:"image-name"`
	}{
		state.Get("server_name").(string),
		state.Get("hduuid").(string),
		c.ImageName,
	}, false)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	command, err := kamateraClient.WaitCommand(fmt.Sprintf("%.0f", res))
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if _, hasLog := command["log"]; !hasLog {
		err := errors.New("invalid response from Kamatera API: command is missing creation log")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("image: %s", c.ImageName))
	state.Put("image_name", c.ImageName)

	return multistep.ActionContinue
}

func (s *stepCreateImage) Cleanup(state multistep.StateBag) {

}
