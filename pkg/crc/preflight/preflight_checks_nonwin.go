// +build !windows

package preflight

import (
	"fmt"
	"os"

	"github.com/code-ready/crc/pkg/crc/logging"
)

var nonWinPreflightChecks = [...]PreflightCheck{
	{
		configKeySuffix:  "check-root-user",
		checkDescription: "Checking if running as non-root",
		check:            checkIfRunningAsNormalUser,
		fix:              fixRunAsNormalUser,
	},
}

func checkIfRunningAsNormalUser() (bool, error) {
	if os.Geteuid() != 0 {
		return true, nil
	}
	logging.Debug("Ran as root")
	return false, fmt.Errorf("crc should be ran as a normal user")
}

func fixRunAsNormalUser() (bool, error) {
	return false, fmt.Errorf("crc should be ran as a normal user")
}
