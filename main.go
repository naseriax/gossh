package main

import "fmt"

func main() {
	host := "1.1.1.1"
	port := "22"
	uname := "root"
	pword := "pass"

	sshagent, err := Init(host, port, uname, pword, 10)
	if err != nil {
		panic(err)
	}

	defer sshagent.Disconnect()

	cmd := `uname -a`
	res, err := sshagent.Exec(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
