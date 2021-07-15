package cmd

import (
	"fmt"
	"time"

	clicumber "github.com/code-ready/clicumber/testsuite"
	"github.com/code-ready/crc/test/extended/util"
)

const (
	// timeout to wait for cluster to change its state
	clusterStateTimeout       = "900"
	CRCExecutableInstalled    = "installed"
	CRCExecutableNotInstalled = "notInstalled"
)

func SetConfigPropertyToValueSucceedsOrFails(property string, value string, expected string) error {
	cmd := "crc config set " + property + " " + value
	return clicumber.ExecuteCommandSucceedsOrFails(cmd, expected)
}

func UnsetConfigPropertySucceedsOrFails(property string, expected string) error {
	cmd := "crc config unset " + property
	return clicumber.ExecuteCommandSucceedsOrFails(cmd, expected)
}

func WaitForClusterInState(state string) error {
	retryCount := 15
	iterationDuration, extraDuration, err :=
		util.GetRetryParametersFromTimeoutInSeconds(retryCount, clusterStateTimeout)
	if err != nil {
		return err
	}
	for i := 0; i < retryCount; i++ {
		err := CheckCRCStatus(state)
		if err == nil {
			return nil
		}
		time.Sleep(iterationDuration)
	}
	if extraDuration != 0 {
		time.Sleep(extraDuration)
	}
	return fmt.Errorf("the did not reach the %s state", state)
}

func CheckCRCStatus(state string) error {
	expression := `.*OpenShift: .*Running \(v\d+\.\d+\.\d+.*\).*`
	if state == "stopped" {
		expression = ".*OpenShift: .*Stopped.*"
	}

	err := clicumber.ExecuteCommand("crc status")
	if err != nil {
		return err
	}
	return clicumber.CommandReturnShouldMatch("stdout", expression)
}

func CheckCRCExecutableState(state string) error {
	command := "which crc"
	// Create a new shell session to reload envs
	if err := clicumber.StartHostShellInstance(""); err != nil {
		return err
	}
	switch state {
	case CRCExecutableInstalled:
		return clicumber.ExecuteCommandSucceedsOrFails(command, "succeeds")
	case CRCExecutableNotInstalled:
		return clicumber.ExecuteCommandSucceedsOrFails(command, "fails")
	default:
		return fmt.Errorf("%s state is not defined as valid crc executable state", state)
	}
}

func CheckMachineNotExists() error {
	expression := `.*Machine does not exist.*`
	err := clicumber.ExecuteCommand("crc status")
	if err != nil {
		return err
	}
	return clicumber.CommandReturnShouldMatch("stderr", expression)
}

func DeleteCRC() error {

	command := "crc delete"
	_ = clicumber.ExecuteCommand(command)

	fmt.Printf("Deleted CRC instance (if one existed).\n")
	return nil
}
