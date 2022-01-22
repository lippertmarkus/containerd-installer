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
   --containerd-version value  Set containerd version to install (default: "1.6.0-rc.1")
   --path value                Set path where to install containerd to (default: "C:\\Program Files\\containerd")
   --cni-plugin-version value  Set version of the CNI plugins to install (default: "0.2.0")
   --no-cni-plugins            Do not install CNI plugins (default: false)
   --no-ctr-symlinks           Disable creating symlinks in /etc/cni/net.d and /opt/cni/bin (currently required for using ctr) (default: false)
   --help, -h                  show help (default: false)

```

## Run
```powershell
PS C:\> .\containerd-installer.exe --debug

time="2021-12-09T14:08:55Z" level=debug msg="Checking for admin privileges"
time="2021-12-09T14:08:55Z" level=debug msg="Enabling Windows features: Containers, Microsoft-Hyper-V, Microsoft-Hyper-V-Management-PowerShell"
time="2021-12-09T14:08:58Z" level=debug msg="Downloading from https://github.com/containerd/containerd/releases/download/v1.6.0-rc.1/containerd-1.6.0-rc.1-windows-amd64.tar.gz"
time="2021-12-09T14:08:58Z" level=debug msg="Downloading from https://github.com/microsoft/windows-container-networking/releases/download/v0.2.0/windows-container-networking-cni-amd64-v0.2.0.zip"
time="2021-12-09T14:08:59Z" level=debug msg="Check for existing HNS NAT network"
time="2021-12-09T14:08:59Z" level=debug msg="No existing NAT network can be retrieved, creating a new one"
time="2021-12-09T14:09:00Z" level=debug msg="Created NAT network: gateway: 172.22.16.1, subnet: 172.22.16.0/20"
time="2021-12-09T14:09:00Z" level=debug msg="Creating CNI config for NAT network"
time="2021-12-09T14:09:00Z" level=debug msg="Creating symlinks to CNI binaries and configs for ctr to work"
time="2021-12-09T14:09:00Z" level=debug msg="Creating containerd config"
time="2021-12-09T14:09:00Z" level=debug msg="Creating containerd service and starting it"
time="2021-12-09T14:09:01Z" level=debug msg="Installation finished"
```
