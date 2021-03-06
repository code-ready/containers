package machine

import (
	"github.com/code-ready/crc/pkg/crc/machine/state"
	"github.com/code-ready/crc/pkg/crc/machine/types"
	"github.com/pkg/errors"
)

// Return console URL if the VM is present.
func (client *client) GetConsoleURL() (*types.ConsoleResult, error) {
	// Here we are only checking if the VM exist and not the status of the VM.
	// We might need to improve and use crc status logic, only
	// return if the Openshift is running as part of status.
	libMachineAPIClient, cleanup := createLibMachineClient()
	defer cleanup()
	host, err := libMachineAPIClient.Load(client.name)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot load machine")
	}

	vmState, err := host.Driver.GetState()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting the state for host")
	}

	crcBundleMetadata, err := getBundleMetadataFromDriver(host.Driver)
	if err != nil {
		return nil, errors.Wrap(err, "Error loading bundle metadata")
	}

	clusterConfig, err := getClusterConfig(crcBundleMetadata)
	if err != nil {
		return nil, errors.Wrap(err, "Error loading cluster configuration")
	}

	return &types.ConsoleResult{
		ClusterConfig: *clusterConfig,
		State:         state.FromMachine(vmState),
	}, nil
}
