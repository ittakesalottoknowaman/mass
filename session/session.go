package session

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

/*
3种情况
1.有密码 无key
2.有密码 有key
3.无密码 有key
*/

type sshSession struct {
	ip      string
	session *ssh.Session
}

func getAuthMethod(key, password string) ([]ssh.AuthMethod, error) {
	if key == "" {
		return []ssh.AuthMethod{ssh.Password(password)}, nil
	}

	b, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if password == "" {
		signer, err = ssh.ParsePrivateKey(b)
		if err != nil {
			return nil, err
		}
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(b, []byte(password))
		if err != nil {
			return nil, err
		}
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}

func (s *sshSession) Run(command string) error {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	s.session.Stdout = &stdoutBuf
	s.session.Stderr = &stderrBuf
	err := s.session.Run(command)
	fmt.Printf("out:%s err:%s", stdoutBuf.String(), stderrBuf.String())
	return err
}

func New(ip, port, username, password, key string, timeout int) (*sshSession, error) {
	authMethod, err := getAuthMethod(key, password)
	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: authMethod,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Duration(timeout) * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", ip, port), clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}
	return &sshSession{
		ip:      ip,
		session: session,
	}, nil
}
