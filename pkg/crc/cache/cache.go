package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/download"
	"github.com/code-ready/crc/pkg/embed"
	"github.com/code-ready/crc/pkg/extract"
	crcos "github.com/code-ready/crc/pkg/os"
	"github.com/pkg/errors"
)

type Cache struct {
	binaryName string
	archiveURL string
	destDir    string
	version    string
	getVersion func() (string, error)
}

type VersionMismatchError struct {
	ExpectedVersion string
	CurrentVersion  string
}

func (e *VersionMismatchError) Error() string {
	return fmt.Sprintf("expected: %s but got: %s", e.ExpectedVersion, e.CurrentVersion)
}

func New(binaryName string, archiveURL string, destDir string, version string, getVersion func() (string, error)) *Cache {
	return &Cache{binaryName: binaryName, archiveURL: archiveURL, destDir: destDir, version: version, getVersion: getVersion}
}

func NewOcCache(version string, getVersion func() (string, error)) *Cache {
	return New(constants.OcBinaryName, constants.GetOcURL(), constants.CrcOcBinDir, version, getVersion)
}

func NewPodmanCache(version string, getVersion func() (string, error)) *Cache {
	return New(constants.PodmanBinaryName, constants.GetPodmanURL(), constants.CrcBinDir, version, getVersion)
}

func NewGoodhostsCache(version string, getVersion func() (string, error)) *Cache {
	return New(constants.GoodhostsBinaryName, constants.GetGoodhostsURL(), constants.CrcBinDir, version, getVersion)
}

func (c *Cache) IsCached() bool {
	if _, err := os.Stat(filepath.Join(c.destDir, c.binaryName)); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Cache) EnsureIsCached() error {
	if !c.IsCached() || c.CheckVersion() != nil {
		return c.CacheBinary()
	}
	return nil
}

// CacheBinary downloads and caches the requested binary into the CRC directory
func (c *Cache) CacheBinary() error {
	if c.IsCached() && c.CheckVersion() == nil {
		return nil
	}

	// Create tmp dir to download the requested tarball
	tmpDir, err := ioutil.TempDir("", "crc")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	assetTmpFile, err := c.getBinary(tmpDir)
	if err != nil {
		return err
	}

	var extractedFiles []string
	// Check the file is tarball or not
	if isTarball(assetTmpFile) {
		// Extract the tarball and put it the cache directory.
		extractedFiles, err = extract.UncompressWithFilter(assetTmpFile, tmpDir, false,
			func(filename string) bool { return filepath.Base(filename) == c.binaryName })
		if err != nil {
			return errors.Wrapf(err, "Cannot uncompress '%s'", assetTmpFile)
		}
	} else {
		extractedFiles = append(extractedFiles, assetTmpFile)
		if filepath.Base(assetTmpFile) != c.binaryName {
			logging.Warnf("Binary name is %s but extracted file name is %s", c.binaryName, filepath.Base(assetTmpFile))
		}
	}

	// Copy the requested asset into its final destination
	err = os.MkdirAll(c.destDir, 0750)
	if err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "Cannot create the target directory.")
	}

	for _, extractedFilePath := range extractedFiles {
		finalBinaryPath := filepath.Join(c.destDir, filepath.Base(extractedFilePath))
		// If the file exists then remove it (ignore error) first before copy because with `0500` permission
		// it is not possible to overwrite the file.
		os.Remove(finalBinaryPath)
		err = crcos.CopyFileContents(extractedFilePath, finalBinaryPath, 0500)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) getBinary(destDir string) (string, error) {
	archiveName := filepath.Base(c.archiveURL)
	logging.Debugf("Trying to extract %s from crc binary", archiveName)

	destPath := filepath.Join(destDir, archiveName)
	if err := embed.Extract(archiveName, destPath); err != nil {
		if strings.HasPrefix(c.archiveURL, "file://") {
			return "", err
		}
		logging.Debugf("Trying to download %s", archiveName)
		return download.Download(c.archiveURL, destDir, 0600)
	}

	return destPath, nil
}

func (c *Cache) CheckVersion() error {
	// Check if version string is non-empty
	if c.version == "" {
		return nil
	}
	currentVersion, err := c.getVersion()
	if err != nil {
		return err
	}
	if currentVersion != c.version {
		return &VersionMismatchError{CurrentVersion: currentVersion, ExpectedVersion: c.version}
	}
	return nil
}

func isTarball(filename string) bool {
	tarballExtensions := []string{".tar", ".tar.gz", ".tar.xz", ".zip", ".tar.bz2", ".crcbundle"}
	for _, extension := range tarballExtensions {
		if strings.HasSuffix(strings.ToLower(filename), extension) {
			return true
		}
	}
	return false
}
