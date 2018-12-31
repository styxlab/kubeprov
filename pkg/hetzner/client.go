package hetzner

import (
	"context"
    "os"

    "github.com/hetznercloud/hcloud-go/hcloud"
 	"github.com/go-kit/kit/log/term"
 	"github.com/gosuri/uiprogress"
)

// Client holds the connection data for the Hetzner Cloud interaction
type Client struct {
	context context.Context
	config *Config
	client *hcloud.Client
	sshKey *hcloud.SSHKey
}

// Connect opens a new Hetzner Cloud connection
func Connect() *Client {
	
	config := GetOrCreateConfig()
	
	client :=  &Client {
		client: hcloud.NewClient(hcloud.WithToken(config.token)),
		context: config.context,
		config: config,
	}
	
	return client.getOrCreateSSHKey()
}

func (c *Client) waitForAction(action *hcloud.Action) error {

	progress, errs := c.client.Action.WatchProgress(c.context, action)

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
