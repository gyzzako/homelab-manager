# `cmd` Package Documentation

## Overview

The `cmd` package implements the CLI logic

## Host Command

### Description

The `host` command applies entries to the system's hosts file based on data provided via supported sources.

### Usage

```sh
# Using CLI parameters
homelab-manager host --provider <provider> [--config <config>] [--path <path>] [--token <token>] [--type <type>] [--query <query>] [--push]

# Using config file
homelab-manager host --config <config>

# Using config file with CLI override
homelab-manager host --config <config> --provider <provider>
```

### Flags

| Flag        | Shorthand | Description                                    | Required |
|-------------|-----------|------------------------------------------------|----------|
| `--config`  | `-c`      | Config file path                               | No       |
| `--provider`| `-p`      | Specifies the data provider                    | No*      |
| `--path`    |           | Path/URL to data                               | No*      |
| `--token`   |           | Authentication token                           | No       |
| `--push`    |           | Push host data to remote Git repository        | No       |
| `--type`    |           | Database type (for SQL provider)               | No*      |
| `--query`   |           | SQL query (for SQL provider)                   | No*      |

*Required depending on provider type and whether using config file

### Examples

```sh
# Using CLI parameters
homelab-manager host --provider config --config ./config.yml
homelab-manager host --provider url --path https://api.example.com/hosts --token your-auth-token
homelab-manager host --provider sql --path ./sqlite.db --type sqlite --query "select * from hosts"

# Using config file
homelab-manager host --config ./config.yml

# Using config file with CLI override
homelab-manager host --config ./config.yml --provider url
```

### Supported Providers

- `config`: Uses a YAML configuration file as the source of host entries.
  - Requires `--config` flag to specify the configuration file path.

- `url`: Fetches host data from a remote URL endpoint.
  - Requires `--path` flag with a valid URL.
  - Requires `--token` flag for authentication.

- `sql`: Fetches host data from a SQL database.
  - Requires `--path` flag with the database file path or connection string.
  - Requires `--type` flag to specify the database type (e.g., sqlite).
  - Requires `--query` flag with the SQL query to execute.

### Configuration File Example

```yml
provider:
  type: sql
  #required for url provider
  url-params:
    url: https://example.com/host_entries.yml
    token: token
  #required for sql provider
  sql-params:
    datasource: ./sqlite.sqlite
    type: sqlite
    query:

#required for config provider
data:
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

#required for git integration
git:
  push: false
  url: https://github.com/user/repo.git
  token: token
```
