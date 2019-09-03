package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	globalViper *viper.Viper
	ViperConfig map[string]interface{}
)

// GetBool returns the value of a boolean config key
func GetBool(key string) bool {
	return globalViper.GetBool(key)
}

// Set sets the value for a give config key
func Set(key string, value interface{}) {
	globalViper.Set(key, value)
	ViperConfig[key] = value
}

func syncViperState(viper *viper.Viper) error {
	encodedConfig, err := json.MarshalIndent(ViperConfig, "", " ")
	if err != nil {
		return errors.Newf("Error encoding config to JSON: %v", err)
	}
	err = viper.ReadConfig(bytes.NewBuffer(encodedConfig))
	if err != nil {
		return errors.Newf("Error reading in new config: %s : %v", constants.ConfigFile, err)
	}
	return nil
}

// Unset unsets a given config key
func Unset(key string) error {
	delete(ViperConfig, key)
	return syncViperState(globalViper)
}

// GetString return the value of a key in string
func GetString(key string) string {
	return globalViper.GetString(key)
}

// GetInt return the value of a key in int
func GetInt(key string) int {
	return globalViper.GetInt(key)
}

// EnsureConfigFileExists creates the viper config file if it does not exists
func EnsureConfigFileExists() error {
	_, err := os.Stat(constants.ConfigPath)
	if err != nil {
		f, err := os.Create(constants.ConfigPath)
		if err == nil {
			_, err = f.WriteString("{}")
			f.Close()
		}
		return err
	}
	return nil
}

// InitViper initializes viper
func InitViper() error {
	v := viper.New()
	v.SetConfigFile(constants.ConfigPath)
	v.SetConfigType("json")
	v.SetEnvPrefix(constants.CrcEnvPrefix)
	// Replaces '-' in flags with '_' in env variables
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	v.SetTypeByDefaultValue(true)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Error Reading config file: %s : %v", constants.ConfigFile, err)
	}
	globalViper = v
	return v.Unmarshal(&ViperConfig)
}

// SetDefault sets the default for a config
func SetDefault(key string, value interface{}) {
	globalViper.SetDefault(key, value)
}

// WriteConfig write config to file
func WriteConfig() error {
	// We recreate a new viper instance, as globalViper.WriteConfig()
	// writes both default values and set values back to disk while we only
	// want the latter to be written
	v := viper.New()
	v.SetConfigFile(constants.ConfigPath)
	v.SetConfigType("json")
	err := syncViperState(v)
	if err != nil {
		return err
	}
	return v.WriteConfig()
}

// AllConfigs returns all the configs
func AllConfigs() map[string]interface{} {
	return globalViper.AllSettings()
}

// IsSet returns true if the config property is set
func IsSet(key string) bool {
	ss := AllConfigs()
	_, ok := ss[key]
	return ok
}

// BindFlags binds flags to config properties
func BindFlag(key string, flag *pflag.Flag) error {
	return globalViper.BindPFlag(key, flag)
}

// BindFlagset binds a flagset to their repective config properties
func BindFlagSet(flagSet *pflag.FlagSet) error {
	return globalViper.BindPFlags(flagSet)
}

type setting struct {
	Name          string
	DefaultValue  interface{}
	ValidationFns []ValidationFnType
}

// SettingsList holds all the config settings
var SettingsList = make(map[string]*setting)

// CreateSetting returns a filled struct of ConfigSetting
// takes the config name and default value as arguments
func AddSetting(name string, defValue interface{}, validationFn []ValidationFnType) *setting {
	s := setting{Name: name, DefaultValue: defValue, ValidationFns: validationFn}
	SettingsList[name] = &s
	return &s
}
