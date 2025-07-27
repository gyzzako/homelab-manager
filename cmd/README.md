# `cmd` Package Documentation

## Overview

The `cmd` package implements the CLI logic

## Host Command

### Description

The `host` command applies entries to the system's hosts file based on data provided via supported sources.

### Usage

```sh
homelab-manager host --provider <provider> --path <path> [--token <token>] [--push]
```

### Flags

| Flag        | Shorthand | Description                                    | Required |
|-------------|-----------|------------------------------------------------|----------|
| `--provider`| `-p`      | Specifies the data provider                    | Yes      |
| `--path`    |           | Path or URL to provider data                   | Yes      |
| `--token`   | `-t`      | Authentication token for URL provider          | No       |
| `--push`    |           | Push host data to remote Git repository        | No       |

### Examples

```sh
# Using config provider
homelab-manager host --provider config --path ./config.yaml
homelab-manager host -p config --path ./config.yaml

# Using config provider with Git push
homelab-manager host --provider config --path ./config.yaml --push

# Using URL provider
homelab-manager host --provider url --path https://api.example.com/hosts --token your-auth-token
homelab-manager host -p url --path https://api.example.com/hosts -t your-auth-token
```

### Supported Providers

- `config`: Uses a YAML file as the source of host entries.
  - Path must point to a valid YAML file (e.g., `./config.yaml`).
  - Supports Git push functionality with `--push` flag.
  - Example:
  ```yml
  host:
    - ip: 192.168.1.10
      domain: example.com
    - ip: 192.168.1.10
      domain: example.com
      subdomains:
        - sub1
        - sub2
    - ip: 192.168.1.11
      domain: mysite.local
      subdomains:
        - sub1

  git:
    url: https://github.com/user/repo.git
    token: github_token
  ```

- `url`: Fetches host data from a remote URL endpoint.
  - Path must be a valid URL (e.g., `https://api.example.com/hosts`).
  - Requires authentication token via `--token` flag.
  - Content structure:
  ```yml
    - ip: 192.168.1.10
      domain: example.com
    - ip: 192.168.1.10
      domain: example.com
      subdomains:
        - sub1
        - sub2
    - ip: 192.168.1.11
      domain: mysite.local
      subdomains:
        - sub1
  ```
