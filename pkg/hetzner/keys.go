package hetzner

import (
	"fmt"
    "log"

    "github.com/hetznercloud/hcloud-go/hcloud"
)

func (c *Client) getOrCreateSSHKey() *Client {

	config := c.config

	sshKey, _, err := c.client.SSHKey.Get(c.context, config.publicKeyName)
	if err != nil {
		log.Fatal(err)
	}
    if sshKey == nil {
        sshKey = c.createSSHKey().sshKey
    }
    
    c.sshKey = sshKey
    return c
}

func (c *Client) createSSHKey() *Client {

	config := c.config

	opts := hcloud.SSHKeyCreateOpts{
		Name:      config.publicKeyName,
		PublicKey: config.publicKey,
	}
	sshKey, _, err := c.client.SSHKey.Create(c.context, opts)
	if err != nil {
		log.Fatalf("we got some problem with the SSH-Key: %s.", err)
	}
	fmt.Printf("SSH key %d created\n", sshKey.ID)

	c.sshKey = sshKey
	return c
}

func (c *Client) getSSHKey() *Client {

	config := c.config

	sshKey, _, err := c.client.SSHKey.Get(c.context, config.publicKeyName)
	if err != nil {
		log.Fatal(err)
	}
    if sshKey == nil {
        log.Fatalf("we got some problem with the SSH-Key")
    }
    
    c.sshKey = sshKey
    return c
}

