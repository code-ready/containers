// +build integration

package crcsuite

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/code-ready/crc/pkg/download"
)

// Download bundle for testing
func DownloadBundle(bundleLocation string, bundleDestination string) (string, error) {

	if bundleLocation[:4] != "http" {

		// copy the file locall

		if bundleDestination == "." {
			bundleDestination, _ = os.Getwd()
		}
		fmt.Printf("Copying bundle from %s to %s.\n", bundleLocation, bundleDestination)
		bundleDestination = filepath.Join(bundleDestination, bundleName)

		source, err := os.Open(bundleLocation)
		if err != nil {
			return "", err
		}
		defer source.Close()

		destination, err := os.Create(bundleDestination)
		if err != nil {
			return "", err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			return "", err
		}

		err = destination.Sync()

		return bundleDestination, err
	}

	filename, err := download.Download(bundleLocation, bundleDestination, 0644)
	fmt.Printf("Downloading bundle from %s to %s.\n", bundleLocation, bundleDestination)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func CopyFilesToTestDir() {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error retrieving current dir: %s", err)
		os.Exit(1)
	}

	l := strings.Split(cwd, string(filepath.Separator))
	dataDirPieces := append([]string{string(filepath.Separator)}, l[:len(l)-3]...)
	dataDirPieces = append(dataDirPieces, "testdata")
	dataDir := filepath.Join(dataDirPieces...)

	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		fmt.Printf("Error occurred loading data files: %s", err)
		os.Exit(1)
	}

	destLoc, _ := os.Getwd()
	for _, file := range files {

		sFileName := filepath.Join(dataDir, file.Name())
		fmt.Printf("Copying %s to %s\n", sFileName, destLoc)

		sFile, err := os.Open(sFileName)
		if err != nil {
			fmt.Printf("Error occurred opening file: %s", err)
			os.Exit(1)
		}
		defer sFile.Close()

		dFileName := file.Name()
		dFile, err := os.Create(dFileName)
		if err != nil {
			fmt.Printf("Error occurred creating file: %s", err)
			os.Exit(1)
		}
		defer dFile.Close()

		_, err = io.Copy(dFile, sFile) // ignore num of bytes
		if err != nil {
			fmt.Printf("Error occurred copying file: %s", err)
			os.Exit(1)
		}

		err = dFile.Sync()
		if err != nil {
			fmt.Printf("Error occurred syncing file: %s", err)
			os.Exit(1)
		}
	}
}

func ParseFlags() {

	flag.StringVar(&bundleURL, "bundle-location", "embedded", "Path to the bundle to be used in tests")
	flag.StringVar(&pullSecretFile, "pull-secret-file", "", "Path to the file containing pull secret")
	flag.StringVar(&CRCBinary, "crc-binary", "", "Path to the CRC binary to be tested")
	flag.StringVar(&bundleVersion, "bundle-version", "", "Version of the bundle used in tests")
}
