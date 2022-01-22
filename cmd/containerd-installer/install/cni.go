package install

import (
	"containerd-installer/pkg/hns"
	"containerd-installer/pkg/http"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

const CniPluginUrl = "https://github.com/microsoft/windows-container-networking/releases/download/v%[1]s/windows-container-networking-cni-amd64-v%[1]s.zip"

const NatCniConfig = `{
    "cniVersion": "0.2.0",
    "name": "nat",
    "type": "nat",
    "master": "Ethernet",
    "ipam": {
        "subnet": "%s",
        "routes": [
            {
                "gateway": "%s"
            }
        ]
    },
    "capabilities": {
        "portMappings": true,
        "dns": true
    }
}`

func DownloadCniPlugins(containerdPath string, version string) error {
	cniBinDir := filepath.Join(containerdPath, "cni", "bin")
	if err := os.MkdirAll(cniBinDir, os.ModeDir); err != nil {
		return errors.Wrapf(err, "Couldn't create directory %s", cniBinDir)
	}

	url := fmt.Sprintf(CniPluginUrl, version)
	if err := http.DownloadAndExtract(url, cniBinDir, 0); err != nil {
		return errors.WithMessage(err, "Error downloading CNI binaries")
	}

	return nil
}

func CreateNatCniConfig(containerdPath string, natNetwork hns.NatNetworkDetails) error {
	logrus.Debug("Creating CNI config for NAT network")

	cniConfDir := filepath.Join(containerdPath, "cni", "conf")
	if err := os.MkdirAll(cniConfDir, os.ModeDir); err != nil {
		return errors.Wrapf(err, "Couldn't create directory %s", cniConfDir)
	}

	cniConfigPath := filepath.Join(cniConfDir, "0-containerd-nat.conf")
	if err := ioutil.WriteFile(cniConfigPath, []byte(fmt.Sprintf(NatCniConfig, natNetwork.AddressPrefix, natNetwork.GatewayAddress)), 0644); err != nil {
		return errors.Wrapf(err, "Couldn't write CNI config for NAT to %s", cniConfigPath)
	}

	return nil
}

func CreateCtrSymlinks(containerdPath string) error {
	logrus.Debug("Creating symlinks to CNI binaries and configs for ctr to work")

	if err := os.MkdirAll("/etc/cni/", os.ModeDir); err != nil {
		return errors.Wrapf(err, "Couldn't create directory /etc/cni")
	}
	if err := os.MkdirAll("/opt/cni/", os.ModeDir); err != nil {
		return errors.Wrapf(err, "Couldn't create directory /opt/cni/")
	}

	cniConfDir := filepath.Join(containerdPath, "cni", "conf")
	cniBinDir := filepath.Join(containerdPath, "cni", "bin")

	if err := os.Symlink(cniConfDir, "/etc/cni/net.d"); err != nil {
		return errors.Wrapf(err, "Couldn't create symlink /etc/cni/net.d")
	}
	if err := os.Symlink(cniBinDir, "/opt/cni/bin"); err != nil {
		return errors.Wrapf(err, "Couldn't create symlink /opt/cni/bin")
	}

	return nil
}
