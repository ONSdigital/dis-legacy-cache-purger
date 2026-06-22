#!/bin/sh

# Capture environment variables for cron to use
env > /etc/environment

# start cron
echo "Executing cron on schedule: $(crontab -l 2>/dev/null || echo 'no crontab installed')"
exec cron -f
    