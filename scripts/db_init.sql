CREATE TABLE IF NOT EXISTS categories (
  id bigserial PRIMARY KEY,
  name varchar(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS transactions(
  id bigserial PRIMARY KEY,
  description varchar(255) NOT NULL,
  category_id bigint,
  amount decimal(10, 2) NOT NULL,
  date date NOT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_transactions_category ON transactions(category_id);

INSERT INTO categories (name) VALUES
  ('income'),
  ('interest'),
  ('rent'),
  ('utilities'),
  ('insurance'),
  ('dining'),
  ('groceries'),
  ('shopping'),
  ('entertainment'),
  ('subscriptions'),
  ('travel'),
  ('gifts'),
  ('investment'),
  ('emergency')
ON CONFLICT (name) DO NOTHING;