create table campaigns
(
    id        serial primary key,
    starts_at timestamptz not null,
    ends_at   timestamptz not null,
    price     int
);

create table orders
(
    id              SERIAL PRIMARY KEY,
    campaign_id     INT REFERENCES campaigns (id),

    -- Customer info
    cus_name        TEXT,
    cus_email       TEXT,

    -- Address info
    adr_street1     TEXT,
    adr_street2     TEXT,
    adr_city        TEXT,
    adr_state       TEXT,
    adr_zip         TEXT,
    adr_country     TEXT,
    adr_raw         TEXT,

    -- Payment info
    pay_source      TEXT,
    pay_customer_id TEXT,
    pay_charge_id   TEXT
);

CREATE INDEX cus_email_idx ON orders (cus_email);
CREATE INDEX cus_pay_cus_id_idx ON orders (pay_customer_id);