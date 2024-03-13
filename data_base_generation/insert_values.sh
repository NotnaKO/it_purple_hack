#!/bin/bash
sudo -u postgres psql -d postgres -a -f baseline_matrix_1.sql #Ubuntu only
sudo -u postgres psql -d postgres -a -f baseline_matrix_2.sql
sudo -u postgres psql -d postgres -a -f baseline_matrix_3.sql
# sudo -u postgres psql -d postgres -a -f discount_matrix_1.sql
# sudo -u postgres psql -d postgres -a -f discount_matrix_2.sql
# sudo -u postgres psql -d postgres -a -f discount_matrix_3.sql