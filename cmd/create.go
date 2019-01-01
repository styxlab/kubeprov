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
	Short:   "creates a new kubernetes cluster on hetzner cloud",
	Run: CreateCluster,
}

func CreateCluster(cmd *cobra.Command, args []string) {

	hc := hetzner.Connect()
	imageSpec := createImageForCoreOS(hc)

	//create new CoreOS servers
	core01 := createServer("core01", imageSpec)
	core02 := createServer("core02", imageSpec)

	fmt.Println(core01.Name())
	fmt.Println(core02.Name())

	installKubernetes(core01)
	//installKubernetes(core02)

	hc.ImageDelete(imageSpec)
	//core01.Delete()
	//core02.Delete()
}

func createImageForCoreOS(hc *hetzner.Client) *hetzner.ImageSpec {

	imageSpec := hetzner.ImageByName("centos-7")
	serverSpec := hc.ServerSpec("coreos-install", "cx11", imageSpec)
	serverInst := serverSpec.Create().EnableRescue().PowerOn().WaitForRunning()

	installCoreOS(serverInst)

	// Create the image before reboot in order to preserver ignition.json
	imageSpec = serverInst.CreateSnapshot("CoreOS")
	serverInst.Delete()

	return imageSpec
}

func installCoreOS(s *hetzner.ServerInstance) {

	ipAddress := s.IPv4()
	fmt.Println("Install CoreOS on", ipAddress);

	auth := ssh.AuthKey(s.PublicKeyName(), s.PrivateKeyFile())
	config := auth.Config("root")
	client := config.Client(ipAddress, 22)
	defer client.Close()

	output := client.RunCmd("uname -a")
	fmt.Println(output)

	dir := "./assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	output = client.RunCmd("./install.sh")
	fmt.Println(output)
}

func createServer(name string, image *hetzner.ImageSpec) *hetzner.ServerInstance {
	//TODO: concurrent server starting
	hc := hetzner.Connect()
	serverSpec := hc.ServerSpec(name, "cx11", image)
	return serverSpec.Create().PowerOn().WaitForRunning()
}

func installKubernetes(s *hetzner.ServerInstance){

	ipAddress := s.IPv4()
	fmt.Println("Install Kubernetes on", ipAddress);

	auth := ssh.AuthKey(s.PublicKeyName(), s.PrivateKeyFile())
	config := auth.Config("core")
	client := config.Client(ipAddress, 22)
	defer client.Close()

	dir := "./assets/kubernetes/"
	client.UploadFile(dir+"kubernetes_install.sh", "/home/core", true)

	output := client.RunCmd("chmod +x ./kubernetes_install.sh; sudo ./kubernetes_install.sh " + s.Name())
	fmt.Println(output)

}