package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"repository-mirror/config"
	"repository-mirror/constant"
	"repository-mirror/pkg/preset"
	"repository-mirror/pkg/sync"
)

var rootCmd = &cobra.Command{
	Use:          "repo-sync [flags] [<config-file>]",
	Short:        "Sync apt/rpm repositories to a local cache",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cfg, err := parseConfigFile(args[0])
			if err != nil {
				return err
			}
			return sync.Run(cfg)
		}

		v := config.GetViperInstance()

		if config.IsList() || v.GetString(constant.FlagID) == "help" || v.GetString(constant.FlagID) == "--help" {
			printPresets()
			return nil
		}

		cfg := sync.Config{
			Identifier:     config.GetIdentifier(),
			Sections:       config.GetSections(),
			Architectures:  config.GetArchitectures(),
			Dists:          config.GetDists(),
			MirrorURL:      config.GetMirrorURL(),
			DestinationDir: config.GetDest(),
			Exclude:        config.GetExclude(),
			Progress:       config.IsProgress(),
			DryRun:         config.IsDryRun(),
			Verbose:        config.IsVerbose(),
			IgnoreMissing:  config.IsIgnoreMissing(),
			IgnoreRelease:  config.IsIgnoreRelease(),
			RsyncOptions:   config.GetRsyncOptions(),
			Timeout:        config.GetTimeout(),
			Keyring:        config.GetKeyring(),
		}

		return sync.Run(cfg)
	},
}

func printPresets() {
	names := make([]string, 0, len(preset.Presets))
	for k := range preset.Presets {
		names = append(names, k)
	}
	sort.Strings(names)

	fmt.Println("Available presets:")
	for _, name := range names {
		p := preset.Presets[name]
		fmt.Printf("  %-25s %-4s %-10s %s\n", name, p.Format, p.SyncType, p.Dest)
	}
	fmt.Println("\nUse --list or --id help to see this list. Use --id <name> to select a preset.")
}

func init() {
	v := config.GetViperInstance()

	rootCmd.Flags().String(constant.FlagID, "ubuntu", "Repository identifier (use --list to see presets)")
	rootCmd.Flags().String(constant.FlagSections, "", "Comma-separated list of sections")
	rootCmd.Flags().String(constant.FlagArchitectures, "", "Comma-separated list of architectures")
	rootCmd.Flags().String(constant.FlagDists, "", "Comma-separated list of distributions")
	rootCmd.Flags().String(constant.FlagMirrorURL, "", "Mirror URL")
	rootCmd.Flags().String(constant.FlagDest, "", "Destination directory")
	rootCmd.Flags().String(constant.FlagExclude, "", "Comma-separated rsync/wget exclude patterns")
	rootCmd.Flags().Bool(constant.FlagProgress, false, "Show progress")
	rootCmd.Flags().Bool(constant.FlagDryRun, false, "Dry run (print command only)")
	rootCmd.Flags().Bool(constant.FlagVerbose, false, "Verbose output")
	rootCmd.Flags().Bool(constant.FlagIgnoreMissing, false, "Ignore missing files")
	rootCmd.Flags().Bool(constant.FlagIgnoreRelease, false, "Ignore release gpg errors")
	rootCmd.Flags().String(constant.FlagRsyncOptions, "", "Additional rsync options")
	rootCmd.Flags().String(constant.FlagTimeout, "", "Timeout for debmirror")
	rootCmd.Flags().String(constant.FlagKeyring, constant.DefaultKeyring, "GPG keyring path")
	rootCmd.Flags().Bool(constant.FlagList, false, "List available presets")
	rootCmd.Flags().Lookup(constant.FlagMirrorURL).Hidden = true

	for _, name := range []string{
		constant.FlagID, constant.FlagSections, constant.FlagArchitectures,
		constant.FlagDists, constant.FlagMirrorURL, constant.FlagDest,
		constant.FlagExclude, constant.FlagProgress, constant.FlagDryRun,
		constant.FlagVerbose, constant.FlagIgnoreMissing, constant.FlagIgnoreRelease,
		constant.FlagRsyncOptions, constant.FlagTimeout, constant.FlagKeyring,
		constant.FlagList,
	} {
		v.BindPFlag(name, rootCmd.Flags().Lookup(name))
	}

	v.BindEnv(constant.FlagID, "REPO_SYNC_ID")
	v.BindEnv(constant.FlagSections, "REPO_SYNC_SECTIONS")
	v.BindEnv(constant.FlagDryRun, "REPO_SYNC_DRY_RUN")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
