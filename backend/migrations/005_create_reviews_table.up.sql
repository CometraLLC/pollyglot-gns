-- Create reviews table (one row per card review, feeds progress stats)
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 0 AND 4),
    reviewed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reviews_user_reviewed_at ON reviews(user_id, reviewed_at);
CREATE INDEX idx_reviews_card_id ON reviews(card_id);
