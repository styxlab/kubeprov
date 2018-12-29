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

	flagFalse := false
    serverOpts := hcloud.ServerCreateOpts{
        Name: name,
        ServerType: &hcloud.ServerType{
            Name: stype,
        },
        Image: &hcloud.Image{
            Name: image,
        },
        StartAfterCreate: &flagFalse,
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

// Status retreives the current status
func (s *ServerSpec) Status() *ServerInstance {

	client := s.cc.client

	result, _, err := client.Server.GetByName(context.Background(), s.options.Name)
    if err != nil {
      	log.Fatal(err)
    }
    if result == nil {
    	log.Fatal("empty server result")
    }

 	return &ServerInstance {
    	spec: s,
    	server: result,
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

    s.server.RescueEnabled = true
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
    	//server.Status = result.Status
    	server = result
    	time.Sleep(2 * time.Second)
	} 

 	return s
}

// Reboot the cloud server
func (s *ServerInstance) Reboot() *ServerInstance {

	c := s.spec.cc
	//server := s.server

	server, _, err := c.client.Server.GetByName(context.Background(), s.server.Name)
		if err != nil {
		log.Fatal(err)
	}
	if server == nil {
		log.Fatal("server not found")
	}

	action, _, err := c.client.Server.Reboot(context.Background(), server)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.waitForAction(action); err != nil {
		log.Fatal("could not reboot server")
    }

    fmt.Printf("Server %d rebooted\n", server.ID)

 	return s
}

// WaitForRescueDisabled which is disabled after the reboot
func (s *ServerInstance) WaitForRescueDisabled() *ServerInstance {

	c := s.spec.cc
	server := s.server

	//TODO: Timeout

	for server.RescueEnabled != false {
		result, _, err := c.client.Server.GetByName(context.Background(), server.Name)
    	if err != nil {
        	log.Fatal(err)
    	}
    	if result == nil {
    		log.Fatal("empty server result")
    	}
    	server.RescueEnabled = result.RescueEnabled
    	time.Sleep(2 * time.Second)
	}

 	return s
}

// PowerOn starts the cloud server
func (s *ServerInstance) PowerOn() *ServerInstance {

	c := s.spec.cc
	server := s.server

	action, _, err := c.client.Server.Poweron(context.Background(), server)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.waitForAction(action); err != nil {
		log.Fatal("could not power on the server")
    }

    fmt.Printf("Server %d started\n", server.ID)

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

// Create an image
func (s *ServerInstance) CreateImage() *ServerInstance {

	c := s.spec.cc
	server := s.server

	opts := &hcloud.ServerCreateImageOpts{
		Type: "snapshot",
		Description: "CoreOS",
	}
	result, _, err :=  c.client.Server.CreateImage(context.Background(), server, opts)
	if err != nil {
		return err
	}

	if err := c.waitForAction(result.Action); err != nil {
		log.Fatal("could not create server image")
    }

    fmt.Printf("Server image created.\n")

	return s
}