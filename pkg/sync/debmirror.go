package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"repository-mirror/helper"
	"repository-mirror/pkg/preset"
)

func debMirror(cfg Config, p preset.Preset) error {
	importKeysIfNeeded(p, cfg.Keyring)

	dists := strings.Split(cfg.Dists, ",")

	if len(p.DistMap) > 0 {
		for _, dist := range dists {
			dist = strings.TrimSpace(dist)
			if dist == "" {
				continue
			}
			mirrorURL := p.ExpandDistURL(dist)
			method, host, root, err := helper.ParseMirrorURL(mirrorURL)
			if err != nil {
				return fmt.Errorf("invalid mirror URL for dist %s: %w", dist, err)
			}

			destDir := cfg.DestinationDir + "/" + dist

			args := []string{
				"--method=" + method,
				"--host=" + host,
				"--root=" + root,
				"--arch=" + cfg.Architectures,
				"--dist=" + dist,
				"--rsync-extra=none",
				"--no-source",
				"--getcontents",
			}

			if cfg.Sections != "" {
				args = append(args, "--section="+cfg.Sections)
			}
			if cfg.Progress {
				args = append(args, "--progress")
			}
			if cfg.Verbose {
				args = append(args, "--verbose")
			}
			if cfg.IgnoreMissing {
				args = append(args, "--ignore-missing-release")
			}
			if cfg.IgnoreRelease {
				args = append(args, "--ignore-release-gpg")
			}
			if cfg.RsyncOptions != "" {
				args = append(args, "--rsync-options="+cfg.RsyncOptions)
			}
			if cfg.Timeout != "" {
				args = append(args, "--timeout="+cfg.Timeout)
			}
			if cfg.Keyring != "" {
				args = append(args, "--keyring="+cfg.Keyring)
			}

			args = append(args, destDir)

			if cfg.DryRun {
				fmt.Println("debmirror " + strings.Join(args, " "))
				continue
			}

			if err := os.MkdirAll(destDir, 0755); err != nil {
				return fmt.Errorf("creating destination directory: %w", err)
			}

			cmd := exec.Command("debmirror", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Printf("Starting debmirror sync from %s to %s\n", mirrorURL, destDir)
			if cfg.Sections != "" {
				fmt.Printf("  Sections: %s\n", cfg.Sections)
			}
			fmt.Printf("  Architectures: %s\n", cfg.Architectures)
			fmt.Printf("  Dist: %s\n", dist)

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("debmirror failed for dist %s: %w", dist, err)
			}
		}
		return nil
	}

	method, host, root, err := helper.ParseMirrorURL(cfg.MirrorURL)
	if err != nil {
		return fmt.Errorf("invalid mirror URL: %w", err)
	}

	destDir := cfg.DestinationDir
	archs := strings.Split(cfg.Architectures, ",")
	if len(archs) == 1 && archs[0] != "amd64" {
		destDir = cfg.DestinationDir + "_" + archs[0]
	}

	args := []string{
		"--method=" + method,
		"--host=" + host,
		"--root=" + root,
		"--arch=" + cfg.Architectures,
		"--dist=" + cfg.Dists,
		"--rsync-extra=none",
		"--no-source",
		"--getcontents",
	}

	if cfg.Sections != "" {
		args = append(args, "--section="+cfg.Sections)
	}
	if cfg.Progress {
		args = append(args, "--progress")
	}
	if cfg.Verbose {
		args = append(args, "--verbose")
	}
	if cfg.IgnoreMissing {
		args = append(args, "--ignore-missing-release")
	}
	if cfg.IgnoreRelease {
		args = append(args, "--ignore-release-gpg")
	}
	if cfg.RsyncOptions != "" {
		args = append(args, "--rsync-options="+cfg.RsyncOptions)
	}
	if cfg.Timeout != "" {
		args = append(args, "--timeout="+cfg.Timeout)
	}
	if cfg.Keyring != "" {
		args = append(args, "--keyring="+cfg.Keyring)
	}

	args = append(args, destDir)

	if cfg.DryRun {
		fmt.Println("debmirror " + strings.Join(args, " "))
		return nil
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	cmd := exec.Command("debmirror", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting debmirror sync from %s to %s\n", cfg.MirrorURL, destDir)
	if cfg.Sections != "" {
		fmt.Printf("  Sections: %s\n", cfg.Sections)
	}
	fmt.Printf("  Architectures: %s\n", cfg.Architectures)
	fmt.Printf("  Dists: %s\n", cfg.Dists)

	return cmd.Run()
}
