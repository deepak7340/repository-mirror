package config

import (
	"github.com/spf13/viper"

	"repository-mirror/constant"
)

var v *viper.Viper

func init() {
	v = viper.New()
	v.AutomaticEnv()
}

func GetIdentifier() string    { return v.GetString(constant.FlagID) }
func GetSections() string      { return v.GetString(constant.FlagSections) }
func GetArchitectures() string { return v.GetString(constant.FlagArchitectures) }
func GetDists() string         { return v.GetString(constant.FlagDists) }
func GetMirrorURL() string     { return v.GetString(constant.FlagMirrorURL) }
func GetDest() string          { return v.GetString(constant.FlagDest) }
func GetExclude() string       { return v.GetString(constant.FlagExclude) }
func GetRsyncOptions() string  { return v.GetString(constant.FlagRsyncOptions) }
func GetTimeout() string       { return v.GetString(constant.FlagTimeout) }
func GetKeyring() string       { return v.GetString(constant.FlagKeyring) }
func IsProgress() bool         { return v.GetBool(constant.FlagProgress) }
func IsDryRun() bool           { return v.GetBool(constant.FlagDryRun) }
func IsVerbose() bool          { return v.GetBool(constant.FlagVerbose) }
func IsIgnoreMissing() bool    { return v.GetBool(constant.FlagIgnoreMissing) }
func IsIgnoreRelease() bool    { return v.GetBool(constant.FlagIgnoreRelease) }
func IsList() bool             { return v.GetBool(constant.FlagList) }

func GetViperInstance() *viper.Viper { return v }
