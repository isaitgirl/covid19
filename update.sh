#!/usr/bin/env bash
cp resources/db.json resources/db.json.$(date +"%Y%m%d_%H%M") && \
wget $(cat resources/db.url) -O resources/db.json && \
echo "Database updated"
