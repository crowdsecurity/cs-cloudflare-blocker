<p align="center">
<img src="https://github.com/crowdsecurity/cs-cloudflare-blocker/raw/master/docs/assets/crowdsec_cloudfare_logo.png" alt="CrowdSec" title="CrowdSec" width="280" height="400" />
</p>
<p align="center">
<img src="https://img.shields.io/badge/build-pass-green">
<img src="https://img.shields.io/badge/tests-pass-green">
</p>
<p align="center">
&#x1F4DA; <a href="https://docs.crowdsec.net/blockers/cloudflare/installation/">Documentation</a>
&#x1F4A0; <a href="https://hub.crowdsec.net">Hub</a>
&#128172; <a href="https://discourse.crowdsec.net">Discourse </a>
</p>

# CrowdSec Cloudflare Blocker

A blocker that will call Cloudflare's API when an IP is banned or unbanned.

# How does it work ?

cs-cloudflare-blocker will monitor MySQL or SQLite database and call Cloudflare's API to add/remove access rules for malevolent IPs.

# Installation

## Install script

Download the [latest release](https://github.com/crowdsecurity/cs-cloudflare-blocker/releases).

```bash
tar xzvf cs-cloudflare-blocker.tgz
cd cs-cloudflare-blocker/
sudo ./install.sh
systemctl status cloudflare-blocker
```


## From source

:warning: requires go >= 1.13

```bash
make release
cd cs-cloudflare-blocker-vX.X.X
sudo ./install.sh
systemctl status cloudflare-blocker
```

# Configuration

By default the blocker expects a SQLite backend, and the configuration file is as :

```yaml
# Cloudflare API information
api_key: <API_KEY>
email: <EMAIL_ADDR>
zone_id: <ZONE_ID>
account_id: <ACCOUNT_ID>
# Scope of the access rules
scope: <zone|account>
piddir: /var/run/
# How often the DB is polled for new bans
update_frequency: 30s
# Service-related options
daemonize: true
log_mode: file
log_dir: /var/log/
db_config:
  ## DB type supported (mysql, sqlite)
  ## By default it using sqlite
  type: sqlite

  ## mysql options
  # db_host: localhost
  # db_username: crowdsec
  # db_password: crowdsec
  # db_name: crowdsec

  ## sqlite options
  db_path: /var/lib/crowdsec/data/crowdsec.db
  flush: false

```

<details>
  <summary>MySQL configuration</summary>

```yaml
# Cloudflare API information
api_key: <API_KEY>
email: <EMAIL_ADDR>
zone_id: <ZONE_ID>
account_id: <ACCOUNT_ID>
# Scope of the access rules
scope: <zone|account>
piddir: /var/run/
# How often the DB is polled for new bans
update_frequency: 30s
# Service-related options
daemonize: true
log_mode: file
log_dir: /var/log/
db_config:
  ## DB type supported (mysql, sqlite)
  ## By default it using sqlite
  type: mysql

  ## mysql options
  db_host: localhost
  db_username: crowdsec
  db_password: crowdsec
  db_name: crowdsec

  ## sqlite options
  #db_path: /var/lib/crowdsec/data/crowdsec.db
  flush: false

```
</details>

# How it works

When the `cloudflare-blocker` service starts, it creates cloudflare access rules from new IPs in the SQLite database by using the cloudflare API.

# Troubleshooting

 - Logs are in `/var/log/cloudflare-blocker.log`
 - You can view/interact directly in the ban list either with `cscli`
 - Service can be started/stopped with `systemctl start/stop cloudflare-blocker`

