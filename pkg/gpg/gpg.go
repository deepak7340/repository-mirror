package gpg

import (
	"fmt"
	"os"
	"os/exec"
)

func keyringHome() string {
	if h := os.Getenv("GNUPGHOME"); h != "" {
		return h
	}
	return "/var/cache/packagesign/keyrings"
}

func InitKeyring() {
	home := keyringHome()
	if err := os.MkdirAll(home, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not create keyring dir %s: %v\n", home, err)
		return
	}
	kbx := home + "/trustedkeys.kbx"
	gpgFile := home + "/trustedkeys.gpg"
	if _, err := os.Stat(gpgFile); os.IsNotExist(err) {
		if _, err := os.Stat("/usr/share/keyrings/ubuntu-archive-keyring.gpg"); err == nil {
			importIntoKeyring("/usr/share/keyrings/ubuntu-archive-keyring.gpg")
		}
		os.Remove(kbx)
	}
}

func importIntoKeyring(keyFile string) {
	home := keyringHome()
	cmd := exec.Command("gpg", "--homedir", home, "--no-default-keyring", "--keyring", "trustedkeys.gpg", "--import", keyFile)
	cmd.Run()
}

func fixGpgKeyring(home string) {
	exec.Command("gpg-connect-agent", "--homedir", home, "killagent", "/bye").Run()
	os.Remove(home + "/S.gpg-agent")
	os.Remove(home + "/S.gpg-agent.browser")
	os.Remove(home + "/S.gpg-agent.extra")
	os.Remove(home + "/S.gpg-agent.ssh")
	os.Remove(home + "/pubring.kbx")
}

func ImportGPGKeys(keyURL string, keyIDs []string, keyring string) {
	InitKeyring()

	if keyURL != "" {
		tmpKey, err := os.CreateTemp("", "gpgkey-*.asc")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create temp file for GPG key: %v\n", err)
		} else {
			tmpKey.Close()
			defer os.Remove(tmpKey.Name())

			dl := exec.Command("wget", "-q", "-O", tmpKey.Name(), keyURL)
			if err := dl.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not download GPG key from %s: %v\n", keyURL, err)
			} else {
				importIntoKeyring(tmpKey.Name())
			}
		}
	}

	home := keyringHome()
	for _, keyID := range keyIDs {
		fetch := exec.Command("gpg", "--homedir", home, "--no-default-keyring", "--keyring", "trustedkeys.gpg", "--keyserver", "keyserver.ubuntu.com", "--recv-keys", keyID)
		output, err := fetch.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch GPG key %s: %v\n", keyID, string(output))
		}
	}

	fixGpgKeyring(home)
}
