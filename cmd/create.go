package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/styxlab/kubeprov/pkg/hetzner"
	//"github.com/styxlab/kubeprov/pkg/ssh"
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
    fmt.Printf("Server should be in rescue mode now: ssh root@%s\n", serverInst.IPv4())

}
