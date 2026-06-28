package sync

import (
	"fmt"
	"strings"

	"repository-mirror/pkg/gpg"
	"repository-mirror/pkg/preset"
)

type Config struct {
	Identifier     string
	Sections       string
	Architectures  string
	Dists          string
	MirrorURL      string
	DestinationDir string
	Exclude        string
	Progress       bool
	DryRun         bool
	Verbose        bool
	IgnoreMissing  bool
	IgnoreRelease  bool
	RsyncOptions   string
	Timeout        string
	Keyring        string
	SyncType       string
}

func wgetExcludes(p preset.Preset, cfg Config) (dirExcludes, rejectPatterns []string) {
	all := p.Exclude
	if cfg.Exclude != "" {
		if all != "" {
			all += ","
		}
		all += cfg.Exclude
	}
	if all == "" {
		return nil, []string{"index.html*"}
	}
	for _, ex := range strings.Split(all, ",") {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		if strings.ContainsAny(ex, "*?") || strings.HasSuffix(ex, ".html") {
			rejectPatterns = append(rejectPatterns, ex)
		} else {
			dirExcludes = append(dirExcludes, ex)
		}
	}
	return
}

func Run(cfg Config) error {
	p, ok := preset.Presets[cfg.Identifier]
	if !ok {
		if cfg.SyncType == "" {
			return fmt.Errorf("unknown identifier: %s", cfg.Identifier)
		}
		p = preset.Preset{
			Dest:          cfg.DestinationDir,
			MirrorURL:     cfg.MirrorURL,
			Sections:      cfg.Sections,
			Architectures: cfg.Architectures,
			Dists:         cfg.Dists,
			SyncType:      cfg.SyncType,
			Exclude:       cfg.Exclude,
		}
	}

	if cfg.DestinationDir == "" {
		cfg.DestinationDir = p.Dest
	}
	if cfg.MirrorURL == "" {
		cfg.MirrorURL = p.MirrorURL
	}
	if cfg.Sections == "" {
		cfg.Sections = p.Sections
	}
	if cfg.Architectures == "" {
		cfg.Architectures = p.Architectures
	}
	if cfg.Dists == "" {
		cfg.Dists = p.Dists
	}

	switch p.SyncType {
	case "debmirror":
		return debMirror(cfg, p)
	case "rsync":
		return syncRsync(cfg, p)
	case "wget":
		return syncWget(cfg, p)
	default:
		return fmt.Errorf("unsupported sync type: %s", p.SyncType)
	}
}

func importKeysIfNeeded(p preset.Preset, keyring string) {
	if p.GPGKeyURL != "" || len(p.GPGKeyIDs) > 0 {
		gpg.ImportGPGKeys(p.GPGKeyURL, p.GPGKeyIDs, keyring)
	}
}
