CREATE TABLE IF NOT EXISTS users(
                                    id INT PRIMARY KEY,
                                    active BOOL
);


CREATE TABLE IF NOT EXISTS chats(
                                    id INT PRIMARY KEY,
                                    name STRING,
                                    auto_pick BOOL,
                                    all_phrases BOOL,
                                    last_time_triggered DATE,
                                    selected_user_id INT REFERENCES users (id)
);


CREATE TABLE IF NOT EXISTS game_stats (
                                          id SERIAL PRIMARY KEY,
                                          user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                                          chat_id INT NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
                                          num  INT,
                                          year INT
);

ALTER TABLE game_stats
    ADD CONSTRAINT game_stat  UNIQUE (user_id, chat_id,year);


CREATE TABLE IF NOT EXISTS phrases (
                                       id SERIAL PRIMARY KEY,
                                       user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
                                       phrase  STRING,
                                       type  STRING,
                                       public BOOL
);
