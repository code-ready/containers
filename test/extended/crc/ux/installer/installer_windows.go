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
	installerWindowTitle string        = "CodeReady Containers Setup"
	elementClickDelay    time.Duration = 2 * time.Second
	installationTime     time.Duration = 30 * time.Second
)

type installerElement struct {
	name   string
	id     string
	screen string
}

var (
	welcomeNextButton         = installerElement{name: "Next", id: "1", screen: "welcome"}
	licenseAcceptCheck        = installerElement{name: "accept", id: "1", screen: "license"}
	licenseNextButton         = installerElement{name: "Next", id: "3", screen: "license"}
	destinantionNextButton    = installerElement{name: "Next", id: "1", screen: "destination"}
	installationInstallButton = installerElement{name: "Install", id: "1", screen: "installation"}
	finalizationFinishButton  = installerElement{name: "Finish", id: "1", screen: "finalization"}
)

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
	if err := runInstaller(*g.InstallerPath); err != nil {
		return err
	}
	// Welcome screen
	if err := clickButton(welcomeNextButton); err != nil {
		return err
	}
	// License screen
	if err := clickButton(licenseAcceptCheck); err != nil {
		return err
	}
	if err := clickButton(licenseNextButton); err != nil {
		return err
	}
	// Destination
	if err := clickButton(destinantionNextButton); err != nil {
		return err
	}
	// Installation
	if err := clickButton(installationInstallButton); err != nil {
		return err
	}
	// wait installation process
	time.Sleep(installationTime)
	// Finalization
	if err := clickButton(finalizationFinishButton); err != nil {
		return err
	}
	// TODO which should be executed from a new cmd (to ensure ENVs)
	// Finalize context
	win32waf.Finalize()
	return nil
}

func runInstaller(installerPath string) error {
	cmd := exec.Command("msiexec.exe", fmt.Sprintf("/i %s /qf", installerPath))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting %v with error %v", cmd, err)
	}
	return nil
}

func clickButton(element installerElement) error {
	// Ensure the installer is the active window
	// Get action center window
	if installerWindow, err := ux.GetActiveElement(installerWindowTitle, ux.WINDOW); err != nil {
		return err
	} else {
		button, err := installerWindow.GetElement(element.name, ux.BUTTON)
		if err != nil {
			logging.Error(err)
			return err
		}
		if err := button.Click(); err != nil {
			return err
		}
	}
	return nil
}
