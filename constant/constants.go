package constant

const (
	FlagID            = "id"
	FlagSections      = "sections"
	FlagArchitectures = "architectures"
	FlagDists         = "dists"
	FlagMirrorURL     = "mirror-url"
	FlagDest          = "dest"
	FlagExclude       = "exclude"
	FlagProgress      = "progress"
	FlagDryRun        = "dry-run"
	FlagVerbose       = "verbose"
	FlagIgnoreMissing = "ignore-missing"
	FlagIgnoreRelease = "ignore-release"
	FlagRsyncOptions  = "rsync-options"
	FlagTimeout       = "timeout"
	FlagKeyring       = "keyring"
	FlagList          = "list"
)

const (
	DefaultKeyring    = "/var/cache/packagesign/keyrings/trustedkeys.gpg"
	DefaultKeyringDir = "/var/cache/packagesign/keyrings"
)
