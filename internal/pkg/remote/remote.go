package remote

import (
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
	remote := Remote{
		Config: ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				publicKeyFile(privateKey),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: change
		},
		Address: ip + ":" + port,
	}

	return &remote
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
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
