-- Share codes let another user preview and clone a deck
ALTER TABLE decks ADD COLUMN share_code VARCHAR(12) UNIQUE;
