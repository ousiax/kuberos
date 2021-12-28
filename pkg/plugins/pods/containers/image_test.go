package containers

import (
	"testing"
)

func TestParseImageRef(t *testing.T) {
	var tests = []struct {
		image    string
		registry string
		org      string
		name     string
		tag      string
		digest   string
		parsed   bool
	}{
		{
			image:    "k8s.gcr.io/ingress-nginx/controller:v1.1.0@sha256:f766669fdcf3dc26347ed273a55e754b427eb4411ee075a53f30718b4499076a",
			registry: "k8s.gcr.io",
			org:      "ingress-nginx",
			name:     "controller",
			tag:      "v1.1.0",
			digest:   "sha256:f766669fdcf3dc26347ed273a55e754b427eb4411ee075a53f30718b4499076a",
			parsed:   true,
		},
		{
			image:    "k8s.gcr.io/ingress-nginx/controller:v1.1.0",
			registry: "k8s.gcr.io",
			org:      "ingress-nginx",
			name:     "controller",
			tag:      "v1.1.0",
			parsed:   true,
		},
		{
			image:  "k8s.gcr.io/ingress-nginx/controller",
			parsed: false,
		},
		{
			image:  "controller",
			parsed: false,
		},
		{
			image:  "controller:latest",
			parsed: false,
		},
	}
	for _, test := range tests {
		ref, err := ParseImageRef(test.image)
		if err != nil && test.parsed {
			t.Error(err)
		}

		if err != nil && !test.parsed {
			return
		}

		if ref.Registry != test.registry {
			t.Errorf("image.registry, got: (%q),want: %q", ref.Registry, test.registry)
		}

		if ref.Org != test.org {
			t.Errorf("image.org, got: (%q),want: %q", ref.Org, test.org)
		}

		if ref.Name != test.name {
			t.Errorf("image.name, got: (%q),want: %q", ref.Name, test.name)
		}

		if ref.Tag != test.tag {
			t.Errorf("image.tag, got: (%q),want: %q", ref.Tag, test.tag)
		}

		if ref.Digest != test.digest {
			t.Errorf("image.digest, got: (%q),want: %q", ref.Digest, test.digest)
		}
	}
}
