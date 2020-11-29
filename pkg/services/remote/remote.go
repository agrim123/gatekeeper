package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/agrim123/gatekeeper/pkg/logger"
	"golang.org/x/crypto/ssh"
)

type Remote struct {
	address string

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
	connection, err := ssh.Dial("tcp", r.address, &r.Config)
	if err != nil {
		logger.Fatalf("Failed to dial. Error: %s", err.Error())
	}

	r.Client = connection
}

func (r *Remote) SpawnShell() error {
	session, _ := r.Client.NewSession()

	if err := setupPty(session); err != nil {
		logger.Errorf("Failed to set up pseudo terminal. Error: %s", err.Error())
		return err
	}

	c := make(chan os.Signal)
	// Ctrl-C exists the shell
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(session *ssh.Session) {
		<-c
		logger.Info("Ctrl+C pressed. Exiting remote shell")
		session.Signal(ssh.SIGTERM)
	}(session)

	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		logger.Errorf("Failed to start interactive shell. Error: %s", err.Error())
		return err
	}
	return session.Wait()
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
		logger.Errorf("Request for pseudo terminal failed. Error: %s", err.Error())
		return err
	}
	return nil
}
