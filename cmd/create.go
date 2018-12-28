package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/styxlab/kubeprov/pkg/hetzner"
	"github.com/styxlab/kubeprov/pkg/ssh"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "creates a new kubernetes cluster",
	Run: CreateCluster,
}

func CreateCluster(cmd *cobra.Command, args []string) {
	
	hc := hetzner.Connect()

	serverSpec := hc.ServerSpec("cws@home", "demo", "cx11", "centos-7")

	//serverInst := serverSpec.Create() //.EnableRescue().PowerOn().WaitForRunning()
	/*ssh.ExecCmdLocal("hcloud", "server", "create", "--image", "centos-7", "--name", "demo", "--type", "cx11", "--ssh-key", "/home/cws/.ssh/id_ed25519.pub")
	ssh.ExecCmdLocal("hcloud", "server", "poweroff", "demo")
	ssh.ExecCmdLocal("hcloud", "server", "enable-rescue", "demo", "--ssh-key", "cws@home")
	ssh.ExecCmdLocal("hcloud", "server", "poweron", "demo")*/

	serverInst := serverSpec.Status()
	ipAddress := serverInst.IPv4()

	//fmt.Printf("NewIP = %s\n", serverInst.IPv4())

	//ipAddress := "116.203.36.158"

    //fmt.Printf("Created node '%s' with IP %s\n", serverInst.Name(), ipAddress)
    fmt.Printf("Server should be in rescue mode now: ssh -oStrictHostKeyChecking=no root@%s\n", ipAddress)

	time.Sleep(1 * time.Second)

	command := "uname -a"
	if err := ssh.ExecCmd("root", "22", ipAddress, command); err != nil {
		 fmt.Printf("Error executing remote command: %s\n", err)
	}

	auth := ssh.AuthKey("cws@home", "/home/cws/.ssh/id_ed25519")
	config := auth.Config("root")
	client := config.Client(ipAddress, "22")
	defer client.Close()

	output := client.RunCmd("uname -a")
	fmt.Println(output)

	/*dir := "/home/cws/go/src/kubeprov/assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	serverInst.Reboot()

	config2 := auth.Config("core")
	client2 := config2.Client(ipAddress, "22")
	defer client2.Close()

	output2 := client2.RunCmd("uname -a")
	fmt.Println(output2)

	fmt.Printf("CoreOs should be installed: ssh -oStrictHostKeyChecking=no core@%s\n", ipAddress)
	*/
}
