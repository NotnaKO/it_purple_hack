#!/bin/bash
pg_dump --inserts --data-only --load-via-partition-root -t "$1" -d "$2" > "$3"