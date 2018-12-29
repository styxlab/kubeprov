package cmd

import (
	"fmt"
	"time"
	"log"

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
	//serverInst := serverSpec.Status()
	fmt.Printf("Created node '%s' with IP %s\n", serverInst.Name(), serverInst.IPv4())

	installCoreOS(serverInst.IPv4())
	serverInst.Reboot()

	if !ssh.ScanPort(serverInst.IPv4(), 22, 2 * time.Second, 120 * time.Second) {
		log.Fatal("Portscan failed after timeout.")
	}

	auth := ssh.AuthKey("cws@home", "/home/cws/.ssh/id_ed25519")
	config := auth.Config("core")
	client := config.Client(serverInst.IPv4(), "22")
	defer client.Close()

	output := client.RunCmd("uname -a")
	fmt.Println(output)

	fmt.Printf("CoreOs installed: ssh -oStrictHostKeyChecking=no core@%s\n", serverInst.IPv4())

	serverInst.CreateImage()
}

func installCoreOS(ipAddress string) {

	//wait for open port
	if !ssh.ScanPort(ipAddress, 22, 2 * time.Second, 120 * time.Second){
		log.Fatal("Portscan failed after timeout.")
	}

	auth := ssh.AuthKey("cws@home", "/home/cws/.ssh/id_ed25519")
	config := auth.Config("root")
	client := config.Client(ipAddress, "22")
	defer client.Close()

	output := client.RunCmd("uname -a")
	fmt.Println(output)

	dir := "/home/cws/go/src/kubeprov/assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	output2 := client.RunCmd("./install.sh")
	fmt.Println(output2)
}
