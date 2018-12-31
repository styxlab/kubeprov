package hetzner

import (
	"context"
    "os"

    "github.com/hetznercloud/hcloud-go/hcloud"
 	"github.com/go-kit/kit/log/term"
 	"github.com/thcyron/uiprogress"
 	"github.com/wayneashleyberry/terminal-dimensions"
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

	progressCh, errCh := c.client.Action.WatchProgress(c.context, action)

	if term.IsTerminal(os.Stdout) {
		progress := uiprogress.New()

		progress.Start()
		bar := progress.AddBar(100).AppendCompleted().PrependElapsed()
		w, _ := terminaldimensions.Width()
		bar.Width = int(w)-20
		bar.Empty = ' '

		for {
			select {
			case err := <-errCh:
				if err == nil {
					bar.Set(100)
				}
				progress.Stop()
				return err
			case p := <-progressCh:
				bar.Set(p)
			}
		}
	}
	return <-errCh
}