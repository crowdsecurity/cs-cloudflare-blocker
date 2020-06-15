<p align="center">
<img src="https://github.com/crowdsecurity/cs-cloudfare-blocker/raw/master/docs/assets/crowdsec_cloudfare_logo.png" alt="CrowdSec" title="CrowdSec" width="280" height="400" />
</p>
<p align="center">
<img src="https://img.shields.io/badge/build-pass-green">
<img src="https://img.shields.io/badge/tests-pass-green">
</p>
<p align="center">
&#x1F4DA; <a href="https://docs.crowdsec.net/blockers/cloudfare/installation/">Documentation</a>
&#x1F4A0; <a href="https://hub.crowdsec.net">Hub</a>
&#128172; <a href="https://discourse.crowdsec.net">Discourse </a>
</p>

# CrowdSec Cloudfare Blocker

This repository contains a cloudfare-blocker, written in golang, that will bans IP address tagged as malevolent in the SQLite database by pushing access rules to Cloudfare API.

## Requirements

A cloudfare account is required for this blocker.
Please provide in the configuration the following information:
 - email used for cloudfare account
 - a valid API key
 - zone or account ID (depending on what you want to block IPs)


## Installation

Download the [latest release](https://github.com/crowdsecurity/cs-cloudfare-blocker/releases).

```bash
tar xzvf cs-cloudfare-blocker.tgz
cd cs-cloudfare-blocker/
sudo ./install.sh
```

## Documentation

Please find the documentation [here](https://docs.crowdsec.net/blockers/cloudfare/installation/).

# How it works

When the `cloudfare-blocker` service starts, it creates cloudfare access rules from new IPs in the SQLite database by using the cloudfare API.
:warning: the `dbpath` in the [blocker's configuration file](https://github.com/crowdsecurity/cs-cloudfare-blocker/blob/master/config/cloudfare-blocker.yaml#L2) must be consistent with the one used by crowdwatch.

# Troubleshooting

 - Logs are in `/var/log/cloudfare-blocker.log`
 - You can view/interact directly in the ban list either with `cscli` or direct at ipset level
 - Service can be started/stopped with `systemctl start/stop cloudfare-blocker`

