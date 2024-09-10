-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS employee
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE IF NOT EXISTS organization
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_id         UUID REFERENCES employee (id) ON DELETE CASCADE
);

CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
    );

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
    );

CREATE TABLE IF NOT EXISTS tenders
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(100)                                        NOT NULL,
    description     VARCHAR(500),
    service_type    service_type,
    status          tender_status,
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE NOT NULL,
    creator_user_id uuid REFERENCES employee (id) ON DELETE CASCADE     NOT NULL,
    version         INT,
    created_at      TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE bid_decision AS ENUM (
    'Approved',
    'Rejected'
    );

CREATE TYPE bid_status AS ENUM (
    'Approved',
    'Canceled',
    'Created',
    'Published',
    'Rejected'
    );

CREATE TYPE bid_author AS ENUM (
    'user',
    'organization'
    );

CREATE TABLE IF NOT EXISTS bids
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id       uuid REFERENCES tenders (id) ON DELETE CASCADE      NOT NULL,
    creator_id      uuid REFERENCES employee (id) ON DELETE CASCADE     NOT NULL,
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE NOT NULL,
    decision        bid_decision,
    status          bid_status,
    author_type     bid_author,
    version         INT,
    created_at      TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS tenders;
DROP TABLE IF EXISTS organization_responsible;
DROP TABLE IF EXISTS organization;
DROP TABLE IF EXISTS employee;

DROP TYPE IF EXISTS service_type;
DROP TYPE IF EXISTS tender_status;
DROP TYPE IF EXISTS organization_type;
DROP TYPE IF EXISTS bid_decision;
DROP TYPE IF EXISTS bid_status;
DROP TYPE IF EXISTS bid_author;
-- +goose StatementEnd
