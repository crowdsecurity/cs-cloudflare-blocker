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

This repository contains a cloudflare-blocker, written in golang, that will bans IP address tagged as malevolent in the SQLite database by pushing access rules to Cloudflare API.

## Requirements

A cloudflare account is required for this blocker.
Please provide in the configuration the following information:
 - email used for cloudflare account
 - a valid API key : you can create it in your cloudflare dashboard in : "My Profile" => "Api token"
 - zone or account ID (depending on what you want to block IPs)


## Installation

Download the [latest release](https://github.com/crowdsecurity/cs-cloudflare-blocker/releases).

```bash
tar xzvf cs-cloudflare-blocker.tgz
cd cs-cloudflare-blocker/
sudo ./install.sh
```

## Documentation

Please find the documentation [here](https://docs.crowdsec.net/blockers/cloudflare/installation/).


### Configuration

The configuration file (located under `/etc/crowdsec/cloudflare-blocker/cloudflare-blocker.yaml`) support those options:

```yaml
api_key: <API_KEY>                             # your cloudflare api key
email: <EMAIL_ADDR>                            # your cloudflare email address
scope: <account|zone>                          # the cloudflare access rule scope : account or zone
zone_id: <ZONE_ID>                             # your cloudflare zone ID if if the selected scope is "zone"
account_id: <ACCOUNT_ID>                       # your cloudflare account ID if the selected scope is "account
piddir: /var/run/
update_frequency: 30s
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

```

# How it works

When the `cloudflare-blocker` service starts, it creates cloudflare access rules from new IPs in the SQLite database by using the cloudflare API.

:warning: the `db_config` in the [blocker's configuration file](https://github.com/crowdsecurity/cs-cloudflare-blocker/blob/master/config/cloudflare-blocker.yaml#L6) must be consistent with the one used by crowdsec.

# Troubleshooting

 - Logs are in `/var/log/cloudflare-blocker.log`
 - You can view/interact directly in the ban list either with `cscli` or direct at ipset level
 - Service can be started/stopped with `systemctl start/stop cloudflare-blocker`

