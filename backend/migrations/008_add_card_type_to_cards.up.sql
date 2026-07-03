-- Card types: basic (front→back) or cloze ({{c1::...}} deletions in front)
ALTER TABLE cards ADD COLUMN card_type VARCHAR(10) NOT NULL DEFAULT 'basic' CHECK (card_type IN ('basic', 'cloze'));
