package hetzner

import (
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


