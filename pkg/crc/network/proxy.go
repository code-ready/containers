package network

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/code-ready/crc/pkg/crc/logging"
)

var (
	globalProxy      ProxyConfig
	defaultNoProxies = []string{"127.0.0.1", "localhost"}
)

// ProxyConfig keeps the proxy configuration for the current environment
type ProxyConfig struct {
	HTTPProxy   string
	HTTPSProxy  string
	NoProxy     []string
	ProxyCAFile string
	ProxyCACert string
}

// SetProxyDefaults specifies proxy configuration. If an empty string is passed the corresponding environment variable
// is checked.
func SetProxyDefaults(httpProxy, httpsProxy, noProxy, proxyCAFile string) error {
	proxyCACert, err := getProxyCAData(proxyCAFile)
	if err != nil {
		return err
	}

	if httpProxy == "" {
		httpProxy = getProxyFromEnv("http_proxy")
	}
	if httpsProxy == "" {
		httpsProxy = getProxyFromEnv("https_proxy")
	}
	if noProxy == "" {
		noProxy = getProxyFromEnv("no_proxy")
	}

	if err := ValidateProxyURL(httpProxy); err != nil {
		return err
	}

	if err := ValidateProxyURL(httpsProxy); err != nil {
		return err
	}

	globalProxy = ProxyConfig{
		HTTPProxy:   httpProxy,
		HTTPSProxy:  httpsProxy,
		ProxyCAFile: proxyCAFile,
		ProxyCACert: proxyCACert,
	}
	if noProxy != "" {
		globalProxy.NoProxy = strings.Split(noProxy, ",")
	}
	return nil
}

// GetProxyConfig retrieve the proxy configuration previously set.
func GetProxyConfig() *ProxyConfig {
	return &ProxyConfig{
		HTTPProxy:   globalProxy.HTTPProxy,
		HTTPSProxy:  globalProxy.HTTPSProxy,
		ProxyCAFile: globalProxy.ProxyCAFile,
		ProxyCACert: globalProxy.ProxyCACert,
		NoProxy:     append(defaultNoProxies, globalProxy.NoProxy...),
	}
}

func getProxyCAData(proxyCAFile string) (string, error) {
	if proxyCAFile == "" {
		return "", nil
	}
	proxyCACert, err := ioutil.ReadFile(proxyCAFile)
	if err != nil {
		return "", err
	}
	// Before passing string back to caller function, remove the empty lines in the end
	return trimTrailingEOL(string(proxyCACert)), nil
}

func trimTrailingEOL(s string) string {
	return strings.TrimRight(s, "\n")
}

func getProxyFromEnv(proxyScheme string) string {
	p := os.Getenv(proxyScheme)
	if p == "" {
		p = os.Getenv(strings.ToUpper(proxyScheme))
	}
	return p
}

func setProxyEnv(key, value string) {
	setEnv(strings.ToLower(key), value)
	setEnv(strings.ToUpper(key), value)
}

func setEnv(key, value string) {
	before := os.Getenv(key)
	if before != value {
		logging.Warnf("Overriding environment variable %s: previously %s, now %s", key, before, value)
	}
	if err := os.Setenv(key, value); err != nil {
		logging.Errorf("Cannot set environment variable %s: %v", key, err)
	}
}

func (p *ProxyConfig) String() string {
	httpProxy, _ := hidePassword(p.HTTPProxy)
	httpsProxy, _ := hidePassword(p.HTTPSProxy)
	return fmt.Sprintf("HTTPProxy: %s, HTTPSProxy: %s, NoProxy: %s, ProxyCAFile: %s",
		httpProxy, httpsProxy, p.GetNoProxyString(), p.ProxyCAFile)
}

// AddNoProxy appends the specified host to the list of no proxied hosts.
func (p *ProxyConfig) AddNoProxy(host ...string) {
	p.NoProxy = append(p.NoProxy, host...)
}

func (p *ProxyConfig) GetNoProxyString() string {
	return strings.Join(p.NoProxy, ",")
}

// Sets the current config as environment variables in the current process.
func (p *ProxyConfig) ApplyToEnvironment() {
	if !p.IsEnabled() {
		return
	}

	if p.HTTPProxy != "" {
		setProxyEnv("http_proxy", p.HTTPProxy)
	}
	if p.HTTPSProxy != "" {
		setProxyEnv("https_proxy", p.HTTPSProxy)
	}
	if len(p.NoProxy) != 0 {
		setProxyEnv("no_proxy", p.GetNoProxyString())
	}
}

// Enabled returns true if at least one proxy (HTTP or HTTPS) is configured. Returns false otherwise.
func (p *ProxyConfig) IsEnabled() bool {
	return p.HTTPProxy != "" || p.HTTPSProxy != ""
}

// ValidateProxyURL validates that the specified proxyURL is valid
func ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil
	}

	if strings.HasPrefix(proxyURL, "https://") {
		return fmt.Errorf("Proxy URL '%s' is not valid: https is not supported", proxyURL)
	}
	if !strings.HasPrefix(proxyURL, "http://") {
		return fmt.Errorf("Proxy URL '%s' is not valid: url should start with http://", proxyURL)
	}
	if !govalidator.IsURL(proxyURL) {
		return fmt.Errorf("Proxy URL '%s' is not valid", proxyURL)
	}
	return nil
}
