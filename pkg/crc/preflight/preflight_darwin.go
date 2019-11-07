package preflight

import (
	"fmt"
	cmdConfig "github.com/code-ready/crc/cmd/crc/cmd/config"
)

var genericPreflightChecks = [...]PreflightCheck{
	{
		skipConfigName:   cmdConfig.SkipCheckRootUser.Name,
		warnConfigName:   cmdConfig.WarnCheckRootUser.Name,
		checkDescription: "Checking if running as non-root",
		check:            checkIfRunningAsNormalUser,
		fix:              fixRunAsNormalUser,
	},
	{
		checkDescription: "Checking if oc binary is cached",
		check:            checkOcBinaryCached,
		fixDescription:   "Caching oc binary",
		fix:              fixOcBinaryCached,
	},
	{
		skipConfigName:   cmdConfig.SkipCheckBundleCached.Name,
		warnConfigName:   cmdConfig.WarnCheckBundleCached.Name,
		checkDescription: "Unpacking bundle from the CRC binary",
		check:            checkBundleCached,
		fix:              fixBundleCached,
		flags:            SetupOnly,
	},
}

// SetupHost performs the prerequisite checks and setups the host to run the cluster
var hyperkitPreflightChecks = [...]PreflightCheck{
	{
		skipConfigName:   cmdConfig.SkipCheckHyperKitInstalled.Name,
		warnConfigName:   cmdConfig.WarnCheckHyperKitInstalled.Name,
		checkDescription: "Checking if HyperKit is installed",
		check:            checkHyperKitInstalled,
		fixDescription:   "Setting up virtualization with HyperKit",
		fix:              fixHyperKitInstallation,
	},
	{
		skipConfigName:   cmdConfig.SkipCheckHyperKitDriver.Name,
		warnConfigName:   cmdConfig.WarnCheckHyperKitDriver.Name,
		checkDescription: "Checking if crc-driver-hyperkit is installed",
		check:            checkMachineDriverHyperKitInstalled,
		fixDescription:   "Installing crc-machine-hyperkit",
		fix:              fixMachineDriverHyperKitInstalled,
	},
}

var dnsPreflightChecks = [...]PreflightCheck{
	{
		skipConfigName:   cmdConfig.SkipCheckResolverFilePermissions.Name,
		warnConfigName:   cmdConfig.WarnCheckResolverFilePermissions.Name,
		checkDescription: fmt.Sprintf("Checking file permissions for %s", resolverFile),
		check:            checkResolverFilePermissions,
		fixDescription:   fmt.Sprintf("Setting file permissions for %s", resolverFile),
		fix:              fixResolverFilePermissions,
	},
	{
		skipConfigName:   cmdConfig.SkipCheckHostsFilePermissions.Name,
		warnConfigName:   cmdConfig.WarnCheckHostsFilePermissions.Name,
		checkDescription: fmt.Sprintf("Checking file permissions for %s", hostFile),
		check:            checkHostsFilePermissions,
		fixDescription:   fmt.Sprintf("Setting file permissions for %s", hostFile),
		fix:              fixHostsFilePermissions,
	},
}

func getPreflightChecks() []PreflightCheck {
	checks := []PreflightCheck{}

	checks = append(checks, genericPreflightChecks[:]...)
	checks = append(checks, hyperkitPreflightChecks[:]...)
	checks = append(checks, dnsPreflightChecks[:]...)

	return checks
}

// StartPreflightChecks performs the preflight checks before starting the cluster
func StartPreflightChecks() {
	doPreflightChecks(getPreflightChecks())
}

// SetupHost performs the prerequisite checks and setups the host to run the cluster
func SetupHost() {
	doFixPreflightChecks(getPreflightChecks())
}
