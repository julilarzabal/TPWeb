-- 1. User table
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    age DATE NULL,
    industry VARCHAR(100) NULL,
    role VARCHAR(50) NOT NULL
);

-- 2. Category table
CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- 3. N:M relationship between USER and CATEGORY (user preferences)
CREATE TABLE user_preference (
    user_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (user_id, category_id),
    CONSTRAINT fk_user_preference_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_preference_category FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
);

-- 4. Event table
CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    location VARCHAR(200) NOT NULL,
    age_range VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_event_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- 5. N:M relationship between EVENT and CATEGORY
CREATE TABLE event_category (
    event_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (event_id, category_id),
    CONSTRAINT fk_event_category_event FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE,
    CONSTRAINT fk_event_category_category FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
);