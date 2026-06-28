package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"repository-mirror/pkg/preset"
)

func syncRsync(cfg Config, p preset.Preset) error {
	mirrorURL := strings.TrimRight(cfg.MirrorURL, "/")

	var excludes []string
	if cfg.Exclude != "" {
		excludes = strings.Split(cfg.Exclude, ",")
	}

	if p.Flat {
		destDir := cfg.DestinationDir
		args := []string{"-avz", "--delete"}
		if cfg.Progress {
			args = append(args, "--progress")
		}
		for _, ex := range excludes {
			ex = strings.TrimSpace(ex)
			if ex != "" {
				args = append(args, "--exclude="+ex)
			}
		}
		if cfg.Verbose {
			args = append(args, "-v")
		}
		args = append(args, mirrorURL+"/", destDir)

		if cfg.DryRun {
			fmt.Printf("rsync %s\n", strings.Join(args, " "))
			return nil
		}
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("creating destination directory %s: %w", destDir, err)
		}
		cmd := exec.Command("rsync", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("Syncing %s to %s\n", cfg.MirrorURL, destDir)
		return cmd.Run()
	}

	dists := strings.Split(cfg.Dists, ",")
	archs := strings.Split(cfg.Architectures, ",")
	sections := strings.Split(cfg.Sections, ",")

	for _, dist := range dists {
		dist = strings.TrimSpace(dist)
		if dist == "" {
			continue
		}
		archList := archs
		if len(archs) == 1 && archs[0] == "" {
			archList = []string{""}
		}
		for _, arch := range archList {
			arch = strings.TrimSpace(arch)
			sectionList := sections
			if len(sections) == 1 && sections[0] == "" {
				sectionList = []string{""}
			}
			for _, section := range sectionList {
				section = strings.TrimSpace(section)
				path := p.RsyncPath
				path = strings.ReplaceAll(path, "{dist}", dist)
				path = strings.ReplaceAll(path, "{arch}", arch)
				path = strings.ReplaceAll(path, "{section}", section)
				for strings.Contains(path, "//") {
					path = strings.ReplaceAll(path, "//", "/")
				}

				srcURL := fmt.Sprintf("%s/%s", mirrorURL, strings.TrimLeft(path, "/"))
				destSubDir := cfg.DestinationDir
				parts := []string{}
				if dist != "" {
					parts = append(parts, dist)
				}
				if section != "" {
					parts = append(parts, section)
				}
				if arch != "" {
					parts = append(parts, arch)
				}
				for _, p := range parts {
					destSubDir = fmt.Sprintf("%s/%s", destSubDir, p)
				}

				rsyncArgs := []string{"-avz", "--delete"}
				if cfg.Progress {
					rsyncArgs = append(rsyncArgs, "--progress")
				}
				for _, ex := range excludes {
					ex = strings.TrimSpace(ex)
					if ex != "" {
						rsyncArgs = append(rsyncArgs, "--exclude="+ex)
					}
				}
				if cfg.Verbose {
					rsyncArgs = append(rsyncArgs, "-v")
				}
				rsyncArgs = append(rsyncArgs, srcURL, destSubDir)

				if cfg.DryRun {
					fmt.Printf("rsync %s\n", strings.Join(rsyncArgs, " "))
					continue
				}
				if err := os.MkdirAll(destSubDir, 0755); err != nil {
					return fmt.Errorf("creating destination directory %s: %w", destSubDir, err)
				}
				cmd := exec.Command("rsync", rsyncArgs...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				fmt.Printf("Syncing from %s to %s\n", srcURL, destSubDir)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("rsync failed: %w", err)
				}
			}
		}
	}

	return nil
}
