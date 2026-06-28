package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"repository-mirror/pkg/sync"
)

func parseConfigFile(path string) (sync.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return sync.Config{}, fmt.Errorf("reading config file %s: %w", path, err)
	}
	defer f.Close()

	vars := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "=(") {
			parts := strings.SplitN(line, "=(", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimRight(strings.TrimSpace(parts[1]), ")")
			val = strings.TrimSpace(val)
			if val == "" {
				vars[key] = ""
			} else {
				vars[key] = strings.Join(strings.Fields(val), ",")
			}
		} else {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, "'\"")
			vars[key] = val
		}
	}
	if err := scanner.Err(); err != nil {
		return sync.Config{}, fmt.Errorf("reading config file %s: %w", path, err)
	}

	cfg := sync.Config{
		Identifier:     vars["IDENTIFIER"],
		DestinationDir: vars["DESTINATION_DIR"],
		MirrorURL:      vars["MIRROR_URL"],
		Exclude:        vars["EXCLUDE"],
		Sections:       vars["SECTIONS"],
		Architectures:  vars["ARCHITECTURES"],
		Dists:          vars["DISTS"],
		Keyring:        vars["KEYRING"],
		SyncType:       vars["SYNC"],
	}

	if syncStr := vars["SYNC"]; syncStr == "false" || syncStr == "no" || syncStr == "0" {
		fmt.Println("SYNC is disabled, skipping")
		cfg.DryRun = true
	}

	for _, v := range []string{"true", "yes", "1"} {
		if vars["DRY_RUN"] == v {
			cfg.DryRun = true
		}
		if vars["VERBOSE"] == v {
			cfg.Verbose = true
		}
		if vars["PROGRESS"] == v {
			cfg.Progress = true
		}
	}

	return cfg, nil
}
