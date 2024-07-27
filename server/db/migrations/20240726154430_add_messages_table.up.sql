CREATE TABLE "messages" (
    id serial primary key,
    message varchar NOT NULL,
    status varchar DEFAULT 'idle',
    CONSTRAINT valid_status CHECK (status IN ('idle', 'send', 'received'))
);


