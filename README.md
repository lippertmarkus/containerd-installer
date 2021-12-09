# containerd-installer

Installs containerd on Windows, optionally with default CNI plugins

## Usage
```   
NAME:
   containerd-installer.exe - Install containerd on Windows, optionally with default CNI plugins

USAGE:
   containerd-installer.exe [global options] [arguments...]

GLOBAL OPTIONS:
   --debug                     Run in debug mode (default: false)
   --containerd-version value  Set containerd version to install (default: "1.6.0-beta.3")
   --path value                Set path where to install containerd to (default: "C:\\Program Files\\containerd")
   --cni-plugin-version value  Set version of the CNI plugins to install (default: "0.2.0")
   --no-cni-plugins            Do not install CNI plugins (default: false)
   --help, -h                  show help (default: false)

```

## Run
```powershell
PS C:\> .\containerd-installer.exe --debug

time="2021-12-09T14:59:05+01:00" level=debug msg="Checking for admin privileges"
time="2021-12-09T14:59:05+01:00" level=debug msg="Enabling Windows features: Containers, Microsoft-Hyper-V, Microsoft-Hyper-V-Management-PowerShell"
time="2021-12-09T14:59:07+01:00" level=debug msg="Downloading from https://github.com/containerd/containerd/releases/download/v1.6.0-beta.3/containerd-1.6.0-beta.3-windows-amd64.tar.gz"
time="2021-12-09T14:59:07+01:00" level=debug msg="Downloading from https://github.com/microsoft/windows-container-networking/releases/download/v0.2.0/windows-container-networking-cni-amd64-v0.2.0.zip"
time="2021-12-09T14:59:09+01:00" level=debug msg="Check for existing HNS NAT network"
time="2021-12-09T14:59:09+01:00" level=debug msg="Creating CNI config for NAT network"
time="2021-12-09T14:59:10+01:00" level=debug msg="Creating containerd config"
time="2021-12-09T14:59:11+01:00" level=debug msg="Creating containerd service and starting it"
time="2021-12-09T14:59:12+01:00" level=debug msg="Installation finished"
```