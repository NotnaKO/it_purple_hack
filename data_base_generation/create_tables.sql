DROP SCHEMA IF EXISTS matrix CASCADE;
CREATE SCHEMA matrix;
create table matrix.baseline_matrix_1(
    microcategory_id int,
    location_id int,
    price int
);

create table matrix.baseline_matrix_2(
    microcategory_id int,
    location_id int,
    price int
);

create table matrix.baseline_matrix_3(
    microcategory_id int,
    location_id int,
    price int
);

create table matrix.discount_matrix_1(
    microcategory_id int,
    location_id int,
    price int
);

create table matrix.discount_matrix_2(
    microcategory_id int,
    location_id int,
    price int
);

create table matrix.discount_matrix_3(
     microcategory_id int,
     location_id int,
     price int
);
