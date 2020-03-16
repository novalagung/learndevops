# Terraform | CLI Installation

### 1. Download Installer

To install Terraform CLI, navigate to https://www.terraform.io/downloads.html, then download the installer. Pick the one that match your Operating System.

Then after that put the binary somewhere, and do add the binary path into `PATH` environment variable.

#### 1.1. Windows `PATH` Environment

Search on start menu using keyword `Environment Variables` → click `Environment Variables` → search `PATH` → click `Edit` → append the path where `terraform` cli is located.

![Terraform | CLI Installation | path environment variables](https://i.imgur.com/xU3fbTe.jpg)

Then open CMD/PowerShell, run `terraform -v` command. If a version number is appear, then everything is good.

![Terraform | CLI Installation | test terraform command](https://i.imgur.com/XOdec43.png)

#### 1.1 Linux/Unix `PATH` Environment

Download the binary, and then put it into `/usr/local/bin/`. Then run `terraform -v` command. If a version number is appear, then everything is good.

[Terraform | CLI Installation | download terraform linux](https://i.imgur.com/cuvt0hv.png)

Or you can also put it anywhere but make sure the directory path (where this terraform binary is placed) is added into `PATH` variable.
