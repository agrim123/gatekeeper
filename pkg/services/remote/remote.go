package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

type Remote struct {
	Address string

	Config ssh.ClientConfig

	Client *ssh.Client
}

func NewRemoteConnection(user, ip, port, privateKey string) *Remote {
	pubkey, err := publicKeyFile(privateKey)
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
		},
		Address: ip + ":" + port,
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

func publicKeyFile(file string) (ssh.AuthMethod, error) {
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

func (r *Remote) Close() error {
	return r.Client.Close()
}

func (r *Remote) RunCommand(cmd string) {
	sess, err := r.Client.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}

	go io.Copy(os.Stdout, sessStdOut)

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}

	go io.Copy(os.Stderr, sessStderr)

	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func (r *Remote) MakeNewConnection() {
	connection, err := ssh.Dial("tcp", r.Address, &r.Config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	r.Client = connection
}
