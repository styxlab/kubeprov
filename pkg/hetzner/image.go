package hetzner

import (
	"context"
	"fmt"
	"log"

    "github.com/hetznercloud/hcloud-go/hcloud"
)

// ImageSpec hold an (incomplete) image specification
type ImageSpec struct {
	spec *hcloud.Image
}

func ImageByName(name string) *ImageSpec {
	return &ImageSpec {
		spec: &hcloud.Image{
			Name: name,
		},
	}
}

func ImageByID(id int) *ImageSpec {
	return &ImageSpec {
		spec: &hcloud.Image{
			ID: id,
		},
	}
}

func (c *Client) ImageDelete(image *ImageSpec) *Client {

	_, err := c.client.Image.Delete(context.Background(), image.spec)
	if err != nil {
		log.Fatal(err)
	}

    fmt.Printf("Image %d deleted\n", image.spec.ID)

    return c
}