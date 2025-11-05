CREATE TABLE IF NOT EXISTS subscriptions
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    serviceName TEXT NOT NULL,
    price INT NOT NULL,
    userId TEXT NOT NULL,
    startDate DATE,
    UNIQUE (userId, serviceName, startDate)
);