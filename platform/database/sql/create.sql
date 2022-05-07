-- DROP TABLE IF EXISTS "public"."revenues";
-- DROP TABLE IF EXISTS "public"."flats";
-- DROP TABLE IF EXISTS "public"."user_building_membership";
-- DROP TABLE IF EXISTS "public"."blocks";
-- DROP TABLE IF EXISTS "public"."buildings";
-- DROP TABLE IF EXISTS "public"."districts";
-- DROP TABLE IF EXISTS "public"."cities";
-- DROP TABLE IF EXISTS "public"."users";

CREATE TABLE IF NOT EXISTS "public"."cities"
(
    "id"        SERIAL PRIMARY KEY,
    "city_name" varchar(255)
);

CREATE TABLE IF NOT EXISTS "public"."districts"
(
    "id"            SERIAL PRIMARY KEY,
    "city_id"       int4         NOT NULL,
    "district_name" varchar(255) NOT NULL,
    CONSTRAINT fk_city FOREIGN KEY (city_id) REFERENCES cities (id)
);

CREATE TABLE IF NOT EXISTS "public"."users"
(
    "id"           SERIAL PRIMARY KEY,
    "created_by"   INTEGER            NOT NULL DEFAULT 0,
    "type"         SMALLINT           NOT NULL DEFAULT 0,
    "email"        VARCHAR(62) UNIQUE NOT NULL DEFAULT '',
    "password"     VARCHAR(64)        NOT NULL DEFAULT '',
    "name"         VARCHAR(50)        NOT NULL DEFAULT '',
    "surname"      VARCHAR(50)        NOT NULL DEFAULT '',
    "phone_number" VARCHAR(16)        NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS "public"."buildings"
(
    "id"           SERIAL PRIMARY KEY,
    "city_id"      SMALLINT      NOT NULL DEFAULT 0,
    "district_id"  SMALLINT      NOT NULL DEFAULT 0,
    "cash_amount"  REAL          NOT NULL DEFAULT 0.00,
    "name"         VARCHAR(96)   NOT NULL DEFAULT '',
    "phone_number" VARCHAR(16)   NOT NULL DEFAULT '',
    "tax_number"   VARCHAR(16)   NOT NULL DEFAULT '',
    "address"      VARCHAR(1000) NOT NULL DEFAULT '',
    CONSTRAINT fk_city
        FOREIGN KEY (city_id)
            REFERENCES cities (id),
    CONSTRAINT fk_district
        FOREIGN KEY (district_id)
            REFERENCES districts (id)
);

CREATE TABLE IF NOT EXISTS "public"."blocks"
(
    "id"          SERIAL PRIMARY KEY,
    "building_id" INTEGER    NOT NULL,
    "letter"      VARCHAR(6) NOT NULL DEFAULT '',
    "d_number"    VARCHAR(6) NOT NULL DEFAULT '',
    CONSTRAINT uniq_letter
        UNIQUE (building_id, letter),
    CONSTRAINT fk_bui
        FOREIGN KEY (building_id)
            REFERENCES buildings (id)
);

CREATE TABLE IF NOT EXISTS "public"."user_building_membership"
(
    "user_id"     INTEGER,
    "building_id" INTEGER,
    "rank"        SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT uid_ref FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT aid_ref FOREIGN KEY (building_id) REFERENCES buildings (id),
    CONSTRAINT uap_pkey
        PRIMARY KEY (user_id, building_id)
);

CREATE TABLE IF NOT EXISTS "public"."flats"
(
    "id"          SERIAL PRIMARY KEY,
    "building_id" INTEGER,
    "block_id"    INTEGER,
    "owner_id"    INTEGER,
    "tenant_id"   INTEGER,
    "type"        SMALLINT   NOT NULL DEFAULT 0,
    "number"      VARCHAR(4) NOT NULL DEFAULT '',
    CONSTRAINT owte_chk CHECK (owner_id <> tenant_id),
    CONSTRAINT flt_uniq UNIQUE (block_id, number),
    CONSTRAINT aid_key FOREIGN KEY (building_id) REFERENCES buildings (id),
    CONSTRAINT blk_key FOREIGN KEY (block_id) REFERENCES blocks (id),
    CONSTRAINT ou_key FOREIGN KEY (owner_id) REFERENCES users (id),
    CONSTRAINT tu_key FOREIGN KEY (tenant_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS "public"."revenues"
(
    "id"              SERIAL PRIMARY KEY,
    "building_id"     INTEGER,
    "flat_id"         INTEGER,
    "rid"             INTEGER      NOT NULL,
    "time"            timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "paid_time"       timestamptz,
    "total"           REAL         NOT NULL DEFAULT 0.00,
    "paid_type"       SMALLINT     NOT NULL DEFAULT 0,
    "payer_full_name" VARCHAR(100) NOT NULL DEFAULT '',
    "payer_email"     VARCHAR(62)  NOT NULL DEFAULT '',
    "payer_phone"     VARCHAR(16)  NOT NULL DEFAULT '',
    "paid_status"     BOOL         NOT NULL DEFAULT false,
    "details"         VARCHAR(128) NOT NULL DEFAULT '',
    CONSTRAINT arid_uniq UNIQUE (building_id, rid),
    CONSTRAINT aid_key FOREIGN KEY (building_id) REFERENCES buildings (id),
    CONSTRAINT fid_key FOREIGN KEY (flat_id) REFERENCES flats (id)
);

CREATE TABLE IF NOT EXISTS "public"."expenses"
(
    "id"          SERIAL PRIMARY KEY,
    "building_id" INTEGER,
    "eid"         INTEGER      NOT NULL,
    "total"       REAL         NOT NULL DEFAULT 0.00,
    "time"        timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "paid_time"   timestamptz,
    "paid_type"   SMALLINT     NOT NULL DEFAULT 0,
    "to_name"     VARCHAR(100) NOT NULL DEFAULT '',
    "to_email"    VARCHAR(62)  NOT NULL DEFAULT '',
    "to_phone"    VARCHAR(16)  NOT NULL DEFAULT '',
    "paid_status" BOOL         NOT NULL DEFAULT false,
    "details"     VARCHAR(128) NOT NULL DEFAULT '',
    CONSTRAINT aeid_uniq UNIQUE (building_id, eid),
    CONSTRAINT aid_ekey FOREIGN KEY (building_id) REFERENCES buildings (id)
);