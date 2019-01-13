package hetzner

import (
	"context"
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
	action *hcloud.Action
	lastop string
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

// PublicKeyName
func (s *ServerInstance) PublicKeyName() string {
	return s.spec.cc.config.publicKeyName
}

// PrivateKeyFile
func (s *ServerInstance) PrivateKeyFile() string {
	return s.spec.cc.config.privateKeyFile
}


// Create makes a new cloud server instance
func (s *ServerSpec) Create() *ServerInstance {

	client := s.cc.client

    result, _, err := client.Server.Create(s.cc.context, s.options)
	if err != nil {
		log.Fatal(err)
	}
    return &ServerInstance {
    	context: s.cc.context,
    	spec: s,
    	server: result.Server,
    	action: result.Action,
    	lastop: "created",
    }
}

func (s *ServerInstance) Delete() *ServerInstance {
	
	c := s.spec.cc
	server := s.server

	_, err := c.client.Server.Delete(s.context, server)
	if err != nil {
		log.Fatal(err)
	}
	s.lastop = "deleted"
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
    	lastop: "status",
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
	s.action = rescue.Action
	s.lastop = "rescueEnabled"
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
	s.action = action
	s.lastop = "powerOn"
	return s
}

// Reboot the cloud server
func (s *ServerInstance) Reboot() *ServerInstance {

	c := s.spec.cc

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
	s.action = action
	s.lastop = "rebooted"
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

	return &ImageSpec {
		spec: result.Image,
		action: result.Action,
		lastop: "created",
	}
}


