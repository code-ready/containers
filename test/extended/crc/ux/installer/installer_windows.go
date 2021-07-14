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

var installFlow = []element{
	{"Next", elementClickTime, ux.BUTTON},
	{"I accept the terms in the License Agreement", elementClickTime, ux.CHECKBOX},
	{"Next", elementClickTime, ux.BUTTON},
	{"Next", elementClickTime, ux.BUTTON},
	{"Install", installationTime, ux.BUTTON},
	{"Finish", elementClickTime, ux.BUTTON}}

type element struct {
	id          string
	delay       time.Duration
	elementType string
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
	time.Sleep(elementClickTime)
	for _, action := range installFlow {
		// delay to get window as active
		element, err := installer.GetElement(action.id, action.elementType)
		if err != nil {
			err = fmt.Errorf("error getting %s with error %v", action.id, err)
			logging.Error(err)
			return err
		}
		if err := element.Click(); err != nil {
			err = fmt.Errorf("error clicking %s with error %v", action.id, err)
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
