CREATE TABLE IF NOT EXISTS posts (
    id TEXT NOT NULL PRIMARY KEY DEFAULT ('post-' || gen_random_uuid() ),
    title VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    post_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS experiences (
    id TEXT NOT NULL PRIMARY KEY DEFAULT ('experience-' || gen_random_uuid()),
    start_date DATE NOT NULL,
    end_date DATE, 
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    entity VARCHAR(255) NOT NULL,
    url VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_posts_type ON posts(post_type);
