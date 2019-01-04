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





	//https://gobyexample.com/worker-pools

	r1 := make(chan *hetzner.ServerInstance)
	go func() {
		core01 := startRescue("core01")
		fmt.Println(core01.Name())
		r1 <- core01
	}()

	r2 := make(chan *hetzner.ServerInstance)
	go func() {
		core02 := startRescue("core02")
		fmt.Println(core02.Name())
		r2 <- core02
	}()

	core01 := <- r1
	fmt.Println("received core01")

	c1 := make(chan string)
	go func() {
		installCoreOS(core01)
		startKubernetes(core01, core01, "master")
		c1 <- "done 01"
	}()

	core02 := <- r2
	fmt.Println("received core02")

	c2 := make(chan string)
	go func() {
		installCoreOS(core02)
		c2 <- "done 02"
	}()

   for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }

    fmt.Println("join node")
    startKubernetes(core02, core01, "worker")

	//core01.Delete()
	//core02.Delete()
}

func startRescue(name string) *hetzner.ServerInstance {

	hc := hetzner.Connect()
	imageSpec := hetzner.ImageByName("centos-7")
	serverSpec := hc.ServerSpec(name, "cx11", imageSpec)
	serverInst := serverSpec.Create().EnableRescue().PowerOn()

	return serverInst
}

func installCoreOS(s *hetzner.ServerInstance) {

	fmt.Println("Install CoreOS on", s.IPv4());

	client := openClient("root", s)
	defer client.Close()

	output := client.RunCmd("uname -a")
	fmt.Println(output)

	dir := "./assets/coreos/"
	client.UploadFile(dir+"ignition.json", "/root", false)
	client.UploadFile(dir+"install.sh", "/root", true)

	dir = "./assets/kubernetes/"
	client.UploadFile(dir+"kubeadm_preinst.sh", "/root", true)

	output = client.RunCmd("./install.sh")
	fmt.Println(output)

	output = client.RunCmd("./kubeadm_preinst.sh")
	fmt.Println(output)

	s.Reboot()
}

func startKubernetes(s *hetzner.ServerInstance, m *hetzner.ServerInstance, role string) {

	fmt.Println("Install Kubernetes on", s.IPv4());

	client := openClient("core", s)
	defer client.Close()

	dir := "./assets/kubernetes/"
	client.UploadFile(dir+"hostname.sh", "/home/core", true)

	output := client.RunCmd("sudo ./hostname.sh " + s.Name())
	fmt.Println(output)

	if role == "master" {
		client.UploadFile(dir+"kubeadm_master.sh", "/home/core", true)
		output = client.RunCmd("./kubeadm_master.sh")
		fmt.Println(output)
	}else {
		master := openClient("core", m)
		defer master.Close()

		cmd := "until $(ncat -z " + m.IPv4() + " 6443); do echo 'Waiting for API server to respond'; sleep 5; done"
		output = master.RunCmd(cmd)
		fmt.Println(output)
		joinCmd := master.RunCmd("sudo kubeadm token create --print-join-command")
		fmt.Println(joinCmd)
		client.RunCmd("sudo " + joinCmd)
	}
}

func openClient(name string, s *hetzner.ServerInstance) *ssh.Client {

	s.WaitForRunning()
	
	auth := ssh.AuthKey(s.PublicKeyName(), s.PrivateKeyFile())
	config := auth.Config(name)
	return config.Client(s.IPv4(), 22)
}
