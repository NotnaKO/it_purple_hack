DROP SCHEMA IF EXISTS matrix CASCADE;
CREATE SCHEMA matrix;
create table matrix.baseline_matrix_1(
    microcategory_id bigint NOT NULL,
    location_id bigint NOT NULL,
    price float,
    CONSTRAINT nm1 PRIMARY KEY (microcategory_id, location_id)
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
