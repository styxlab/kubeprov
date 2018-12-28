package hetzner

import (
    "context"
    "fmt"
    "time"
    "log"
    "os"
    "strings"

    "github.com/hetznercloud/hcloud-go/hcloud"
    "github.com/go-kit/kit/log/term"
    "github.com/gosuri/uiprogress"
)

// CloudClient holds the connection data for the Hetzner Cloud interaction
type CloudClient struct {
	client *hcloud.Client
	sshKey *hcloud.SSHKey
}

// ServerSpec specifies the properties of the cloud server
type ServerSpec struct {
	cc *CloudClient
	options hcloud.ServerCreateOpts
}

// ServerInstance represents cloud server taht was previously created
type ServerInstance struct {
	spec *ServerSpec
	server *hcloud.Server
}

// Connect opens a new Hetzner Cloud connection
func Connect() *CloudClient {
	return &CloudClient {
		client: hcloud.NewClient(hcloud.WithToken(strings.TrimSpace(os.Getenv("HCLOUD_TOKEN")))),
	}
}

// ServerSpec creates a cloud server specification
func (c *CloudClient) ServerSpec(keyid string, name string, stype string, image string) *ServerSpec {

	if c.sshKey == nil {
		c.getSSHKey(keyid)
	}

    serverOpts := hcloud.ServerCreateOpts{
        Name: name,
        ServerType: &hcloud.ServerType{
            Name: stype,
        },
        Image: &hcloud.Image{
            Name: image,
        },
    }
    serverOpts.SSHKeys = append(serverOpts.SSHKeys, c.sshKey)

    return &ServerSpec {
    	cc: c,
    	options: serverOpts,
    }
}

func (c *CloudClient) getSSHKey(keyid string) {

	sshKey, _, err := c.client.SSHKey.Get(context.Background(), keyid)
	if err != nil {
		log.Fatal(err)
	}
    if sshKey == nil {
        log.Fatalf("we got some problem with the SSH-Key, chances are you are in the wrong context")
    }
    c.sshKey = sshKey
}

func (c *CloudClient) waitForAction(action *hcloud.Action) error {

	progress, errs := c.client.Action.WatchProgress(context.Background(), action)

	if term.IsTerminal(os.Stdout) {
		prog := uiprogress.New()

		prog.Start()
		bar := prog.AddBar(100).AppendCompleted().PrependElapsed()
		bar.Width = 40
		bar.Empty = ' '

		for {
			select {
			case err := <-errs:
				if err == nil {
					bar.Set(100)
				}
				prog.Stop()
				return err
			case p := <-progress:
				bar.Set(p)
			}
		}
	}
	return <-errs
}

// Create makes a new cloud server instance
func (s *ServerSpec) Create() *ServerInstance {

	client := s.cc.client

	fmt.Printf("creating server '%s'...", s.options.Name)
    result, _, err := client.Server.Create(context.Background(), s.options)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.cc.waitForAction(result.Action); err != nil {
		log.Fatal("could not create server")
    }
    return &ServerInstance {
    	spec: s,
    	server: result.Server,
    }
}

// EnableRescue activates the rescue mode
func (s *ServerInstance) EnableRescue() *ServerInstance {

	c := s.spec.cc

	rescueOpts := hcloud.ServerEnableRescueOpts{
        Type: "linux64",
    }
    rescueOpts.SSHKeys = append(rescueOpts.SSHKeys, c.sshKey)

    rescue, _, err := c.client.Server.EnableRescue(context.Background(), s.server, rescueOpts)
    if err != nil {
		log.Fatal(err)
	}
    if err := c.waitForAction(rescue.Action); err != nil {
		log.Fatal("could not enable rescue")
    }

    s.server.EnableRescue = true
    return s
}

// WaitForRunning waits for the server to attain the "running" status
func (s *ServerInstance) WaitForRunning() *ServerInstance {

	c := s.spec.cc
	server := s.server

	//TODO: Timeout

	for server.Status != "running" {
		result, _, err := c.client.Server.GetByName(context.Background(), server.Name)
    	if err != nil {
        	log.Fatal(err)
    	}
    	if result == nil {
    		log.Fatal("empty server result")
    	}
    	server.Status = result.Status
    	time.Sleep(2 * time.Second)
	} 

	fmt.Printf("Status = '%v'...\n", s.server.Status)
 	return s
}

// Reboot the cloud server
func (s *ServerInstance) Reboot() *ServerInstance {

	c := s.spec.cc
	server := s.server

	fmt.Printf("RescueEnabled = '%v'...\n", server.RescueEnabled)

	fmt.Printf("rebooting server '%s'...", server.Name)

	action, _, err := c.client.Server.Reboot(context.Background(), server)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.waitForAction(action); err != nil {
		log.Fatal("could not reboot server")
    }

    fmt.Printf("reboot completed?\n")

 	return s
}

// WaitForRescueDisabled which is disabled after the reboot
func (s *ServerInstance) WaitForRescueDisabled() *ServerInstance {

	c := s.spec.cc
	server := s.server

	//TODO: Timeout

	fmt.Printf("WaitForRescueDisabled for server '%s'...\n", server.Name)
	fmt.Printf("RescueEnabled = '%v'...\n", server.RescueEnabled)

	for server.RescueEnabled != false {
		result, _, err := c.client.Server.GetByName(context.Background(), server.Name)
    	if err != nil {
        	log.Fatal(err)
    	}
    	if result == nil {
    		log.Fatal("empty server result")
    	}
    	server.RescueEnabled = result.RescueEnabled
		fmt.Printf("RescueEnabled = '%v'...\n", server.RescueEnabled)

    	time.Sleep(2 * time.Second)
	}
	fmt.Printf("reboot completed?\n")



 	return s
}

// Name of server instance for convenience
func (s *ServerInstance) Name() string {
	return s.server.Name
}

// IPv4 address of server instance for convenience
func (s *ServerInstance) IPv4() string {
	return s.server.PublicNet.IPv4.IP.String()
}
