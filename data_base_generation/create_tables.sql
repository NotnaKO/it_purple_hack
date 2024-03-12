DROP SCHEMA IF EXISTS matrix CASCADE;
CREATE SCHEMA matrix;
create table matrix.baseline_matrix_1_1(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm1 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_2(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm10 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_3(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm11 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_4(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm12 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_5(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm13 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_6(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm14 PRIMARY KEY (microcategory_id, location_id)
);


create table matrix.baseline_matrix_1_7(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm15 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_8(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm16 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_9(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm17 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_1_10(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm18 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_2(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm2 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.baseline_matrix_3(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm3 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.discount_matrix_1(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm4 PRIMARY KEY (microcategory_id, location_id)
);

create table matrix.discount_matrix_2(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm5 PRIMARY KEY (microcategory_id, location_id)
);

CREATE TABLE matrix.discount_matrix_3 (
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm6 PRIMARY KEY (microcategory_id, location_id)
);
