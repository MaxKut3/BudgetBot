-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE BudgetBotData (
    id          serial primary key,
    chatId      bigint,
    category    varchar(50),
    amount      integer,
    cur         varchar(3),
    amountRub   integer,
    month       integer,
    year        integer

);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE BudgetBotData;
-- +goose StatementEnd
