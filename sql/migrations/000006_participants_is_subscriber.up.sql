ALTER TABLE IF EXISTS participants
ADD COLUMN is_subscribe boolean DEFAULT false NOT NULL;