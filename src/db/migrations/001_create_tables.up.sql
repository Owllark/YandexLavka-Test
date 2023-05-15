CREATE TABLE IF NOT EXISTS couriers (
        courier_id    BIGSERIAL PRIMARY KEY,
        type          text,
        regions       integer[],
        working_hours text[]
);

CREATE TABLE IF NOT EXISTS orders (
      order_id       BIGSERIAL PRIMARY KEY,
      weight        real,
      region        integer,
      delivery_hours text[],
      cost          integer,
      completed_time timestamp
);

CREATE TABLE IF NOT EXISTS completed_orders (
        id BIGSERIAL PRIMARY KEY,
        courier_id INTEGER,
        order_id INTEGER,
        complete_time TIMESTAMP,
        FOREIGN KEY (courier_id) REFERENCES couriers (courier_id),
        FOREIGN KEY (order_id) REFERENCES orders (order_id)
);


