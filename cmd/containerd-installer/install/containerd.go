package install

import (
	"containerd-installer/pkg/http"
	"fmt"
	"github.com/google/glazier/go/helpers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
)

const ContainerdUrl = "https://github.com/containerd/containerd/releases/download/v%[1]s/containerd-%[1]s-windows-amd64.tar.gz"
const ContainerdServiceName = "containerd"

func DownloadContainerd(targetPath string, version string) error {
	url := fmt.Sprintf(ContainerdUrl, version)
	if err := http.DownloadAndExtract(url, targetPath, 1); err != nil {
		return errors.WithMessage(err, "Error downloading containerd")
	}

	return nil
}

func CreateContainerdConfig(containerdPath string) error {
	logrus.Debug("Creating containerd config")

	configFilePath := filepath.Join(containerdPath, "config.toml")
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return errors.Wrapf(err, "Couldn't create config file %s", configFilePath)
	}
	defer configFile.Close()

	binaryPath := filepath.Join(containerdPath, "containerd.exe")
	cmd := exec.Command(binaryPath, "config", "default")
	cmd.Stdout = configFile
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "Couldn't run containerd binary %s", binaryPath)
	}

	return nil
}

func RegisterAndStartContainerdService(containerdPath string) error {
	logrus.Debug("Creating containerd service and starting it")

	binaryPath := filepath.Join(containerdPath, "containerd.exe")
	if err := exec.Command(binaryPath, "--register-service").Run(); err != nil {
		return errors.Wrapf(err, "Failed to run containerd command to register service")
	}

	if err := helpers.StartService(ContainerdServiceName); err != nil {
		return errors.Wrapf(err, "Couldn't start service %s", ContainerdServiceName)
	}

	return nil
}
