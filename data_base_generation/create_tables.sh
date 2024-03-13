#!/bin/bash
sudo -u postgres psql -d postgres -a -f create_tables.sql  # Ubuntu only