package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/styxlab/kubeprov/pkg/hetzner"
	"github.com/styxlab/kubeprov/pkg/ssh"
)

var (
	createCmd = &cobra.Command {
		Use:     "create",
		Short:   "creates a new kubernetes cluster on hetzner cloud",
		Run: CreateCluster,
	}
)

func CreateCluster(cmd *cobra.Command, args []string) {

	//master
	core01 :=  installCoreOS("core01")
	fmt.Println(core01.Name())

	joinCmd := ""
	joinCmd = installKubernetes(core01, "master", joinCmd)

	core02 :=  installCoreOS("core02")
	fmt.Println(core02.Name())

	result := installKubernetes(core02, "worker", joinCmd)
	fmt.Println(result)

	//core01.Delete()
	//core02.Delete()
}

func installCoreOS(name string) *hetzner.ServerInstance {

	hc := hetzner.Connect()
	imageSpec := hetzner.ImageByName("centos-7")
	serverSpec := hc.ServerSpec(name, "cx11", imageSpec)
	serverInst := serverSpec.Create().EnableRescue().PowerOn().WaitForRunning()

	installCoreOSonServer(serverInst)

	return serverInst.Reboot().WaitForRunning()
}

func installCoreOSonServer(s *hetzner.ServerInstance) {

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

func installKubernetes(s *hetzner.ServerInstance, role string, joinCmd string) string {

	ipAddress := s.IPv4()
	fmt.Println("Install Kubernetes on", ipAddress);

	auth := ssh.AuthKey(s.PublicKeyName(), s.PrivateKeyFile())
	config := auth.Config("core")
	client := config.Client(ipAddress, 22)
	defer client.Close()

	dir := "./assets/kubernetes/"
	client.UploadFile(dir+"kubeadm_install.sh", "/home/core", true)

	output := client.RunCmd("sudo ./kubeadm_install.sh " + s.Name())
	fmt.Println(output)

	if role == "master" {
		client.UploadFile(dir+"kubeadm_master.sh", "/home/core", true)
		output = client.RunCmd("./kubeadm_master.sh")
		fmt.Println(output)
		output = client.RunCmd("sudo kubeadm token create --print-join-command")
		fmt.Println(output)
		return output
	}else if 0 < len(joinCmd) {
		return client.RunCmd("sudo " + joinCmd)
	}
	return ""
}

/*
func createImageForCoreOS(hc *hetzner.Client) *hetzner.ImageSpec {

	imageSpec := hetzner.ImageByName("centos-7")
	serverSpec := hc.ServerSpec("coreos-install", "cx11", imageSpec)
	serverInst := serverSpec.Create().EnableRescue().PowerOn().WaitForRunning()

	installCoreOSonServer(serverInst)

	// Create the image before reboot in order to preserver ignition.json
	imageSpec = serverInst.CreateSnapshot("CoreOS")
	serverInst.Delete()

	return imageSpec
}
*/

/*
	//hc := hetzner.Connect()
	//imageSpec := createImageForCoreOS(hc)

	


	//master
	core01 := createServer("core01", imageSpec)
	fmt.Println(core01.Name())
	
	joinCmd := ""
	joinCmd = installKubernetes(core01, "master", joinCmd)
	
	//worker
	core02 := createServer("core02", imageSpec)
	fmt.Println(core02.Name())

	result := installKubernetes(core02, "worker", joinCmd)
	fmt.Println(result)

	hc.ImageDelete(imageSpec)
	*/
/*
	func createServer(name string, image *hetzner.ImageSpec) *hetzner.ServerInstance {
	//TODO: concurrent server starting
	hc := hetzner.Connect()
	serverSpec := hc.ServerSpec(name, "cx11", image)
	return serverSpec.Create().PowerOn().WaitForRunning()
}
*/