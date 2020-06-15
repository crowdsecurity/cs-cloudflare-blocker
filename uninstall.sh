#!/bin/bash

BIN_PATH_INSTALLED="/usr/local/bin/cloudfare-blocker"
CONFIG_DIR="/etc/crowdsec/cloudfare-blocker/"
PID_DIR="/var/run/crowdsec/"
SYSTEMD_PATH_FILE="/etc/systemd/system/cloudfare-blocker.service"

uninstall() {
	systemctl stop cloudfare-blocker
	rm -rf "${CONFIG_DIR}"
	rm -f "${SYSTEMD_PATH_FILE}"
	rm -f "${PID_DIR}cloudfare-blocker.pid"
	rm -f "${BIN_PATH_INSTALLED}"
}

uninstall