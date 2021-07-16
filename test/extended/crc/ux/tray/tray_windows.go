// +build windows

package tray

import (
	"fmt"

	win32waf "github.com/adrianriobo/gowinx/pkg/win32/api/user-interface/windows-accesibility-features"
	"github.com/adrianriobo/gowinx/pkg/win32/desktop/notificationarea"
	"github.com/adrianriobo/gowinx/pkg/win32/interaction"
	"github.com/adrianriobo/gowinx/pkg/win32/ux"
	clicumber "github.com/code-ready/clicumber/testsuite"
	// "github.com/adrianriobo/gowinx/pkg/win32/interaction"
	// "github.com/adrianriobo/gowinx/pkg/win32/ux"
)

type gowinxHandler struct {
	bundleLocation         *string
	pullSecretFileLocation *string
}

const (
	trayAssemblyName string = "crc-tray.exe"

	notificationIcon string = "CodeReady Containers"
	contextMenu      string = "menu"
	loginMenu        string = "copy-oc-login"
)

var (
	elements = map[string]string{
		actionStart:   "start",
		actionStop:    "stop",
		actionDelete:  "delete",
		actionQuit:    "exit",
		fieldState:    "status",
		userKubeadmin: "kubeadmin",
		userDeveloper: "developer"}
)

func NewTray(bundleLocationValue, pullSecretFileLocationValue *string) Tray {
	return gowinxHandler{
		bundleLocation:         bundleLocationValue,
		pullSecretFileLocation: pullSecretFileLocationValue}
}

func RequiredResourcesPath() (string, error) {
	return "", nil
}

func (g gowinxHandler) Install() error {
	return clicumber.ExecuteCommandSucceedsOrFails("crc setup", "succeeds")
}

func (g gowinxHandler) IsInstalled() error {
	command := fmt.Sprintf("tasklist /NH /FI \"IMAGENAME eq %s\"", trayAssemblyName)
	output := fmt.Sprintf("%s*", trayAssemblyName)
	return clicumber.CommandReturnShouldContain(command, output)
}

func (g gowinxHandler) IsAccessible() (err error) {
	initialize()
	err = notificationarea.ShowHiddenNotificationArea()
	if err != nil {
		_, _, err = notificationarea.GetIconPositionByTitle(notificationIcon)
	}
	finalize()
	return
}

func (g gowinxHandler) ClickStart() error {
	return click(actionStart)
}

func (g gowinxHandler) ClickStop() error {
	return click(actionStop)
}

func (g gowinxHandler) ClickDelete() error {
	return click(actionDelete)
}

func (g gowinxHandler) ClickQuit() error {
	return click(actionQuit)
}

func (g gowinxHandler) SetPullSecret() error {
	return fmt.Errorf("not implemented yet")
}

func (g gowinxHandler) IsClusterRunning() error {
	return waitTrayShowsFieldWithValue(stateRunning, checkTrayShowsStatusValue)
}

func (g gowinxHandler) IsClusterStopped() error {
	return waitTrayShowsFieldWithValue(stateStopped, checkTrayShowsStatusValue)
}

func (g gowinxHandler) CopyOCLoginCommandAsKubeadmin() error {
	return clickOnSubmenu(loginMenu, userKubeadmin)
}

func (g gowinxHandler) CopyOCLoginCommandAsDeveloper() error {
	return clickOnSubmenu(loginMenu, userDeveloper)
}

func (g gowinxHandler) ConnectClusterAsKubeadmin() error {
	return connectClusterAs(userKubeadmin)
}

func (g gowinxHandler) ConnectClusterAsDeveloper() error {
	return connectClusterAs(userDeveloper)
}

func initialize() {
	// Initialize context
	win32waf.Initalize()
}

func finalize() {
	// Finalize context
	win32waf.Finalize()
}

func click(action string) (err error) {
	// Initialize base elements
	initialize()
	_, err = clickOnContextMenu(action)
	// Finalize context
	finalize()
	return
}

func clickOnSubmenu(submenuElement, action string) (err error) {
	// Initialize base elements
	initialize()
	if submenu, clickErr := clickOnContextMenu(action); clickErr != nil {
		err = clickErr
	} else {
		_, err = clickAction(action, submenu)
	}
	// Finalize context
	finalize()
	return
}

func clickOnContextMenu(element string) (*ux.UXElement, error) {
	contextMenu, err := getContextMenu()
	if err != nil {
		return nil, err
	}
	return clickAction(element, contextMenu)
}

func clickAction(element string, menu *ux.UXElement) (*ux.UXElement, error) {
	buttonID, err := getElement(element, elements)
	if err != nil {
		return nil, err
	}
	button, err := menu.GetElement(buttonID, ux.MENUITEM)
	if err != nil {
		return nil, err
	}
	return button, button.Click()
}

func getContextMenu() (*ux.UXElement, error) {
	err := notificationarea.ShowHiddenNotificationArea()
	if err != nil {
		return nil, err
	}
	x, y, err := notificationarea.GetIconPositionByTitle(notificationIcon)
	if err != nil {
		return nil, err
	}
	if err := interaction.Click(int32(x), int32(y)); err != nil {
		return nil, err
	}
	return ux.GetActiveElement(contextMenu, ux.MENU)
}

func connectClusterAs(connectedUser string) error {
	//  Get oc
	err := clicumber.ExecuteCommand("crc oc-env | Invoke-Expression")
	if err != nil {
		return err
	}
	// Copy command from clipboard
	err = clicumber.ExecuteCommand("$clipboardCommand=Get-Clipboard")
	if err != nil {
		return err
	}
	// Run copied command
	err = clicumber.ExecuteCommand("Invoke-Expression -Command $clipboardCommand")
	if err != nil {
		return err
	}
	// Clear
	err = clicumber.ExecuteCommand("clear")
	if err != nil {
		return err
	}
	// Check user
	return clicumber.CommandReturnShouldContain("oc whoami", connectedUser)
}

func checkTrayShowsStatusValue(expectedValue string) (err error) {
	// Initialize base elements
	initialize()
	contextMenu, errContextMenu := getContextMenu()
	if errContextMenu != nil {
		err = errContextMenu
	} else {
		if statusLabel, errStatusLabel := contextMenu.GetElementByType(ux.TEXT); errStatusLabel != nil {
			err = errStatusLabel
		} else if statusLabel.GetName() != expectedValue {
			err = fmt.Errorf("Tray does not show value %s", expectedValue)
		}
	}
	// Finalize context
	finalize()
	return
}
