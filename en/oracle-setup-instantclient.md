# Setup Oracle Instant Client

In this post, we are going to learn how to setup Oracle instant client on Linux, Windows, and MacOS.

---

### A. Table of Contents

 - [Setup Oracle Instant Client on Linux (Ubuntu 16.04)](#a-setup-oracle-instant-client-on-linux-ubuntu-1604)
 - [Setup Oracle Instant Client on Windows 10](#b-setup-oracle-instant-client-on-windows-10)

---

### B. Setup Oracle Instant Client on Linux (Ubuntu 16.04)

First of all, download these three files from https://www.oracle.com/technetwork/topics/linuxx86-64soft-092277.html.

- instantclient-basic-linux.x64-12.2.0.1.0.zip
- instantclient-sdk-linux.x64-12.2.0.1.0.zip
- instantclient-sqlplus-linux.x64-12.2.0.1.0.zip

Then update os package repository, continue with install the required essentials tools/

```bash
sudo apt-get update
sudo apt-get install build-essential libaio1 unzip git pkg-config
```

Next, we shall setup the oracle client.

```bash
mkdir -p /home/novalagung/oracle && cd /home/novalagung/oracle
cp /where/instantclient-*.zip .

unzip instantclient-basic-linux.x64-12.2.0.1.0.zip
unzip instantclient-sdk-linux.x64-12.2.0.1.0.zip
unzip instantclient-sqlplus-linux.x64-12.2.0.1.0.zip
rm -rf instantclient-*.zip

echo 'export ORACLE_HOME=/home/novalagung/oracle/instantclient_12_2' >> /home/novalagung/.bashrc
echo 'export PATH=$PATH:$ORACLE_HOME' >> /home/novalagung/.bashrc
source /home/novalagung/.bashrc

cd /home/novalagung/oracle/instantclient_12_2
ln -s libclntsh.so.12.1 libclntsh.so
ln -s libocci.so.12.1 libocci.so

sudo sh -c 'echo '/home/novalagung/oracle/instantclient_12_2' >> /etc/ld.so.conf.d/oracle-instantclient.conf'

echo 'export DYLD_LIBRARY_PATH=$ORACLE_HOME' >> /home/novalagung/.bashrc
echo 'export LD_LIBRARY_PATH=$ORACLE_HOME' >> /home/novalagung/.bashrc

sudo ldconfig
```

Next, create `oci8.pc` file. This file is required later by go oracle driver to be able to communicate with the oracle database server. If you plan only to connect to the oracle database server by using `sqlplus` only, then the file is not necessarily required.

```bash
sudo nano /usr/lib/pkgconfig/oci8.pc
```

Fill it with this content:

```bash
instantclient=/home/novalagung/oracle/instantclient_12_2
libdir=${instantclient}
includedir=${instantclient}/sdk/include/

Name: oci8
Description: oci8 library
Version: 12.1
Libs: -L${libdir} -lclntsh
Cflags: -I${includedir}
```

Last, try to connect to the oracle db server using `sqlplus`.

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

---

### C. Setup Oracle Instant Client on Windows 10

First of all, download these three files from https://www.oracle.com/technetwork/topics/winx64soft-089540.html.

- instantclient-basic-windows.x64-12.2.0.1.0.zip
- instantclient-sdk-windows.x64-12.2.0.1.0.zip
- instantclient-sqlplus-windows.x64-12.2.0.1.0.zip

Create `Oracle` folder at `C:\Oracle`, put all downloaded archives into this folder.

```bash
cd \
mkdir Oracle
```

Then extract zip files, all of them. By default it'll be extracted into `instantclient_12_2` folder under `C:\Oracle`, and let it be.

Append the `C:\Oracle\instantclient_12_2` path into `%PATH%` variable.

Next, set these CGO_ variables.

```bash
setx CGO_CFLAGS "-IC:\OtherPrograms\Oracle\instantclient_12_2\sdk\include"
setx CGO_LDFLAGS "-LC:\OtherPrograms\Oracle\instantclient_12_2 -loci"
```

Now we need to install **GCC**, and in this tutorial we'll use MSYS2 64bit. Download the installer `msys2-x86_64-*.exe` from https://www.msys2.org. After download process finished, run the installer. Pick any directory path you want, but make sure to remember it. In my place, I install it here.

```bash
C:\msys64
```

Next, run the **MSYS2 MinGW 64-bit** application. Then execute these commands.

```bash
# Update pacman
pacman -Su
# Install pkg-config and gcc
pacman -S mingw64/mingw-w64-x86_64-pkg-config mingw64/mingw-w64-x86_64-gcc
```

Now, set the `PKG_CONFIG_PATH` variable to points to the oci8.pc file inside msys64 (mingw64) pkgconfig folder.

```bash
setx PKG_CONFIG_PATH "C:\msys64\mingw64\lib\pkgconfig\oci8.pc"
```

Then add msys64 (mingw64) binary path into `%PATH%` variable.

```bash
C:\msys64\mingw64\bin
```

Next create the `oci8.pc` file on inside msys64 (mingw64) pkgconfig folder. This file is required later by go oracle driver to be able to communicate with the oracle database server. If you plan only to connect to the oracle database server by using `sqlplus` only, then the file is not necessarily required.

```bash
C:\msys64\mingw64\lib\pkgconfig\oci8.pc
```

Below is the content of the file.

```
oralib="C:/Oracle/instantclient_12_2/sdk/lib/msvc"
orainclude="C:/Oracle/instantclient_12_2/sdk/include"
gcclib="C:/msys64/mingw64/lib"
gccinclude="C:/msys64/mingw64/include"

Name: oci8
Version: 12.2
Description: oci8 library
Libs: -L${oralib} -L${gcclib} -loci
Libs.private:
Cflags: -I${orainclude} -I${gccinclude}
```

**REMINDER:** You need to adjust the `oralib`, `orainclude`, `gcclib`, and `gccinclude` value to match your settings. And also replace the backslash (`\`) into slash (`/`).

OK, the oracle client setup is done. Last step, try to connect to the oracle db server using `sqlplus`.

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
