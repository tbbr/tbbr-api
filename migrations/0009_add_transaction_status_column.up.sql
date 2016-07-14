-- Force existing rows to have status of Confirmed
ALTER TABLE transactions
ADD COLUMN status text NOT NULL DEFAULT 'Confirmed';

-- Set default value for newly created rows to Pending
ALTER TABLE transactions
ALTER COLUMN status SET DEFAULT 'Pending';
