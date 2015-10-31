#!/bin/sh
BIN_DIR=/usr/bin
SCRIPT_DIR=/usr/lib/telegraf/scripts
LOG_DIR=/var/log/telegraf

if ! id telegraf >/dev/null 2>&1; then
    useradd --help 2>&1| grep -- --system > /dev/null 2>&1
    old_useradd=$?
    if [ $old_useradd -eq 0 ]
    then
	useradd --system -U -M telegraf
    else
	groupadd telegraf && useradd -M -g telegraf telegraf
    fi
fi

chown telegraf:telegraf $BIN_DIR/telegraf
chmod a+rX $BIN_DIR/telegraf

# Systemd
if which systemctl > /dev/null 2>&1 ; then
    cp -f $SCRIPT_DIR/telegraf.service /lib/systemd/system/telegraf.service
    systemctl enable telegraf
# Sysv
else
    # Remove legacy symlink
    test -h /etc/init.d/telegraf && rm -f /etc/init.d/telegraf

    cp -f $SCRIPT_DIR/init.sh /etc/init.d/telegraf
    chmod +x /etc/init.d/telegraf
    # update-rc.d sysv service:
    if which update-rc.d > /dev/null 2>&1 ; then
	update-rc.d -f telegraf remove
	update-rc.d telegraf defaults
	# CentOS-style sysv:
    else
	chkconfig --add telegraf
    fi
    mkdir -p $LOG_DIR
    chown -R -L telegraf:telegraf $LOG_DIR
fi
