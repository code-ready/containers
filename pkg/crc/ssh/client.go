package ssh

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	"github.com/code-ready/machine/libmachine/log"
	"golang.org/x/crypto/ssh"
)

type Client interface {
	Output(command string) (string, error)
	Close()
}

type NativeClient struct {
	User     string
	Hostname string
	Port     int
	Auth     *Auth

	conn *ssh.Client
}

type Auth struct {
	Keys []string
}

func NewClient(user string, host string, port int, auth *Auth) (Client, error) {
	return &NativeClient{
		User:     user,
		Hostname: host,
		Port:     port,
		Auth:     auth,
	}, nil
}

func NewNativeConfig(user string, auth *Auth) (ssh.ClientConfig, error) {
	var (
		privateKeys []ssh.Signer
		authMethods []ssh.AuthMethod
	)

	for _, k := range auth.Keys {
		key, err := ioutil.ReadFile(k)
		if err != nil {
			log.Debugf("Cannot read private ssh key %s", k)
			continue
		}

		privateKey, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return ssh.ClientConfig{}, err
		}

		privateKeys = append(privateKeys, privateKey)
	}

	if len(privateKeys) > 0 {
		authMethods = append(authMethods, ssh.PublicKeys(privateKeys...))
	}

	return ssh.ClientConfig{
		User: user,
		Auth: authMethods,
		// #nosec G106
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Minute,
	}, nil
}

func (client *NativeClient) session() (*ssh.Session, error) {
	if client.conn == nil {
		var err error
		config, err := NewNativeConfig(client.User, client.Auth)
		if err != nil {
			return nil, fmt.Errorf("Error getting config for native Go SSH: %s", err)
		}
		client.conn, err = ssh.Dial("tcp", net.JoinHostPort(client.Hostname, strconv.Itoa(client.Port)), &config)
		if err != nil {
			return nil, err
		}
	}
	session, err := client.conn.NewSession()
	if err != nil {
		return nil, err
	}
	return session, err
}

func (client *NativeClient) Output(command string) (string, error) {
	session, err := client.session()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (client *NativeClient) Close() {
	if client.conn == nil {
		return
	}
	err := client.conn.Close()
	if err != nil {
		log.Debugf("Error closing SSH Client: %s", err)
	}
}
