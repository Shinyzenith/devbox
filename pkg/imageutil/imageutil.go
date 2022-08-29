package imageutil

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
)

func GetImage(ctx context.Context, client *containerd.Client, image_name string, image_tag string) (containerd.Image, error) {
	logrus.Debugf("Fetching %s:%s image", image_name, image_tag)
	image_url, err := resolveShortnameUrl(image_name, image_tag)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	image, err := client.GetImage(ctx, image_url)
	if err != nil {
		logrus.Debugf("%s:%s not found, Pulling from source.\n", image_name, image_tag)
		image, err := client.Pull(ctx, image_url, containerd.WithPullUnpack)
		if err != nil {
			return nil, err
		}
		return image, nil
	}
	logrus.Debugf("Found %s:%s image, not pulling.\n", image_name, image_tag)
	return image, nil
}

func resolveShortnameUrl(image_name string, image_tag string) (string, error) {
	if image_url, exists := getShortNames()[image_name]; !exists {
		return "", fmt.Errorf("Failed to resolve shortname `%s` to image destination url.", image_name)
	} else {
		image_url = image_url + ":" + image_tag
		return image_url, nil
	}
}

func getShortNames() map[string]string {
	val := map[string]string{
		"almalinux":                    "docker.io/library/almalinux",
		"almalinux-minimal":            "docker.io/library/almalinux-minimal",
		"archlinux":                    "docker.io/archlinux/archlinux",
		"centos":                       "quay.io/centos/centos",
		"alpine":                       "docker.io/library/alpine",
		"fedora-minimal":               "registry.fedoraproject.org/fedora-minimal",
		"fedora":                       "registry.fedoraproject.org/fedora",
		"opensuse/tumbleweed":          "registry.opensuse.org/opensuse/tumbleweed",
		"opensuse/tumbleweed-dnf":      "registry.opensuse.org/opensuse/tumbleweed-dnf",
		"opensuse/tumbleweed-microdnf": "registry.opensuse.org/opensuse/tumbleweed-microdnf",
		"opensuse/leap":                "registry.opensuse.org/opensuse/leap",
		"opensuse/busybox":             "registry.opensuse.org/opensuse/busybox",
		"tumbleweed":                   "registry.opensuse.org/opensuse/tumbleweed",
		"tumbleweed-dnf":               "registry.opensuse.org/opensuse/tumbleweed-dnf",
		"tumbleweed-microdnf":          "registry.opensuse.org/opensuse/tumbleweed-microdnf",
		"leap":                         "registry.opensuse.org/opensuse/leap",
		"leap-dnf":                     "registry.opensuse.org/opensuse/leap-dnf",
		"leap-microdnf":                "registry.opensuse.org/opensuse/leap-microdnf",
		"tw-busybox":                   "registry.opensuse.org/opensuse/busybox",
		"suse/sle15":                   "registry.suse.com/suse/sle15",
		"suse/sles12sp5":               "registry.suse.com/suse/sles12sp5",
		"suse/sles12sp4":               "registry.suse.com/suse/sles12sp4",
		"suse/sles12sp3":               "registry.suse.com/suse/sles12sp3",
		"sle15":                        "registry.suse.com/suse/sle15",
		"sles12sp5":                    "registry.suse.com/suse/sles12sp5",
		"sles12sp4":                    "registry.suse.com/suse/sles12sp4",
		"sles12sp3":                    "registry.suse.com/suse/sles12sp3",
		"rhel":                         "registry.access.redhat.com/rhel",
		"rhel6":                        "registry.access.redhat.com/rhel6",
		"rhel7":                        "registry.access.redhat.com/rhel7",
		"rhel7.9":                      "registry.access.redhat.com/rhel7.9",
		"rhel-atomic":                  "registry.access.redhat.com/rhel-atomic",
		"rhel-minimal":                 "registry.access.redhat.com/rhel-minimum",
		"rhel-init":                    "registry.access.redhat.com/rhel-init",
		"rhel7-atomic":                 "registry.access.redhat.com/rhel7-atomic",
		"rhel7-minimal":                "registry.access.redhat.com/rhel7-minimum",
		"rhel7-init":                   "registry.access.redhat.com/rhel7-init",
		"rhel7/rhel":                   "registry.access.redhat.com/rhel7/rhel",
		"rhel7/rhel-atomic":            "registry.access.redhat.com/rhel7/rhel7/rhel-atomic",
		"ubi7/ubi":                     "registry.access.redhat.com/ubi7/ubi",
		"ubi7/ubi-minimal":             "registry.access.redhat.com/ubi7-minimal",
		"ubi7/ubi-init":                "registry.access.redhat.com/ubi7-init",
		"ubi7":                         "registry.access.redhat.com/ubi7",
		"ubi7-init":                    "registry.access.redhat.com/ubi7-init",
		"ubi7-minimal":                 "registry.access.redhat.com/ubi7-minimal",
		"rhel8":                        "registry.access.redhat.com/ubi8",
		"rhel8-init":                   "registry.access.redhat.com/ubi8-init",
		"rhel8-minimal":                "registry.access.redhat.com/ubi8-minimal",
		"rhel8-micro":                  "registry.access.redhat.com/ubi8-micro",
		"ubi8":                         "registry.access.redhat.com/ubi8",
		"ubi8-minimal":                 "registry.access.redhat.com/ubi8-minimal",
		"ubi8-init":                    "registry.access.redhat.com/ubi8-init",
		"ubi8-micro":                   "registry.access.redhat.com/ubi8-micro",
		"ubi8/ubi":                     "registry.access.redhat.com/ubi8/ubi",
		"ubi8/ubi-minimal":             "registry.access.redhat.com/ubi8-minimal",
		"ubi8/ubi-init":                "registry.access.redhat.com/ubi8-init",
		"ubi8/ubi-micro":               "registry.access.redhat.com/ubi8-micro",
		"rhel9":                        "registry.access.redhat.com/ubi9",
		"rhel9-init":                   "registry.access.redhat.com/ubi9-init",
		"rhel9-minimal":                "registry.access.redhat.com/ubi9-minimal",
		"rhel9-micro":                  "registry.access.redhat.com/ubi9-micro",
		"ubi9":                         "registry.access.redhat.com/ubi9",
		"ubi9-minimal":                 "registry.access.redhat.com/ubi9-minimal",
		"ubi9-init":                    "registry.access.redhat.com/ubi9-init",
		"ubi9-micro":                   "registry.access.redhat.com/ubi9-micro",
		"ubi9/ubi":                     "registry.access.redhat.com/ubi9/ubi",
		"ubi9/ubi-minimal":             "registry.access.redhat.com/ubi9-minimal",
		"ubi9/ubi-init":                "registry.access.redhat.com/ubi9-init",
		"ubi9/ubi-micro":               "registry.access.redhat.com/ubi9-micro",
		"rockylinux":                   "docker.io/library/rockylinux",
		"debian":                       "docker.io/library/debian",
		"ubuntu":                       "docker.io/library/ubuntu",
		"oraclelinux":                  "container-registry.oracle.com/os/oraclelinux",
	}
	return val
}
