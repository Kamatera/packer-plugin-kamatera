package kamatera

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/uuid"

	"packer-plugin-kamatera/httpclient"
)

type createServerPostValues struct {
	Name             string `json:"name"`
	Password         string `json:"password"`
	PasswordValidate string `json:"passwordValidate"`
	SSHKey           string `json:"ssh-key"`
	Datacenter       string `json:"datacenter"`
	Image            string `json:"image"`
	CPU              string `json:"cpu"`
	RAM              string `json:"ram"`
	Disk             string `json:"disk"`
	DailyBackup      string `json:"dailybackup"`
	Managed          string `json:"managed"`
	Network          string `json:"network"`
	Quantity         string `json:"quantity"`
	BillingCycle     string `json:"billingcycle"`
	MonthlyPackage   string `json:"monthlypackage"`
	PowerOn          string `json:"poweronaftercreate"`
}

var defaultServerOption = struct {
	Datacenter     string `json:"datacenter"`
	Image          string `json:"image"`
	CPU            string `json:"cpu"`
	RAM            string `json:"ram"`
	Password       string `json:"password"`
	Disk           string `json:"disk"`
	DailyBackup    string `json:"dailybackup"`
	Managed        string `json:"managed"`
	Network        string `json:"network"`
	Quantity       string `json:"quantity"`
	BillingCycle   string `json:"billingcycle"`
	MonthlyPackage string `json:"monthlypackage"`
}{
	"IL",
	"ubuntu_server_18.04_64-bit",
	"1A",
	"1024",
	"__generate__",
	"size=10",
	"no",
	"no",
	"name=wan,ip=auto",
	"1",
	"hourly",
	"",
}

type stepCreateServer struct {
}

func (s *stepCreateServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	kamateraClient := state.Get("kamateraClient").(*httpclient.Kamatera)

	ui := state.Get("ui").(packersdk.Ui)
	c := state.Get("config").(*Config)
	pubKey := state.Get("public_key").(string)

	// Create the server based on configuration
	ui.Say("Creating server ...")

	values := createServerPostValues{
		Name:             fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())[:40],
		SSHKey:           pubKey,
		Datacenter:       c.Datacenter,
		Image:            c.Image,
		CPU:              c.CPU,
		RAM:              c.RAM,
		Password:         defaultServerOption.Password,
		PasswordValidate: defaultServerOption.Password,
		Disk:             c.Disk,
		DailyBackup:      defaultServerOption.DailyBackup,
		Managed:          defaultServerOption.Managed,
		Network:          defaultServerOption.Network,
		Quantity:         defaultServerOption.Quantity,
		BillingCycle:     defaultServerOption.BillingCycle,
		MonthlyPackage:   defaultServerOption.MonthlyPackage,
	}

	result, err := kamateraClient.Request("POST", "service/server", values)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("Waiting creation ...")

	var commandIDs []interface{}
	if r, ok := result.([]interface{}); ok {
		commandIDs = r
	} else {
		response := result.(map[string]interface{})
		commandIDs = response["commandIds"].([]interface{})
	}

	if len(commandIDs) != 1 {
		err := errors.New("invalid response from Kamatera API: did not return expected command ID")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	commandID := commandIDs[0].(string)
	command, err := kamateraClient.WaitCommand(commandID)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	createLog, hasCreateLog := command["log"]
	if !hasCreateLog {
		err := errors.New("invalid response from Kamatera API: command is missing creation log")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	createdServerName := ""
	for _, line := range strings.Split(createLog.(string), "\n") {
		if strings.HasPrefix(line, "Name: ") {
			createdServerName = strings.Replace(line, "Name: ", "", 1)
		}
	}
	if createdServerName == "" {
		err := errors.New("invalid response from Kamatera API: failed to get created server name")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("server_name", createdServerName)

	//result, err = kamateraClient.Request("POST", "service/server/ssh", struct {
	//	Name string `json:"name"`
	//}{
	//	createdServerName,
	//})
	//if err != nil {
	//	state.Put("error", err)
	//	ui.Error(err.Error())
	//	return multistep.ActionHalt
	//}
	//servers := result.([]interface{})
	//if len(servers) != 1 {
	//	err := fmt.Errorf("wrong number of server, got %v", len(servers))
	//	state.Put("error", err)
	//	ui.Error(err.Error())
	//	return multistep.ActionHalt
	//}
	//server := servers[0].(map[string]interface{})
	//state.Put("server_ip", server["externalIp"].(string))

	result, err = kamateraClient.Request("POST", "server/network", struct {
		Name string `json:"name"`
	}{
		createdServerName,
	})
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	servers := result.([]interface{})
	if len(servers) != 1 {
		err := fmt.Errorf("wrong number of server, got %v", len(servers))
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	server := servers[0].(map[string]interface{})
	state.Put("server_ip", server["ips"].([]interface{})[0].(string))

	return multistep.ActionContinue
}

func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	serverName := state.Get("server_name")
	if _, ok := serverName.(string); !ok {
		return
	}

	kamateraClient := state.Get("kamateraClient").(*httpclient.Kamatera)

	ui := state.Get("ui").(packersdk.Ui)

	// Destroy the server we just created
	ui.Say("Destroying server...")
	result, err := kamateraClient.Request("POST", "service/server/terminate", struct {
		Name  string `json:"name"`
		Force bool   `json:"force"`
	}{
		serverName.(string),
		true,
	})
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying server. Please destroy it manually: %s", err))
		return
	}

	commandIds := result.([]interface{})
	_, err = kamateraClient.WaitCommand(commandIds[0].(string))
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error destroying server. Please destroy it manually: %s", err))
	}
}
