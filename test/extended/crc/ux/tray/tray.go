package tray

import (
	"fmt"
	"time"

	"github.com/code-ready/crc/test/extended/util"
)

type Tray interface {
	Install() error
	IsInstalled() error
	IsAccessible() error
	ClickStart() error
	ClickStop() error
	ClickDelete() error
	ClickQuit() error
	SetPullSecret() error
	IsClusterRunning() error
	IsClusterStopped() error
	CopyOCLoginCommandAsKubeadmin() error
	CopyOCLoginCommandAsDeveloper() error
	// TODO check if make sense create a new ux component
	ConnectClusterAsKubeadmin() error
	ConnectClusterAsDeveloper() error
}

func getElement(name string, elements map[string]string) (string, error) {
	identifier, ok := elements[name]
	if ok {
		return identifier, nil
	}
	return "", fmt.Errorf("element '%s', Can not be accessed from the tray", name)
}

func waitTrayShowsFieldWithValue(expectedValue string, fn func(string) error) error {
	retryCount := 15
	iterationDuration, extraDuration, err :=
		util.GetRetryParametersFromTimeoutInSeconds(retryCount, trayClusterStateTimeout)
	if err != nil {
		return err
	}
	for i := 0; i < retryCount; i++ {
		err := fn(expectedValue)
		if err == nil {
			return nil
		}
		time.Sleep(iterationDuration)
	}
	if extraDuration != 0 {
		time.Sleep(extraDuration)
	}
	return fmt.Errorf("Tray did not showed %s ", expectedValue)
}
