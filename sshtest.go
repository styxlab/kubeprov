package main

//Note: /etc/ssh/sshd_config: MaxSessions 10

import (
	
	"fmt"
	"kubeprov/pkg/ssh"
)

func main() {

    ipAddress := "116.203.46.235"

    auth := ssh.SSHKey("cws@home", "/home/cws/.ssh/id_ed25519") 
    config := auth.SSHConfig("root")
    client := config.SSHClient(ipAddress, "22")
    defer client.Close()
    
    output := client.RunCmd("uname -a; ls -l")
    fmt.Println(output)

    dir := "/home/cws/go/src/kubeprov/assets/coreos/"
    client.UploadFile(dir + "ignition.json", "/root", false);
    client.UploadFile(dir + "install.sh", "/root", true);

    output2 := client.RunCmd("ls")
    fmt.Println(output2)
}
