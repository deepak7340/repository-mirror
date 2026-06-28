package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"repository-mirror/helper"
	"repository-mirror/pkg/preset"
)

func syncWget(cfg Config, p preset.Preset) error {
	mirrorURL := strings.TrimRight(cfg.MirrorURL, "/")

	dirExcludes, rejectPatterns := wgetExcludes(p, cfg)
	cutDirs := helper.CountPathSegments(cfg.MirrorURL)

	buildArgs := func(srcURL string, overrideCutDirs ...int) []string {
		cd := cutDirs
		if len(overrideCutDirs) > 0 {
			cd = overrideCutDirs[0]
		}
		args := []string{
			"--mirror", "--no-parent", "-nH",
			"--execute", "robots=off",
			fmt.Sprintf("--cut-dirs=%d", cd),
			"-P", cfg.DestinationDir, srcURL,
		}
		if cfg.Progress {
			args = append(args, "--progress=dot:giga")
		}
		var rejectRegex []string
		for _, ex := range dirExcludes {
			rejectRegex = append(rejectRegex, ex+"/.*")
		}
		if len(rejectRegex) > 0 {
			args = append(args, "--reject-regex="+strings.Join(rejectRegex, "|"))
		}
		if len(rejectPatterns) > 0 {
			args = append(args, "--reject="+strings.Join(rejectPatterns, ","))
		}
		return args
	}

	if !cfg.DryRun {
		for _, ex := range dirExcludes {
			cmd := exec.Command("find", cfg.DestinationDir, "-type", "d", "-name", ex, "-exec", "rm", "-rf", "{}", "+")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}

	if p.Flat {
		if len(p.DistMap) > 0 || strings.Contains(cfg.MirrorURL, "{dist}") {
			dists := strings.Split(cfg.Dists, ",")
			sections := strings.Split(cfg.Sections, ",")
			hasDistMap := strings.Contains(cfg.MirrorURL, "{dist_map}")
			for _, dist := range dists {
				dist = strings.TrimSpace(dist)
				if dist == "" {
					continue
				}
				expandedURL := p.ExpandDistURL(dist)
				sectionList := sections
				if len(sections) == 1 && sections[0] == "" {
					sectionList = []string{""}
				}
				for _, section := range sectionList {
					section = strings.TrimSpace(section)
					var srcURL string
					if hasDistMap {
						srcURL = expandedURL + "/" + dist
						if section != "" {
							srcURL += "/" + section
						}
						srcURL += "/"
					} else {
						srcURL = expandedURL
						if section != "" {
							srcURL += "/" + section
						}
						srcURL += "/"
					}
					cd := cutDirs
					if !hasDistMap {
						cd = 0
					}
					args := buildArgs(srcURL, cd)

					if cfg.DryRun {
						fmt.Printf("wget %s\n", strings.Join(args, " "))
						continue
					}
					if err := os.MkdirAll(cfg.DestinationDir, 0755); err != nil {
						return fmt.Errorf("creating destination directory %s: %w", cfg.DestinationDir, err)
					}
					cmd := exec.Command("wget", args...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					fmt.Printf("Syncing %s to %s\n", srcURL, cfg.DestinationDir)
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("wget failed: %w", err)
					}
				}
			}
			return nil
		}
		srcURL := mirrorURL + "/"
		args := buildArgs(srcURL)

		if cfg.DryRun {
			fmt.Printf("wget %s\n", strings.Join(args, " "))
			return nil
		}
		if err := os.MkdirAll(cfg.DestinationDir, 0755); err != nil {
			return fmt.Errorf("creating destination directory %s: %w", cfg.DestinationDir, err)
		}
		cmd := exec.Command("wget", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("Syncing %s to %s\n", mirrorURL, cfg.DestinationDir)
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
				args := buildArgs(srcURL)

				if cfg.DryRun {
					fmt.Printf("wget %s\n", strings.Join(args, " "))
					continue
				}
				if err := os.MkdirAll(cfg.DestinationDir, 0755); err != nil {
					return fmt.Errorf("creating destination directory %s: %w", cfg.DestinationDir, err)
				}
				cmd := exec.Command("wget", args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				fmt.Printf("Syncing from %s to %s\n", srcURL, cfg.DestinationDir)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("wget failed: %w", err)
				}
			}
		}
	}

	return nil
}
