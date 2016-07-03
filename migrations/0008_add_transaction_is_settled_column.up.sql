ALTER TABLE transactions
ADD is_settled boolean NOT NULL DEFAULT(FALSE);
