# Kubernetes | Kubectl Installation

Kubernetes CLI or `kubectl` is a command line tool that used to run operation againts the kubernetes clusters, including inspecting and managing the cluster resources, view logs, etc.

### 1. Download Installer

To install `kubectl`, navigate to https://kubernetes.io/docs/tasks/tools/install-kubectl/, then download the installer. Pick the one that match your Operating System.

Then after that put the binary somewhere locally, and do add the binary path into `PATH` environment variable.

#### 1.1. Windows `PATH` Environment

Search on start menu using keyword `Environment Variables` → click `Environment Variables` → search `PATH` → click `Edit` → append the path where `terraform` cli is located.

![Kubernetes | Kubectl Installation | path environment variables](https://i.imgur.com/xU3fbTe.jpg)

Then open CMD/PowerShell, run `kubectl version --client` command. If a version number of kubernetes client is appear, then everything is good.

![Kubernetes | Kubectl Installation | download kubectl linux](https://i.imgur.com/YkNJ31p.png)

#### 1.2. Linux/Unix `PATH` Environment

Download the binary, and then put it into `/usr/local/bin/`. Then run `kubectl version --client` command. If a version number of kubernetes client is appear, then everything is good.

![Kubernetes | Kubectl Installation | download kubectl linux](https://i.imgur.com/HGxoUBM.png)

Or you can also put it anywhere but make sure the directory path (where this terraform binary is placed) is added into `PATH` variable.
