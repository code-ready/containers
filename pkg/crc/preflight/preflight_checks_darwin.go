package preflight

import (
	"errors"
	"fmt"
	"github.com/code-ready/crc/pkg/crc/oc"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	crcos "github.com/code-ready/crc/pkg/os"
)

const (
	virtualBoxDownloadURL   = "https://download.virtualbox.org/virtualbox/6.0.4/VirtualBox-6.0.4-128413-OSX.dmg"
	virtualBoxMountLocation = "/Volumes/VirtualBox"

	resolverFile = "/etc/resolver/testing"
)

var (
	virtualBoxPkgLocation = fmt.Sprintf("%s/VirtualBox.pkg", virtualBoxMountLocation)
)

// Add darwin specific checks
func checkVirtualBoxInstalled() (bool, error) {
	path, err := exec.LookPath("VBoxManage")
	if err != nil {
		return false, errors.New("VirtualBox cli VBoxManage is not found in the path")
	}
	fi, _ := os.Stat(path)
	if fi.Mode()&os.ModeSymlink != 0 {
		path, err = os.Readlink(path)
		if err != nil {
			return false, errors.New("VirtualBox cli VBoxManage is not found in the path")
		}
	}
	return true, nil
}

func fixVirtualBoxInstallation() (bool, error) {
	// Download the driver binary in /tmp
	tempFilePath := filepath.Join(os.TempDir(), "virtualbox.dmg")
	out, err := os.Create(tempFilePath)
	if err != nil {
		return false, err
	}
	defer out.Close()
	resp, err := http.Get(virtualBoxDownloadURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	stdOut, stdErr, err := crcos.RunWithPrivilege("hdiutil", "attach", tempFilePath)
	if err != nil {
		return false, fmt.Errorf("Could not mount the virtualbox.dmg file: %s %v: %s", stdOut, err, stdErr)
	}
	stdOut, stdErr, err = crcos.RunWithPrivilege("installer", "-package", virtualBoxPkgLocation, "-target", "/")
	if err != nil {
		return false, fmt.Errorf("Could not install VirtualBox.pkg: %s %v: %s", stdOut, err, stdErr)
	}
	stdOut, stdErr, err = crcos.RunWithPrivilege("hdiutil", "detach", virtualBoxMountLocation)
	if err != nil {
		return false, fmt.Errorf("Could not install VirtualBox.pkg: %s %v: %s", stdOut, err, stdErr)
	}
	return true, nil
}

func checkResolverFilePermissions() (bool, error) {
	info, err := os.Stat(resolverFile)
	if err != nil {
		return false, fmt.Errorf("Unable to get permissions of the resolver file: %s", err)
	}

	m := info.Mode()

	// 16 is checking for user write permissions
	return m&16 != 1, nil
}

func fixResolverFilePermissions() (bool, error) {
	stdOut, stdErr, err := crcos.RunWithPrivilege("touch", resolverFile)
	if err != nil {
		return false, fmt.Errorf("Unable to create the resolver file: %s %v: %s", stdOut, err, stdErr)
	}

	stdOut, stdErr, err = crcos.RunWithPrivilege("chown", string(os.Getuid()), resolverFile)
	if err != nil {
		return false, fmt.Errorf("Unable to change permissions of the resolver file: %s %v: %s", stdOut, err, stdErr)
	}

	return true, nil
}

// Check if oc binary is cached or not
func checkOcBinaryCached() (bool, error) {
	oc := oc.OcCached{}
	if !oc.IsCached() {
		return false, errors.New("oc binary is not cached.")
	}
	return true, nil
}

func fixOcBinaryCached() (bool, error) {
	oc := oc.OcCached{}
	if err := oc.EnsureIsCached(); err != nil {
		return false, fmt.Errorf("Not able to download oc %v", err)
	}
	return true, nil
}
