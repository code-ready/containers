// +build windows

package installer

import (
	"fmt"
	"os/exec"
	"time"

	win32waf "github.com/adrianriobo/gowinx/pkg/win32/api/user-interface/windows-accesibility-features"
	"github.com/adrianriobo/gowinx/pkg/win32/ux"
	"github.com/code-ready/crc/pkg/crc/logging"
)

const (
	installerWindowTitle string = "CodeReady Containers Setup"

	elementClickTime time.Duration = 2 * time.Second
	installationTime time.Duration = 30 * time.Second
)

var installFlow = map[string]button{
	"welcomeNext":         {"Next", elementClickTime},
	"licenseAccept":       {"accept", elementClickTime},
	"licenseNext":         {"Next", elementClickTime},
	"destinantionNext":    {"Next", elementClickTime},
	"installationInstall": {"Install", installationTime},
	"finalizationFinish":  {"Finish", elementClickTime}}

type button struct {
	id    string
	delay time.Duration
}

type gowinxHandler struct {
	CurrentUserPassword *string
	InstallerPath       *string
}

func NewInstaller(currentUserPassword, installerPath *string) Installer {
	// TODO check parameters as they are mandatory otherwise exit
	return gowinxHandler{
		CurrentUserPassword: currentUserPassword,
		InstallerPath:       installerPath}

}

func RequiredResourcesPath() (string, error) {
	return "", nil
}

func (g gowinxHandler) Install() error {
	// Initialize context
	win32waf.Initalize()
	installer, err := runInstaller(*g.InstallerPath)
	if err != nil {
		return err
	}
	for _, action := range installFlow {
		actionButton, err := installer.GetElement(action.id, ux.BUTTON)
		if err != nil {
			logging.Error(err)
			return err
		}
		if err := actionButton.Click(); err != nil {
			logging.Error(err)
			return err
		}
		// Delay after action
		time.Sleep(action.delay)
	}
	// TODO which should be executed from a new cmd (to ensure ENVs)
	// Finalize context
	win32waf.Finalize()
	return nil
}

func runInstaller(installerPath string) (*ux.UXElement, error) {
	cmd := exec.Command("msiexec.exe", "/i", installerPath, "/qf")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting %v with error %v", cmd, err)
	}
	// delay to get window as active
	time.Sleep(1 * time.Second)
	return ux.GetActiveElement(installerWindowTitle, ux.WINDOW)
}
