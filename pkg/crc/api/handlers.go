package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/code-ready/crc/cmd/crc/cmd/config"
	crcConfig "github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/machine"
	"github.com/code-ready/crc/pkg/crc/validation"
	"github.com/code-ready/crc/pkg/crc/version"
)

func statusHandler(_ ArgsType) string {
	statusConfig := machine.ClusterStatusConfig{Name: constants.DefaultName}
	clusterStatus, _ := machine.Status(statusConfig)
	return encodeStructToJSON(clusterStatus)
}

func stopHandler(_ ArgsType) string {
	stopConfig := machine.StopConfig{
		Name:  constants.DefaultName,
		Debug: true,
	}
	commandResult, _ := machine.Stop(stopConfig)
	return encodeStructToJSON(commandResult)
}

func startHandler(_ ArgsType) string {
	startConfig := machine.StartConfig{
		Name:          constants.DefaultName,
		BundlePath:    crcConfig.GetString(config.Bundle.Name),
		Memory:        crcConfig.GetInt(config.Memory.Name),
		CPUs:          crcConfig.GetInt(config.CPUs.Name),
		NameServer:    crcConfig.GetString(config.NameServer.Name),
		GetPullSecret: getPullSecretFileContent,
		Debug:         true,
	}
	status, _ := machine.Start(startConfig)
	return encodeStructToJSON(status)
}

func versionHandler(_ ArgsType) string {
	v := &machine.VersionResult{
		CrcVersion:       version.GetCRCVersion(),
		CommitSha:        version.GetCommitSha(),
		OpenshiftVersion: version.GetBundleVersion(),
		Success:          true,
	}
	return encodeStructToJSON(v)
}

func getPullSecretFileContent() (string, error) {
	data, err := ioutil.ReadFile(crcConfig.GetString(config.PullSecretFile.Name))
	if err != nil {
		return "", err
	}
	pullsecret := string(data)
	if err := validation.ImagePullSecret(pullsecret); err != nil {
		return "", err
	}
	return pullsecret, nil
}

func deleteHandler(_ ArgsType) string {
	delConfig := machine.DeleteConfig{Name: constants.DefaultName}
	r, _ := machine.Delete(delConfig)
	return encodeStructToJSON(r)
}

func webconsoleURLHandler(_ ArgsType) string {
	consoleConfig := machine.ConsoleConfig{Name: constants.DefaultName}
	r, _ := machine.GetConsoleURL(consoleConfig)
	return encodeStructToJSON(r)
}

func encodeStructToJSON(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		logging.Error(err.Error())
		err := commandError{
			Err: "Failed while encoding JSON to string",
		}
		s, _ := json.Marshal(err)
		return string(s)
	}
	return string(s)
}

func encodeErrorToJSON(errMsg string) string {
	err := commandError{
		Err: errMsg,
	}
	return encodeStructToJSON(err)
}
