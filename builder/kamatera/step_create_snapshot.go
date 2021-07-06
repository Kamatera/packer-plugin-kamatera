package kamatera

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"packer-plugin-kamatera/httpclient"
)

type stepCreateSnapshot struct {
}

func (s *stepCreateSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	kamateraClient := state.Get("kamateraClient").(*httpclient.Kamatera)

	ui := state.Get("ui").(packersdk.Ui)
	c := state.Get("config").(*Config)

	ui.Say("Creating snapshot ...")

	res, err := kamateraClient.Request("POST", "server/hdlib", struct {
		Name      string `json:"name"`
		Clone     string `json:"clone"`
		ImageName string `json:"image_name"`
	}{
		state.Get("server_name").(string),
		state.Get("hduuid").(string),
		c.SnapshotName,
	})
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

	res, err = kamateraClient.Request("POST", "server/hdlib", struct {
		Name      string `json:"name"`
		Clone     string `json:"clone"`
		ImageName string `json:"image_name"`
	}{
		state.Get("server_name").(string),
		state.Get("hduuid").(string),
		c.SnapshotName,
	})
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	result := res.(map[string]interface{})
	snapshotUUID := result["uuid"].(string)

	ui.Say(fmt.Sprintf("snapshot: %s", snapshotUUID))
	state.Put("snapshot_uuid", snapshotUUID)

	return multistep.ActionContinue
}

func (s *stepCreateSnapshot) Cleanup(state multistep.StateBag) {

}
