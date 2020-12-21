package remote

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/agrim123/gatekeeper/pkg/logger"
	"golang.org/x/crypto/ssh"
)

func NewRemoteConnection(user, ip, port, privateKey string) *Remote {
	pubkey, err := getPubKeyAuthMethod(privateKey)
	if err != nil {
		panic(err)
	}

	remote := Remote{
		Config: ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				pubkey,
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: change
			Timeout:         5 * time.Second,
		},
		address: ip + ":" + port,
	}

	return &remote
}

func verifyPrivateKeyPermissions(privateKey string) error {
	info, err := os.Stat(privateKey)
	if err != nil {
		return err
	}

	allowedPerm := uint32(0400)
	if uint32(info.Mode()) & ^allowedPerm == 0 {
		return nil
	}

	return fmt.Errorf("Check private key: '%s' permissions. Have %v, want %v", privateKey, info.Mode(), os.FileMode(allowedPerm))
}

func getPubKeyAuthMethod(file string) (ssh.AuthMethod, error) {
	logger.InfofL("Reading private key")
	if err := verifyPrivateKeyPermissions(file); err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(key), nil
}

// pty = pseudo terminal
func setupPty(session *ssh.Session) error {
	modes := ssh.TerminalModes{
		ssh.ECHO: 0, // disable echoing
		// ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		// ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		logger.Error("Request for pseudo terminal failed. Error: %s", err.Error())
		return err
	}
	return nil
}
