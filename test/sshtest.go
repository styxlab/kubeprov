package main

//Note: /etc/ssh/sshd_config: MaxSessions 10

import (
	"fmt"
	"github.com/styxlab/kubeprov/pkg/ssh"
)

func main() {

	ipAddress := "1.2.3.4"

	auth := ssh.AuthKey("demo", "/home/demo/.ssh/id_ed25519")
	config := auth.Config("root")
	client := config.Client(ipAddress, "22")
	defer client.Close()

	output := client.RunCmd("uname -a; ls -l")
	fmt.Println(output)

	dir := "/home/demo/go/src/kubeprov/assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	output2 := client.RunCmd("ls")
	fmt.Println(output2)
}
