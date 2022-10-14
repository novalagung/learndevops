# Setup Oracle XE Database Server

In this post, we are going to learn how to setup Oracle XE Database Server on CentOS, Oracle Linux, and using Docker Container.

---

## Table of Contents

 - [Setup Oracle XE Database Server on CentOS 6 (Oracle Linux)](#a-setup-oracle-xe-database-server-on-centos-6-oracle-linux)
 - [Setup Oracle XE Database Server using Docker](#b-setup-oracle-xe-database-server-using-docker)

---

## A. Setup Oracle XE Database Server on CentOS 6 (Oracle Linux)

### A.1. Convert CentOS 6 into Oracle Linux

The easiest way to install Oracle Database Server is through **Oracle Linux** distribution.

Oracle Linux is a Linux distribution packaged and freely distributed by Oracle, available partially under the GNU General Public License since late 2006. It's free, we can easily get it from [Oracle Linux Download Page](https://www.oracle.com/linux/).

There is also an alternative way to get the Oracle Linux, by converting CentOS into Oracle Linux. In this post we'll learn to do that.

OK let's start. First of all, update os package repository.

```bash
sudo yum update
```

Oracle provides us capability to convert CentOS into Oracle Linux, and they make it to be so easy to use. For detailed information just take a look at https://linux.oracle.com/switch/centos.

Ok, let's download the `centos2ol.sh` file then execute it.

```bash
curl -O https://linux.oracle.com/switch/centos2ol.sh 
sudo sh centos2ol.sh

# Checking for required packages...
# Checking your distribution...
# Checking for yum lock...
# Looking for yumdownloader...
# Finding your repository directory...
# Downloading Oracle Linux yum repository file...
# Backing up and removing old repository files...
# Downloading Oracle Linux release package...
# 
# ... will take sometime
#
# Dependency Updated:
#   plymouth-core-libs.x86_64 0:0.8.3-29.0.1.el6
# 
# Replaced:
#   redhat-logos.noarch 0:60.0.14-12.el6.centos
# 
# Finished Transaction
# > Leaving Shell
# Updating initrd...
# Installation successful!
# Run 'yum distro-sync' to synchronize your installed packages
# with the Oracle Linux repository.
```

Next, synchronize the installed packages to the Oracle Linux repository by using command below.

```bash
sudo yum distro-sync

# Loaded plugins: fastestmirror, security
# Setting up Distribution Synchronization Process
# Loading mirror speeds from cached hostfile
# Only Upgrade available on package: sysstat-9.0.4-33.el6_9.1.x86_64
# Resolving Dependencies
# --> Running transaction check
# ---> Package acpid.x86_64 0:1.0.10-3.el6 will be updated
# ---> Package acpid.x86_64 0:2.0.19-6.0.1.el6 will be an update
# 
# ...
# 
# Updated:
#   sos.noarch 0:3.2-63.0.1.el6_10.2
#   system-config-network-tui.noarch 0:1.6.0.el6.3-4.0.1.el6
#   systemtap-runtime.x86_64 0:2.9-9.0.1.el6
#   yum-plugin-security.noarch 0:1.1.30-42.0.1.el6_10
#   yum-utils.noarch 0:1.1.30-42.0.1.el6_10
# 
# Complete!
```

Just that, your Oracle Linux is ready.

### A.2. Setup Oracle XE Database Server on Oracle Linux

You can get Oracle linux from [Oracle Linux download page](https://www.oracle.com/linux/), or by [converting CentOS into Oracle Linux](/convert-linux-centos-into-oracle-linux.md).

Download the **Oracle Database Express Edition 11g Release 2 for Linux x64** from https://www.oracle.com/technetwork/database/database-technologies/express-edition/downloads/xe-prior-releases-5172097.html. You might need to download it from the web browser since the download process require us to log in using oracle account (create one on the website if you haven't).

Unzip the downloaded oracle xe installer.

```bash
mkdir -p /home/novalagung/oracle-xe
cd /home/novalagung/oracle-xe
cp /path/to/file/oracle-xe-11.2.0-1.0.x86_64.rpm.zip .
unzip oracle-xe-11.2.0-1.0.x86_64.rpm.zip

sudo rpm -ivh Disk1/oracle-xe-11.2.0-1.0.x86_64.rpm
```

Now the Oracle XE 11g is installed. Next we need to run the Oracle XE Configuration. In this step few prompts will appear asking certain information like port for Oracle Application Express and for database listener.

One default user will be created, it's `SYSTEM` user. We'll need to put some password for this user (cannot left it empty). In this example we use `MANAGER` as the password.

```bash
sudo /etc/init.d/oracle-xe configure

# Oracle Database 11g Express Edition Configuration
# -------------------------------------------------
# This will configure on-boot properties of Oracle Database 11g Express 
# Edition.  The following questions will determine whether the database should 
# be starting upon system boot, the ports it will use, and the passwords that 
# will be used for database accounts.  Press <Enter> to accept the defaults. 
# Ctrl-C will abort.
# 
# Specify the HTTP port that will be used for Oracle Application Express [8080]:
# 
# Specify a port that will be used for the database listener [1521]:
# 
# Specify a password to be used for database accounts.  Note that the same
# password will be used for SYS and SYSTEM.  Oracle recommends the use of 
# different passwords for each database account.  This can be done after 
# initial configuration: MANAGER
# Confirm the password: MANAGER
# 
# Do you want Oracle Database 11g Express Edition to be started on boot (y/n) [y]:y
# 
# Starting Oracle Net Listener...Done
# Configuring database...Done
# Starting Oracle Database 11g Express Edition instance...Done
# Installation completed successfully.
```

Ok, now our Oracle XE is 100% ready. Next, we need to perform some tests, to make sure everything is working fine. We'll try to connect to the database server using default user `SYSTEM` and password `MANAGER`.

```bash
sqlplus SYSTEM/MANAGER@//localhost:1521/xe

# SQL*Plus: Release 11.2.0.4.0 Production on Wed Oct 3 06:48:54 2018
# Copyright (c) 1982, 2013, Oracle.  All rights reserved.
# 
# Connected to:
# Oracle Database 11g Express Edition Release 11.2.0.2.0 - 64bit Production
```

The result is: **connected**. Try to perform a simple query like getting the database version.

```bash
SQL> SELECT * FROM V$VERSION;

BANNER
--------------------------------------------------------------------------------
Oracle Database 11g Express Edition Release 11.2.0.2.0 - 64bit Production
PL/SQL Release 11.2.0.2.0 - Production
CORE	11.2.0.2.0	Production
TNS for Linux: Version 11.2.0.2.0 - Production
NLSRTL Version 11.2.0.2.0 - Production
```

---

## B. Setup Oracle XE Database Server using Docker

> This tutorial can be implemented in both Windows, Linux, or MacOS operating systems.

Download the **Oracle Database Express Edition 11g Release 2 for Linux x64** from https://www.oracle.com/technetwork/database/database-technologies/express-edition/downloads/xe-prior-releases-5172097.html. You might need to download it from the web browser since the download process require us to log in using oracle account (create one on the website if you haven't).

> REMINDER: Even you perform this installation on Windows or MacOS, you must download the Linux x64 installer! not the windows version or the macos version.

Then clone the official oracle docker images from their github.

```bash
git clone https://github.com/oracle/docker-images.git
```

Move the downloaded oracle xe installer into this path.

```bash
cd docker-images/OracleDatabase/SingleInstance/dockerfiles
cp oracle-xe-11.2.0-1.0.x86_64.rpm.zip /11.2.0.2/
```

Next, execute the `./buildDockerImage.sh` command with several arguments:

- Flag `-v 11.2.0.2` to specify the oracle version (in this case it's 11.2.0.2). The choosen version must match with the installer version.
- Flag `-x` to pick the **Express Edition** image.
- Flag `-i` to skip the md5sum verification.

```bash
./buildDockerImage.sh -v 11.2.0.2 -x -i

# Ignored MD5 checksum.
# ==========================
# DOCKER info:
# Containers: 3
#  Running: 0
#  Paused: 0
#  Stopped: 3
# Images: 10
# Server Version: 18.09.0
# ...
# 
# ==========================
# Building image 'oracle/database:11.2.0.2-xe' ...
# Sending build context to Docker daemon  631.8MB
# Step 1/10 : FROM oraclelinux:7-slim
# 7-slim: Pulling from library/oraclelinux
# a8d84c1f755a: Pulling fs layer
# a8d84c1f755a: Verifying Checksum
# a8d84c1f755a: Download complete
# ...
# 
# Removing intermediate container 51a3bdde4e7e
#  ---> bf56ef57fe4c
# Step 9/10 : HEALTHCHECK --interval=1m --start-period=5m    CMD "$ORACLE_BASE/$CHECK_DB_FILE" >/dev/null || exit 1
#  ---> Running in dcee11bca78e
# Removing intermediate container dcee11bca78e
#  ---> 4fbcb8aec67f
# Step 10/10 : CMD exec $ORACLE_BASE/$RUN_FILE
#  ---> Running in 253bd5706098
# Removing intermediate container 253bd5706098
#  ---> 97fb5f2328d0
# [Warning] One or more build-args [DB_EDITION] were not consumed
# Successfully built 97fb5f2328d0
# Successfully tagged oracle/database:11.2.0.2-xe
# SECURITY WARNING: You are building a Docker image from Windows against a non-Windows Docker host. All files and directories added to build context will have '-rwxr-xr-x' permissions. It is recommended to double check and reset permissions for sensitive files and directories.
# 
#   Oracle Database Docker Image for 'xe' version 11.2.0.2 is ready to be extended:
# 
#     --> oracle/database:11.2.0.2-xe
# 
#   Build completed in 303 seconds.
```

The process will take some time. In the end a new docker image called `oracle/database` will be created.

Next, start a new container using the `oracle/database` image.

```bash
docker run --name my-oracle-db-server \
    -p 1521:1521 \
    -p 5500:5500 \
    -e ORACLE_SID=xe \
    -e ORACLE_PWD=MANAGER \
    -v oradata:/opt/oracle/oradata \
    --shm-size=2g \
    oracle/database:11.2.0.2-xe

# ORACLE PASSWORD FOR SYS AND SYSTEM: MANAGER
# 
# Oracle Database 11g Express Edition Configuration
# -------------------------------------------------
# This will configure on-boot properties of Oracle Database 11g Express
# Edition.  The following questions will determine whether the database should
# be starting upon system boot, the ports it will use, and the passwords that
# will be used for database accounts.  Press <Enter> to accept the defaults.
# Ctrl-C will abort.
# 
# Specify the HTTP port that will be used for Oracle Application Express [8080]:
# Specify a port that will be used for the database listener [1521]:
# Specify a password to be used for database accounts.  Note that the same
# password will be used for SYS and SYSTEM.  Oracle recommends the use of
# different passwords for each database account.  This can be done after
# initial configuration:
# Confirm the password:
# 
# Do you want Oracle Database 11g Express Edition to be started on boot (y/n) [y]:
# Starting Oracle Net Listener...Done
# Configuring database...
# 
# ...
# 
# #########################
# DATABASE IS READY TO USE!
# #########################
# The following output is now a tail of the alert.log:
# QMNC started with pid=24, OS id=685
# Completed: ALTER DATABASE OPEN
# Fri Feb 22 08:17:28 2019
# db_recovery_file_dest_size of 10240 MB is 0.98% used. This is a
# user-specified limit on the amount of space that will be used by this
# database for recovery-related files, and does not reflect the amount of
# space available in the underlying filesystem or ASM diskgroup.
# Starting background process CJQ0
# Fri Feb 22 08:17:28 2019
# CJQ0 started with pid=25, OS id=699
```

Few explanations about above command arguments:

- Flag `-p 1521:1521`, export the oracle listener port.
- Flag `-p 5500:5500`, export the oem express port.
- Flag `-e ORACLE_SID=xe`, specify the oracle SID.
- Flag `-e ORACLE_PWD=MANAGER`, set the default password of `SYS`, `SYSTEM` and `PDB_ADMIN` users.
- Flag `-v oradata:/opt/oracle/oradata`, mirror the volume.
- Flag `--shm-size=2g`, allocate memory size for particular container.

Ok, now lets try to connect to the database server using default user `SYSTEM`.

```bash
sqlplus SYSTEM/MANAGER@//localhost:1521/XE

# SQL*Plus: Release 11.2.0.4.0 Production on Wed Oct 3 06:48:54 2018
# Copyright (c) 1982, 2013, Oracle.  All rights reserved.
# 
# Connected to:
# Oracle Database 11g Express Edition Release 11.2.0.2.0 - 64bit Production
```

The result is: **connected**. Try to perform a simple query like getting the database version.

```bash
SQL> SELECT * FROM V$VERSION;

BANNER
--------------------------------------------------------------------------------
Oracle Database 11g Express Edition Release 11.2.0.2.0 - 64bit Production
PL/SQL Release 11.2.0.2.0 - Production
CORE	11.2.0.2.0	Production
TNS for Linux: Version 11.2.0.2.0 - Production
NLSRTL Version 11.2.0.2.0 - Production
```

If you the container stopped, then you just need to start it. No need to create new container using same specification.
