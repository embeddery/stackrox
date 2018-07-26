package utils

import (
	"fmt"
	"strings"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"github.com/docker/distribution/reference"
)

// GenerateImageFromString generates an image type from a common string format
func GenerateImageFromString(imageStr string) *v1.Image {
	image := v1.Image{
		Name: &v1.ImageName{},
	}

	// Check if its a sha and return if it is
	if strings.HasPrefix(imageStr, "sha256:") {
		image.Name.Sha = imageStr
		return &image
	}

	// Cut off @sha256:
	if idx := strings.Index(imageStr, "@sha256:"); idx != -1 {
		image.Name.Sha = imageStr[idx+1:]
		imageStr = imageStr[:idx]
	}

	named, err := reference.ParseNormalizedNamed(imageStr)
	if err != nil {
		return &image
	}
	tag := "latest"
	namedTagged, ok := named.(reference.NamedTagged)
	if ok {
		tag = namedTagged.Tag()
	}
	image.Name.Registry = reference.Domain(named)
	image.Name.Remote = reference.Path(named)
	image.Name.Tag = tag
	image.Name.FullName = fmt.Sprintf("%s/%s:%s", image.Name.Registry, image.Name.Remote, image.Name.Tag)
	return &image
}

// ExtractImageSha returns the image sha if it exists within the string.
func ExtractImageSha(imageStr string) string {
	if idx := strings.Index(imageStr, "@sha256:"); idx != -1 {
		return imageStr[idx+len("@sha256:"):]
	}

	return ""
}
