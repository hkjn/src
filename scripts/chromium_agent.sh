#!/bin/bash

# Start Chromium configured to pretend it's another user agent.
# USER_AGENT="Mozilla/5.0 (X11; CrOS x86_64 5116.115.5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.152 Safari/537.36"
# USER_AGENT="Windows / Firefox 26: Mozilla/5.0 (Windows NT 6.1; WOW64; rv:26.0) Gecko/20100101 Firefox/26.0"
USER_AGENT="Windows / Chrome 32: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.107 Safari/537.36"
chromium --user-agent="$USER_AGENT"
