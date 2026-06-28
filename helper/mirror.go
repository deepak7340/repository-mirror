package helper

import (
	"fmt"
	"strings"
)

func ParseMirrorURL(rawURL string) (method, host, root string, err error) {
	if strings.HasPrefix(rawURL, "https://") {
		method = "https"
		rawURL = strings.TrimPrefix(rawURL, "https://")
	} else if strings.HasPrefix(rawURL, "http://") {
		method = "http"
		rawURL = strings.TrimPrefix(rawURL, "http://")
	} else if strings.HasPrefix(rawURL, "ftp://") {
		method = "ftp"
		rawURL = strings.TrimPrefix(rawURL, "ftp://")
	} else if strings.HasPrefix(rawURL, "rsync://") {
		method = "rsync"
		rawURL = strings.TrimPrefix(rawURL, "rsync://")
	} else {
		return "", "", "", fmt.Errorf("unsupported protocol in URL: %s", rawURL)
	}

	slashIdx := strings.Index(rawURL, "/")
	if slashIdx == -1 {
		host = rawURL
		root = ""
	} else {
		host = rawURL[:slashIdx]
		root = strings.Trim(rawURL[slashIdx:], "/")
	}

	return method, host, root, nil
}

func CountPathSegments(rawURL string) int {
	if strings.HasPrefix(rawURL, "https://") || strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "ftp://") || strings.HasPrefix(rawURL, "rsync://") {
		parts := strings.SplitN(rawURL, "://", 2)
		if len(parts) == 2 {
			afterProto := parts[1]
			slashIdx := strings.Index(afterProto, "/")
			if slashIdx != -1 {
				path := afterProto[slashIdx+1:]
				path = strings.TrimRight(path, "/")
				if path == "" {
					return 0
				}
				segs := strings.Split(path, "/")
				return len(segs)
			}
		}
	}
	return 0
}
