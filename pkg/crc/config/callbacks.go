package config

import (
	"fmt"
)

func RequiresRestartMsg(key string, _ interface{}) string {
	return fmt.Sprintf("Changes to configuration property '%s' are only applied when the CRC instance is started.\n"+
		"If you already have a running CRC instance, then for this configuration change to take effect, "+
		"stop the CRC instance with 'crc stop' and restart it with 'crc start'.", key)
}

func SuccessfullyApplied(key string, value interface{}) string {
	return fmt.Sprintf("Successfully configured %s to %s", key, value)
}
