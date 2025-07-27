# Homelab Manager

A CLI tool for homelab management that helps you manage your system's hosts file by applying entries from various data sources.

## Overview

Homelab Manager is a Go-based command-line application designed to simplify the management of your system's hosts file. It supports multiple data providers and can automatically update your hosts file with entries from configuration files or remote URLs.

## Features

- **Multiple Data Providers**: Support for YAML configuration files and URL-based data sources
- **Git Integration**: Push host data to remote Git repositories for backup and versioning
- **Authentication Support**: Token-based authentication for secure URL providers
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Installation

```bash
go build -o homelab-manager cmd/main.go
```
## Commands

### Host Command

The primary command for managing hosts file entries.

**📖 For complete documentation, usage examples, and all available parameters, see the [Host Command Documentation](cmd/README.md).**