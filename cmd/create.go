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

	imageSpec := createImageForCoreOS()

	//create new CoreOS servers
	createServer("core01", imageSpec)
	createServer("core02", imageSpec)
}

func createImageForCoreOS() *ImageSpec {

	hc := hetzner.Connect()
	imageSpec := hetzner.ImageByName("centos-7")
	serverSpec := hc.ServerSpec("cws@home", "coreOS-install", "cx11", imageSpec)
	serverInst := serverSpec.Create().EnableRescue().PowerOn().WaitForRunning()
	//serverInst := serverSpec.Status()

	installCoreOS(serverInst.IPv4())
	
	// Create the image before reboot in order to preserver ignition.json
	imageSpec = serverInst.CreateSnapshot("CoreOS")

	// Delete server 
	serverInst.ServerDelete()

	return imageSpec
}

func installCoreOS(ipAddress string) {

	auth := ssh.AuthKey("cws@home", "/home/cws/.ssh/id_ed25519")
	config := auth.Config("root")
	client := config.Client(ipAddress, 22)
	defer client.Close()

	output := client.WaitForOpenPort().RunCmd("uname -a")
	fmt.Println(output)

	dir := "/home/cws/go/src/kubeprov/assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	output = client.RunCmd("./install.sh")
	fmt.Println(output)
}

func createServer(name string, image *ImageSpec){

	hc := hetzner.Connect()
	serverSpec := hc.ServerSpec("cws@home", name, "cx11", image)
	serverInst := serverSpec.Create().PowerOn().WaitForRunning()
}


	/* serverInst.Reboot()

	auth := ssh.AuthKey("cws@home", "/home/cws/.ssh/id_ed25519")
	config := auth.Config("core")
	client := config.Client(serverInst.IPv4(), "22")
	defer client.Close()

	output := client.RunCmd("uname -a")
	fmt.Println(output)

	fmt.Printf("CoreOs installed: ssh -oStrictHostKeyChecking=no core@%s\n", serverInst.IPv4())

	//now you can install a new server based on the new coreos image
	imageID: = imageInst.id
	serverSpec2 := hc.ServerSpec("cws@home", "core02", "cx11", imageID)
	serverInst2 := serverSpec2.Create().PowerOn().WaitForRunning()

	if err := ssh.WaitForOpenPort(serverInst.IPv4(), 22, 2 * time.Second, 60 * time.Second); err != nil {
		log.Fatal("Port closed after timeout.")
	}*/
