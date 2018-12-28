package cmd

import (
	"fmt"

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

	serverInst := serverSpec.Create().EnableRescue().PowerOn().WaitForRunning()

    fmt.Printf("Created node '%s' with IP %s\n", serverInst.Name(), serverInst.IPv4())
    fmt.Printf("Server should be in rescue mode now: ssh -oStrictHostKeyChecking=no root@%s\n", serverInst.IPv4())

   	auth := ssh.AuthKey("cws@home", "/home/cws/.ssh/id_ed25519")
	config := auth.Config("root")
	client := config.Client(serverInst.IPv4(), "22")
	defer client.Close()

	output := client.RunCmd("uname -a; ls -l")
	fmt.Println(output)

	dir := "/home/cws/go/src/kubeprov/assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	output2 := client.RunCmd("ls")
	fmt.Println(output2)

	serverInst.Reboot()
	fmt.Printf("CoreOs should be installed: ssh -oStrictHostKeyChecking=no core@%s\n", serverInst.IPv4())
}
