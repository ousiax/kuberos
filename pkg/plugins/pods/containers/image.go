package containers

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile(`^(?P<registry>[^/]+)/(?P<org>[^/]+)/(?P<name>[^:]+):(?P<tag>[^@]+)@?(?P<digest>.*)?$`)
)

type ImageRef struct {
	Registry string
	Org      string
	Name     string
	Tag      string
	Digest   string
}

func ParseImageRef(image string) (*ImageRef, error) {
	if !re.MatchString(image) {
		return nil, fmt.Errorf("`image: %s` must match (registry)/(org)/(name):(tag)[@digest]", image)
	}

	matches := re.FindStringSubmatch(image)
	ref := &ImageRef{
		Registry: matches[1],
		Org:      matches[2],
		Name:     matches[3],
		Tag:      matches[4],
		Digest:   matches[5],
	}
	return ref, nil
}

func (i *ImageRef) String() string {
	var image strings.Builder

	if i.Registry != "" {
		image.WriteString(i.Registry)
		image.WriteString("/")
	}
	if i.Org != "" {
		image.WriteString(i.Org)
		image.WriteString("/")
	}
	if i.Name != "" {
		image.WriteString(i.Name)
	}
	if i.Tag != "" {
		image.WriteString(":")
		image.WriteString(i.Tag)
	}
	if i.Digest != "" {
		image.WriteString("@")
		image.WriteString(i.Digest)
	}
	return image.String()
}
