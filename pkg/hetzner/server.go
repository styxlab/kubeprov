package hetzner

import (
	"context"
    "fmt"
    "time"
    "log"

    "github.com/hetznercloud/hcloud-go/hcloud"
)

// ServerSpec specifies the properties of the cloud server
type ServerSpec struct {
	cc *Client
	options hcloud.ServerCreateOpts
}

// ServerInstance represents cloud server that was previously created
type ServerInstance struct {
	context context.Context
	spec *ServerSpec
	server *hcloud.Server
}

// ServerSpec creates a cloud server specification
func (c *Client) ServerSpec(name string, stype string, image *ImageSpec) *ServerSpec {

	flagFalse := false
    serverOpts := hcloud.ServerCreateOpts{
        Name: name,
        ServerType: &hcloud.ServerType{
            Name: stype,
        },
        Image: image.spec,
        StartAfterCreate: &flagFalse,
    }
    serverOpts.SSHKeys = append(serverOpts.SSHKeys, c.sshKey)

    return &ServerSpec {
    	cc: c,
    	options: serverOpts,
    }
}

// Name of server instance for convenience
func (s *ServerInstance) Name() string {
	return s.server.Name
}

// IPv4 address of server instance for convenience
func (s *ServerInstance) IPv4() string {
	return s.server.PublicNet.IPv4.IP.String()
}


// Create makes a new cloud server instance
func (s *ServerSpec) Create() *ServerInstance {

	client := s.cc.client

    result, _, err := client.Server.Create(s.cc.context, s.options)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.cc.waitForAction(result.Action); err != nil {
		log.Fatal("could not create server")
    }
    return &ServerInstance {
    	context: s.cc.context,
    	spec: s,
    	server: result.Server,
    }
}

func (s *ServerInstance) Delete() *ServerInstance {
	
	c := s.spec.cc
	server := s.server

	_, err := c.client.Server.Delete(s.context, server)
	if err != nil {
		log.Fatal(err)
	}

    fmt.Printf("Server %d deleted\n", server.ID)
    return s
}

// Status retreives the current status
func (s *ServerSpec) Status() *ServerInstance {

	client := s.cc.client

	result, _, err := client.Server.GetByName(s.cc.context, s.options.Name)
    if err != nil {
      	log.Fatal(err)
    }
    if result == nil {
    	log.Fatal("empty server result")
    }

 	return &ServerInstance {
 		context: s.cc.context,
    	spec: s,
    	server: result,
    }
}

// WaitForRunning waits for the server to attain the "running" status
func (s *ServerInstance) WaitForRunning() *ServerInstance {

	c := s.spec.cc
	server := s.server

	//TODO: Timeout

	for server.Status != "running" {
		result, _, err := c.client.Server.GetByName(s.context, server.Name)
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

// EnableRescue activates the rescue mode
func (s *ServerInstance) EnableRescue() *ServerInstance {

	c := s.spec.cc

	rescueOpts := hcloud.ServerEnableRescueOpts{
        Type: "linux64",
    }
    rescueOpts.SSHKeys = append(rescueOpts.SSHKeys, c.sshKey)

    rescue, _, err := c.client.Server.EnableRescue(s.context, s.server, rescueOpts)
    if err != nil {
		log.Fatal(err)
	}
    if err := c.waitForAction(rescue.Action); err != nil {
		log.Fatal("could not enable rescue")
    }

    s.server.RescueEnabled = true
    return s
}

// WaitForRescueDisabled which is disabled after the reboot
func (s *ServerInstance) WaitForRescueDisabled() *ServerInstance {

	c := s.spec.cc
	server := s.server

	//TODO: Timeout

	for server.RescueEnabled != false {
		result, _, err := c.client.Server.GetByName(s.context, server.Name)
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

	action, _, err := c.client.Server.Poweron(s.context, server)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.waitForAction(action); err != nil {
		log.Fatal("could not power on the server")
    }

    fmt.Printf("Server %d started\n", server.ID)

	return s
}

// Reboot the cloud server
func (s *ServerInstance) Reboot() *ServerInstance {

	c := s.spec.cc
	//server := s.server

	server, _, err := c.client.Server.GetByName(s.context, s.server.Name)
		if err != nil {
		log.Fatal(err)
	}
	if server == nil {
		log.Fatal("server not found")
	}

	action, _, err := c.client.Server.Reboot(s.context, server)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.waitForAction(action); err != nil {
		log.Fatal("could not reboot server")
    }

    fmt.Printf("Server %d rebooted\n", server.ID)

 	return s
}

// Create an image
func (s *ServerInstance) CreateSnapshot(description string) *ImageSpec {

	c := s.spec.cc
	server := s.server

	opts := &hcloud.ServerCreateImageOpts{
		Type: "snapshot",
		Description: &description,
	}
	result, _, err :=  c.client.Server.CreateImage(s.context, server, opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.waitForAction(result.Action); err != nil {
		log.Fatal("could not create server image")
    }

    fmt.Printf("Server image created.\n")

	return &ImageSpec {
		spec: result.Image,
	}
}
