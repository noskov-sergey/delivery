-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY         NOT NULL,
    courier_id UUID,
    location_x SMALLINT         NOT NULL,
    location_y SMALLINT         NOT NULL,
    volume INT                  NOT NULL,
    status VARCHAR(20)          NOT NULL
);

CREATE TABLE IF NOT EXISTS couriers (
    id UUID PRIMARY KEY         NOT NULL,
    name_ TEXT                  NOT NULL,
    speed INT                   NOT NULL,
    location_x SMALLINT         NOT NULL,
    location_y SMALLINT         NOT NULL,
    storage_places TEXT[]       NOT NULL
);

CREATE TABLE IF NOT EXISTS storage_places (
    id UUID PRIMARY KEY         NOT NULL,
    name_ TEXT                  NOT NULL,
    total_volume INT            NOT NULL,
    order_id UUID,
    courier_id UUID             NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders, couriers, storage_places;
-- +goose StatementEnd
