package main

import (
	"containerd-installer/cmd/containerd-installer/install"
	"containerd-installer/pkg/hns"
	"containerd-installer/pkg/osutils"
	"context"
	"github.com/containerd/containerd/defaults"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/windows"
	"os"
)

func main() {
	app := &cli.App{
		Usage:           "Installs containerd on Windows, optionally with default CNI plugins",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug, d",
				Usage: "Run in debug mode",
			},
			&cli.StringFlag{
				Name:  "containerd-version",
				Usage: "Set containerd version to install",
				Value: "1.6.0-beta.3",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Set path where to install containerd to",
				Value: defaults.DefaultConfigDir,
			},
			&cli.StringFlag{
				Name:  "cni-plugin-version",
				Usage: "Set version of the CNI plugins to install",
				Value: "0.2.0",
			},
			&cli.BoolFlag{
				Name:  "no-cni-plugins",
				Usage: "Do not install CNI plugins",
			},
		},
		Before: func(c *cli.Context) error {
			if c.IsSet("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		Action: runCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func runCmd(c *cli.Context) error {
	containerdPath := c.String("path")
	group, _ := errgroup.WithContext(context.Background())

	isAdmin, err := osutils.IsAdmin()
	if err != nil {
		return errors.WithMessage(err, "Failed to check for administrator privileges")
	}
	if !isAdmin {
		return errors.New("You need to have administrator privileges to run installation")
	}

	restartNeeded, err := osutils.EnableFeatures([]string{"Containers", "Microsoft-Hyper-V", "Microsoft-Hyper-V-Management-PowerShell"})
	if err != nil {
		return errors.WithMessage(err, "Failed to enable required Windows features")
	}
	if restartNeeded {
		logrus.Info("Please restart to enable required Windows features and run installation again")
		os.Exit(int(windows.ERROR_FAIL_REBOOT_REQUIRED))
	}

	group.Go(func() error {
		if err := install.DownloadContainerd(containerdPath, c.String("containerd-version")); err != nil {
			return errors.WithMessage(err, "Failed to download and extract containerd")
		}

		if err := install.CreateContainerdConfig(containerdPath); err != nil {
			return errors.WithMessage(err, "Failed to verify downloaded containerd")
		}

		if err := install.RegisterAndStartContainerdService(containerdPath); err != nil {
			return errors.WithMessage(err, "Failed create register or start containerd service")
		}

		return nil
	})

	if !c.IsSet("no-cni-plugins") {
		group.Go(func() error {
			if err := install.DownloadCniPlugins(containerdPath, c.String("cni-plugin-version")); err != nil {
				return errors.WithMessage(err, "Failed to download CNI binaries")
			}

			natNetwork, err := hns.GetOrCreateNatNetwork()
			if err != nil {
				return errors.WithMessage(err, "Failed to create/retrieve NAT network")
			}

			if err := install.CreateNatCniConfig(containerdPath, natNetwork); err != nil {
				return errors.WithMessage(err, "Failed create CNI config for NAT network")
			}

			return nil
		})
	}

	defer logrus.Debug("Installation finished")

	return group.Wait()
}
