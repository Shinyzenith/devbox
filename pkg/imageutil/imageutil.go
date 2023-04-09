/*
* SPDX-License-Identifier: GPL-3.0-only
*
* imageutil.go
*
* Created by:	Aakash Sen Sharma
* Copyright:	(C) 2022, Aakash Sen Sharma & Contributors
 */

package imageutil

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
)

func GetImage(ctx context.Context, client *containerd.Client, image_name string, image_tag string) (containerd.Image, error) {
	logrus.Debugf("Fetching %s:%s image", image_name, image_tag)
	image_url, err := resolveShortnameUrl(image_name, image_tag)
	if err != nil {
		return nil, err
	}

	image, err := client.GetImage(ctx, image_url)
	if err != nil {
		logrus.Debugf("%s:%s not found, Pulling from source", image_name, image_tag)
		image, err := client.Pull(ctx, image_url, containerd.WithPullUnpack)
		if err != nil {
			return nil, err
		}
		return image, nil
	}
	logrus.Debugf("Found %s:%s image, not pulling", image_name, image_tag)
	return image, nil
}

func resolveShortnameUrl(image_name string, image_tag string) (string, error) {
	if image_url, exists := getShortNames()[image_name]; !exists {
		return "", fmt.Errorf("Failed to resolve shortname `%s` to image destination url", image_name)
	} else {
		image_url = image_url + ":" + image_tag
		return image_url, nil
	}
}
