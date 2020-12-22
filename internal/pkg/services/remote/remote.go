package remote

import (
	"fmt"
	"io"
	"log"
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

// RunCommands runs mutiple commands in ssh connection
// FIX: Need to be re looked at
func (r *Remote) RunCommands(cmds []string) {
	sess, err := r.Client.NewSession()
	if err != nil {
		logger.Fatal("Failed to create session: ", err)
	}
	defer sess.Close()

	// StdinPipe for commands
	stdin, err := sess.StdinPipe()
	if err != nil {
		logger.Fatal("Failed", err)
	}

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		logger.Fatal("Failed", err)
	}

	for _, cmd := range cmds {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			logger.Fatal("Failed", err)
		}
	}

	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

// RunCommand runs a single command over a ssh connection
func (r *Remote) RunCommand(cmd string) {
	logger.Info("Running `%s`", cmd)
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

// MakeNewConnection initiates new ssh connection to given host
func (r *Remote) MakeNewConnection() {
	connection, err := ssh.Dial("tcp", r.address, &r.Config)
	if err != nil {
		logger.Fatal("Failed to connect to %s. Error: %s", r.address, err.Error())
	}

	r.Client = connection
}

// SpawnShell spwans shell on remote machine
// Ctrl-C exits the shell
func (r *Remote) SpawnShell() error {
	session, _ := r.Client.NewSession()

	if err := setupPty(session); err != nil {
		logger.Error("Failed to set up pseudo terminal. Error: %s", err.Error())
		return err
	}

	c := make(chan os.Signal)
	// Ctrl-C exists the shell
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(session *ssh.Session) {
		logger.Info("Shell Spawned. Press %s to exit.", logger.Bold("Ctrl+C"))
		<-c
		logger.Info("Ctrl+C pressed. Exiting remote shell")
		session.Signal(ssh.SIGTERM)
	}(session)

	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		logger.Error("Failed to start interactive shell. Error: %s", err.Error())
		return err
	}

	return session.Wait()
}

// Close closes the ssh connection
func (r *Remote) Close() error {
	return r.Client.Close()
}
