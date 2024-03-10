DROP SCHEMA IF EXISTS matrix CASCADE;
CREATE SCHEMA matrix;
create table matrix.baseline_matrix_1(
    microcategory_id int NOT NULL,
    location_id int NOT NULL,
    price int,
    CONSTRAINT nm1 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_2(
    microcategory_id int NOT NULL,
    location_id int NOT NULL,
    price int,
    CONSTRAINT nm2 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_3(
    microcategory_id int NOT NULL,
    location_id int NOT NULL,
    price int,
    CONSTRAINT nm3 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.discount_matrix_1(
    microcategory_id int NOT NULL,
    location_id int NOT NULL,
    price int,
    CONSTRAINT nm4 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.discount_matrix_2(
    microcategory_id int NOT NULL,
    location_id int NOT NULL,
    price int,
    CONSTRAINT nm5 PRIMARY KEY (microcategory_id, location_id)
);

CREATE TABLE matrix.discount_matrix_3 (
    microcategory_id int NOT NULL,
    location_id int NOT NULL,
    price int,
    CONSTRAINT nm6 PRIMARY KEY (microcategory_id, location_id)
);
