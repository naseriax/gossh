package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

//SshAgent object contains all ssh connectivity info and tools.
type SshAgent struct {
	Host     string
	Port     string
	UserName string
	Password string
	Timeout  int
	Client   *ssh.Client
	Session  *ssh.Session
}

//Connect connects to the specified server and opens a session (Filling the Client and Session fields in SshAgent struct)
func (s *SshAgent) Connect() error {
	var err error

	config := &ssh.ClientConfig{
		User: s.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(s.Timeout) * time.Second,
	}

	s.Client, err = ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Host, s.Port), config)
	if err != nil {
		log.Printf("Failed to dial: %v", err)
		return err
	}

	s.Session, err = s.Client.NewSession()
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		return err
	}
	return nil
}

//Exec executed a single command on the ssh session.
func (s *SshAgent) Exec(cmd string) (string, error) {
	var b bytes.Buffer
	s.Session.Stdout = &b
	if err := s.Session.Run(cmd); err != nil {
		return "", fmt.Errorf("failed to run: %v >> %v", cmd, err.Error())
	} else {
		return b.String(), nil
	}
}

//Disconnect closes the ssh sessoin and connection.
func (s *SshAgent) Disconnect() {
	s.Session.Close()
	s.Client.Close()
	log.Println("Closed the ssh session.")
}

//Init initialises the ssh connection and returns the usable ssh agent.
func Init(host, port, username, password string, timeout int) (SshAgent, error) {
	sshagent := SshAgent{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
		Timeout:  timeout,
	}
	err := sshagent.Connect()
	if err != nil {
		return sshagent, err
	} else {
		return sshagent, nil
	}
}
