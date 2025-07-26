# `cmd` Package Documentation

## Overview

The `cmd` package implements the CLI logic

## Host Command

### Description

The `host` command applies entries to the system's `hosts` file based on data provided via supported sources.

### Usage

```sh
homelab-manager host --provider <provider> --path <path>
```

### Flags

| Flag        | Shorthand | Description                    | Required |
|-------------|-----------|--------------------------------|----------|
| `--provider`| `-p`      | Specifies the data provider    | Yes      |
| `--path`    |           | Path or URL to provider data   | Yes      |

### Example

```sh
homelab-manager host --provider config --path ./config.yaml
homelab-manager host -p config --path ./config.yaml
```

### Supported Providers

- `config`: Uses a YAML file as the source of host entries.
  - Path must point to a valid YAML file (e.g., `./config.yaml`).
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
  ```

