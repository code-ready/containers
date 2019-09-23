package preflight

import (
	cmdConfig "github.com/code-ready/crc/cmd/crc/cmd/config"
	"github.com/code-ready/crc/pkg/crc/config"
)

// StartPreflightChecks performs the preflight checks before starting the cluster
func StartPreflightChecks(vmDriver string) {
	preflightCheckSucceedsOrFails(false,
		checkIfRunningAsNormalUser,
		"Checking if running as non-root",
		false,
	)
	preflightCheckSucceedsOrFails(false,
		checkOcBinaryCached,
		"Checking if oc binary is cached",
		false,
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckVirtEnabled.Name),
		checkVirtualizationEnabled,
		"Checking if Virtualization is enabled",
		config.GetBool(cmdConfig.WarnCheckVirtEnabled.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckKvmEnabled.Name),
		checkKvmEnabled,
		"Checking if KVM is enabled",
		config.GetBool(cmdConfig.WarnCheckKvmEnabled.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckLibvirtInstalled.Name),
		checkLibvirtInstalled,
		"Checking if libvirt is installed",
		config.GetBool(cmdConfig.WarnCheckLibvirtInstalled.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckUserInLibvirtGroup.Name),
		checkUserPartOfLibvirtGroup,
		"Checking if user is part of libvirt group",
		config.GetBool(cmdConfig.WarnCheckUserInLibvirtGroup.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckLibvirtEnabled.Name),
		checkLibvirtEnabled,
		"Checking if libvirt is enabled",
		config.GetBool(cmdConfig.WarnCheckLibvirtEnabled.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckLibvirtRunning.Name),
		checkLibvirtServiceRunning,
		"Checking if libvirt daemon is running",
		config.GetBool(cmdConfig.WarnCheckLibvirtRunning.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckLibvirtVersionCheck.Name),
		checkLibvirtVersion,
		"Checking if a supported libvirt version is installed",
		config.GetBool(cmdConfig.WarnCheckLibvirtVersionCheck.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckLibvirtDriver.Name),
		checkMachineDriverLibvirtInstalled,
		"Checking if crc-driver-libvirt is installed",
		config.GetBool(cmdConfig.WarnCheckLibvirtDriver.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckCrcNetwork.Name),
		checkLibvirtCrcNetworkAvailable,
		"Checking if libvirt 'crc' network is available",
		config.GetBool(cmdConfig.WarnCheckCrcNetwork.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckCrcNetworkActive.Name),
		checkLibvirtCrcNetworkActive,
		"Checking if libvirt 'crc' network is active",
		config.GetBool(cmdConfig.WarnCheckCrcNetworkActive.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckNetworkManagerInstalled.Name),
		checkNetworkManagerInstalled,
		"Checking if NetworkManager is installed",
		config.GetBool(cmdConfig.WarnCheckNetworkManagerInstalled.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckNetworkManagerRunning.Name),
		CheckNetworkManagerIsRunning,
		"Checking if NetworkManager service is running",
		config.GetBool(cmdConfig.WarnCheckNetworkManagerRunning.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckCrcNetworkManagerConfig.Name),
		checkCrcNetworkManagerConfig,
		"Checking if /etc/NetworkManager/conf.d/crc-nm-dnsmasq.conf exists",
		config.GetBool(cmdConfig.WarnCheckCrcDnsmasqFile.Name),
	)
	preflightCheckSucceedsOrFails(config.GetBool(cmdConfig.SkipCheckCrcDnsmasqFile.Name),
		checkCrcDnsmasqConfigFile,
		"Checking if /etc/NetworkManager/dnsmasq.d/crc.conf exists",
		config.GetBool(cmdConfig.WarnCheckCrcDnsmasqFile.Name),
	)
}

// SetupHost performs the prerequisite checks and setups the host to run the cluster
func SetupHost(vmDriver string) {
	preflightCheckAndFix(false,
		checkIfRunningAsNormalUser,
		fixRunAsNormalUser,
		"Checking if running as non-root",
		false,
	)
	preflightCheckAndFix(false,
		checkOcBinaryCached,
		fixOcBinaryCached,
		"Caching oc binary",
		false,
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckVirtEnabled.Name),
		checkVirtualizationEnabled,
		fixVirtualizationEnabled,
		"Setting up virtualization",
		config.GetBool(cmdConfig.WarnCheckVirtEnabled.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckKvmEnabled.Name),
		checkKvmEnabled,
		fixKvmEnabled,
		"Setting up KVM",
		config.GetBool(cmdConfig.WarnCheckKvmEnabled.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckLibvirtInstalled.Name),
		checkLibvirtInstalled,
		fixLibvirtInstalled,
		"Installing libvirt service and dependencies",
		config.GetBool(cmdConfig.WarnCheckLibvirtInstalled.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckUserInLibvirtGroup.Name),
		checkUserPartOfLibvirtGroup,
		fixUserPartOfLibvirtGroup,
		"Adding user to libvirt group",
		config.GetBool(cmdConfig.WarnCheckUserInLibvirtGroup.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckLibvirtEnabled.Name),
		checkLibvirtEnabled,
		fixLibvirtEnabled,
		"Enabling libvirt",
		config.GetBool(cmdConfig.WarnCheckLibvirtEnabled.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckLibvirtRunning.Name),
		checkLibvirtServiceRunning,
		fixLibvirtServiceRunning,
		"Starting libvirt service",
		config.GetBool(cmdConfig.WarnCheckLibvirtRunning.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckLibvirtDriver.Name),
		checkMachineDriverLibvirtInstalled,
		fixMachineDriverLibvirtInstalled,
		"Installing crc-driver-libvirt",
		config.GetBool(cmdConfig.WarnCheckLibvirtDriver.Name),
	)
	preflightCheckAndFix(false,
		checkOldMachineDriverLibvirtInstalled,
		fixOldMachineDriverLibvirtInstalled,
		"Removing older system-wide crc-driver-libvirt",
		false,
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckCrcNetwork.Name),
		checkLibvirtCrcNetworkAvailable,
		fixLibvirtCrcNetworkAvailable,
		"Setting up libvirt 'crc' network",
		config.GetBool(cmdConfig.WarnCheckCrcNetwork.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckCrcNetworkActive.Name),
		checkLibvirtCrcNetworkActive,
		fixLibvirtCrcNetworkActive,
		"Starting libvirt 'crc' network",
		config.GetBool(cmdConfig.WarnCheckCrcNetworkActive.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckNetworkManagerInstalled.Name),
		checkNetworkManagerInstalled,
		fixNetworkManagerInstalled,
		"Checking if NetworkManager is installed",
		config.GetBool(cmdConfig.WarnCheckNetworkManagerInstalled.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckNetworkManagerRunning.Name),
		CheckNetworkManagerIsRunning,
		fixNetworkManagerIsRunning,
		"Checking if NetworkManager service is running",
		config.GetBool(cmdConfig.WarnCheckNetworkManagerRunning.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckCrcNetworkManagerConfig.Name),
		checkCrcNetworkManagerConfig,
		fixCrcNetworkManagerConfig,
		"Writing Network Manager config for crc",
		config.GetBool(cmdConfig.WarnCheckCrcDnsmasqFile.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckCrcDnsmasqFile.Name),
		checkCrcDnsmasqConfigFile,
		fixCrcDnsmasqConfigFile,
		"Writing dnsmasq config for crc",
		config.GetBool(cmdConfig.WarnCheckCrcDnsmasqFile.Name),
	)
	preflightCheckAndFix(config.GetBool(cmdConfig.SkipCheckBundleCached.Name),
		checkBundleCached,
		fixBundleCached,
		"Unpacking bundle from the CRC binary",
		config.GetBool(cmdConfig.WarnCheckBundleCached.Name),
	)
}
