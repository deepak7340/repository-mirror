package preset

import "strings"

type Preset struct {
	Dest          string
	MirrorURL     string
	Sections      string
	Architectures string
	Dists         string
	Format        string
	SyncType      string
	RsyncPath     string
	Flat          bool
	Exclude       string
	GPGKeyURL     string
	GPGKeyIDs     []string
	DistMap       map[string]string
}

func (p Preset) ExpandDistURL(dist string) string {
	mapped := dist
	if v, ok := p.DistMap[dist]; ok {
		mapped = v
	}
	result := strings.ReplaceAll(p.MirrorURL, "{dist_map}", mapped)
	result = strings.ReplaceAll(result, "{dist}", dist)
	return result
}

var Presets = map[string]Preset{
	"ubuntu": {
		Dest:          "/var/cache/packagesign/apt/ubuntu",
		MirrorURL:     "https://mirrors.dotsrc.org/ubuntu/",
		Sections:      "main,restricted,universe,multiverse",
		Architectures: "amd64",
		Dists:         "focal,focal-backports,focal-security,focal-updates,jammy,jammy-backports,jammy-security,jammy-updates,noble,noble-backports,noble-security,noble-updates,plucky,plucky-backports,plucky-security,plucky-updates,resolute,resolute-backports,resolute-security,resolute-updates",
		Format:        "deb",
		SyncType:      "debmirror",
	},
	"docker-apt": {
		Dest:          "/var/cache/packagesign/docker/apt",
		MirrorURL:     "https://download.docker.com/linux/ubuntu",
		Sections:      "stable",
		Architectures: "amd64",
		Dists:         "focal,jammy,noble",
		Format:        "deb",
		SyncType:      "debmirror",
		GPGKeyURL:     "https://download.docker.com/linux/ubuntu/gpg",
	},
	"docker-yum": {
		Dest:          "/var/cache/packagesign/docker/yum",
		MirrorURL:     "https://download.docker.com/linux/centos",
		Sections:      "stable",
		Architectures: "x86_64",
		Dists:         "7,8,9",
		Format:        "rpm",
		SyncType:      "wget",
		RsyncPath:     "{dist}/{arch}/{section}/",
	},
	"epel": {
		Dest:          "/var/cache/packagesign/yum/epel",
		MirrorURL:     "rsync://dl.fedoraproject.org/fedora-epel",
		Sections:      "",
		Architectures: "x86_64",
		Dists:         "8,9,10",
		Format:        "rpm",
		SyncType:      "rsync",
		RsyncPath:     "{dist}/Everything/{arch}/",
	},
	"jenkins": {
		Dest:      "/var/cache/packagesign/jenkins",
		MirrorURL: "https://archives.jenkins.io/{dist}",
		Dists:     "debian-stable,redhat-stable",
		Format:    "rpm",
		SyncType:  "wget",
		Flat:      true,
	},
	"microsoft-apt": {
		Dest:          "/var/cache/packagesign/microsoft/apt",
		MirrorURL:     "https://packages.microsoft.com/ubuntu/{dist_map}/prod",
		Sections:      "main",
		Architectures: "amd64",
		Dists:         "focal,jammy,noble,plucky,resolute",
		Format:        "deb",
		SyncType:      "debmirror",
		GPGKeyURL:     "https://packages.microsoft.com/keys/microsoft.asc",
		GPGKeyIDs:     []string{"EE4D7792F748182B"},
		DistMap: map[string]string{
			"focal":    "20.04",
			"jammy":    "22.04",
			"noble":    "24.04",
			"plucky":   "25.10",
			"resolute": "26.04",
		},
	},
	"microsoft-rpm": {
		Dest:      "/var/cache/packagesign/microsoft/yum",
		MirrorURL: "https://packages.microsoft.com/{dist_map}",
		Dists:     "7,8,9,10",
		Sections:  "prod",
		Format:    "rpm",
		SyncType:  "wget",
		Flat:      true,
		DistMap: map[string]string{
			"7":  "centos",
			"8":  "rhel",
			"9":  "rhel",
			"10": "rhel",
		},
	},
	"rocky": {
		Dest:          "/var/cache/packagesign/yum/rocky",
		MirrorURL:     "https://dl.rockylinux.org/pub/rocky",
		Sections:      "BaseOS,AppStream,extras",
		Architectures: "x86_64",
		Dists:         "8,9,10",
		Format:        "rpm",
		SyncType:      "wget",
		RsyncPath:     "{dist}/{section}/{arch}/os/",
		Exclude:       "isolinux,images,EFI",
	},
}
