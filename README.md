# Kredentials
[![Version](https://img.shields.io/github/v/release/SticketInya/kredentials?include_prereleases&sort=semver)](https://github.com/SticketInya/kredentials/releases)
[![License](https://img.shields.io/github/license/SticketInya/kredentials)](https://github.com/SticketInya/kredentials/blob/main/LICENSE)

Kredentials is a simple CLI tool for managing Kubernetes config files. It makes it easy to backup, restore, and switch between different Kubernetes configurations, helping developers and operators streamline their cluster management workflow.

## Installation
- [Automatic install](#option-1-automatic-install)
- [Manual install](#option-2-manual-install)
- [Verifying the installation](#verify-the-installation)

### Option 1: Automatic install
To install `kredentials`, simply run:

```shell
    curl -sS https://raw.githubusercontent.com/SticketInya/kredentials/main/install.sh | bash
```

### Option 2: Manual install
1. Download the appropriate archive for your system from the [releases page](https://github.com/SticketInya/kredentials/releases/).
2. Extract the archive: `tar xzf kredentials_X.Y.Z_OS_ARCH.tar.gz`
3. Make the binary executable: `chmod +x kredentials`
4. Move to your PATH: `sudo mv kredentials /usr/local/bin/`

### Verify the installation
To verify that the installation was successful, run:
```shell
kredentials version
```
