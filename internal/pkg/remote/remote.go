package remote

import (
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHComand struct {
	User       string
	IP         string
	PrivateKey string
}

func PublicKeyFile(file string) ssh.AuthMethod {
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

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
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

func Connect(sshCommand SSHComand) {
	sshConfig := &ssh.ClientConfig{
		User: sshCommand.User,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(sshCommand.PrivateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", sshCommand.IP+":22", sshConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	defer connection.Close()

	// runCommand("bash ~/a.sh", connection)
}
