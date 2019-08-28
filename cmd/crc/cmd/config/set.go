package config

import (
	"github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/errors"
	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.AddCommand(configSetCmd)
}

var configSetCmd = &cobra.Command{
	Use:   "set CONFIG-KEY VALUE",
	Short: "Sets a crc configuration property.",
	Long: `Sets a crc configuration property. Some of the configuration properties are equivalent
to the options that you set when you run the 'crc start' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			errors.ExitWithMessage(1, "Please provide a configuration property and value to set")
		}
		runConfigSet(args[0], args[1])
	},
}

func runConfigSet(key string, value interface{}) {
	_, ok := config.SettingsList[key]
	if !ok {
		errors.ExitWithMessage(1, "Config property does not exist: %s", key)
	}

	ok, expectedValue := runValidations(config.SettingsList[key].ValidationFns, value)
	if !ok {
		errors.ExitWithMessage(1, "Config value is invalid: %s, Expected: %s\n", value, expectedValue)
	}

	config.Set(key, value)
	if err := config.WriteConfig(); err != nil {
		errors.ExitWithMessage(1, "Error Writing config to file %s: %s", constants.ConfigPath, err.Error())
	}
}

func runValidations(validations []config.ValidationFnType, value interface{}) (bool, string) {
	for _, fn := range validations {
		ok, expectedValue := fn(value)
		if !ok {
			return false, expectedValue
		}
	}
	return true, ""
}
