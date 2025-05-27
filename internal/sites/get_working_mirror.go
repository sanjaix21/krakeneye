package sites

import (
	"errors"
	"net/http"
	"time"
)

type MirrorResult struct {
	SiteName string
	Mirror   string
}

// tires all mirrors till it find the first working one
func FindFirstWorkingMirror() (*MirrorResult, error) {
	for _, site := range PiracySites {
		for _, mirror := range site.Mirrors {
			client := http.Client{Timeout: 5 * time.Second}
			resp, err := client.Head(mirror)
			if err == nil && resp.StatusCode == 200 {
				return &MirrorResult{
					SiteName: site.Name,
					Mirror:   mirror,
				}, nil
			}
		}
	}

	return nil, errors.New("no working mirror found for any site")
}
