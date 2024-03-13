#!/bin/bash
psql -d postgres -a -f baseline_matrix_1.sql #Ubuntu only
psql -d postgres -a -f baseline_matrix_2.sql
psql -d postgres -a -f baseline_matrix_3.sql
psql -d postgres -a -f discount_matrix_1.sql
psql -d postgres -a -f discount_matrix_2.sql
psql -d postgres -a -f discount_matrix_3.sql