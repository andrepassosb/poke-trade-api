CREATE TABLE user_cards (
    user_id INTEGER NOT NULL,
    card_id TEXT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    PRIMARY KEY (user_id, card_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);