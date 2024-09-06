-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id       UUID PRIMARY KEY,
    nickname TEXT NOT NULL,
    email    TEXT NOT NULL
);

CREATE TABLE file_services
(
    id         UUID PRIMARY KEY,
    address    TEXT NOT NULL,
    scanned_at TIMESTAMP
);

CREATE TABLE files
(
    id              UUID PRIMARY KEY,
    name            TEXT NOT NULL,
    path            TEXT NOT NULL,
    file_service_id UUID NOT NULL REFERENCES file_services (id)
);
CREATE INDEX idx_files_name
    ON files (name);

CREATE TABLE covers
(
    id      UUID PRIMARY KEY,
    width   SMALLINT NOT NULL,
    height  SMALLINT NOT NULL,
    file_id UUID     NOT NULL REFERENCES files (id)
);

CREATE TABLE genres
(
    id   UUID PRIMARY KEY,
    name TEXT
);
INSERT INTO genres(id, name)
VALUES ('00000000-0000-0000-0000-000000000000', '???');
CREATE INDEX idx_genres_name
    ON genres (name);

CREATE TABLE albums
(
    id    UUID PRIMARY KEY,
    title TEXT
);
INSERT INTO albums(id, title)
VALUES ('00000000-0000-0000-0000-000000000000', '???');
CREATE INDEX idx_album_title
    ON albums (title);

CREATE TABLE artists
(
    id   UUID PRIMARY KEY,
    name TEXT
);
INSERT INTO artists(id, name)
VALUES ('00000000-0000-0000-0000-000000000000', '???');
CREATE INDEX idx_artists_name
    ON artists (name);

CREATE TABLE songs
(
    id             UUID PRIMARY KEY,
    file_id        UUID REFERENCES files (id),
    title          TEXT     NOT NULL,
    duration_ns    BIGINT   NOT NULL,
    cover_id       UUID REFERENCES covers (id),
    genre_id       UUID     NOT NULL REFERENCES genres (id),
    album_id       UUID     NOT NULL REFERENCES albums (id),
    artist_id      UUID     NOT NULL REFERENCES artists (id),
    year           SMALLINT,
    number         SMALLINT,
    comment        TEXT,
    channels       SMALLINT NOT NULL,
    bitrate_kbps   INTEGER  NOT NULL,
    sample_rate_hz INTEGER  NOT NULL,
    md5            CHAR(32) NOT NULL
);
CREATE INDEX ids_songs_title
    ON songs (title);
CREATE INDEX ids_songs_genre_id
    ON songs (genre_id);
CREATE INDEX ids_songs_album_id
    ON songs (album_id);
CREATE INDEX ids_songs_artist_id
    ON songs (artist_id);
CREATE INDEX idx_songs_md5
    ON songs (md5);

CREATE TABLE playlists
(
    id        UUID PRIMARY KEY,
    name      TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users (id)
);
CREATE INDEX idx_playlists_name
    ON playlists (name);

CREATE TABLE playlists_songs
(
    playlist_id UUID    NOT NULL REFERENCES playlists (id),
    song_id     UUID    NOT NULL REFERENCES songs (id),
    number      INTEGER NOT NULL,
    PRIMARY KEY (playlist_id, number)
);

CREATE TABLE saved_playlists
(
    user_id     UUID NOT NULL REFERENCES users (id),
    playlist_id UUID NOT NULL REFERENCES playlists (id),
    PRIMARY KEY (user_id, playlist_id)
);

CREATE TYPE playback_type AS ENUM ('DEFAULT', 'REPEAT_ONE', 'REPEAT_QUEUE', 'RANDOM');
CREATE TABLE rooms
(
    id            UUID PRIMARY KEY,
    owner_id      UUID          NOT NULL REFERENCES users (id),
    is_main       BOOLEAN       NOT NULL,
    playback_type playback_type NOT NULL,
    share_code    CHAR(32) UNIQUE
);

CREATE TABLE roommates
(
    room_id UUID NOT NULL REFERENCES rooms (id),
    user_id UUID NOT NULL REFERENCES users (id),
    PRIMARY KEY (room_id, user_id)
);

CREATE TABLE devices
(
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users (id),
    name            TEXT NOT NULL,
    current_room_id UUID REFERENCES rooms (id)
);

CREATE TABLE queue_items
(
    id           UUID PRIMARY KEY,
    song_id      UUID NOT NULL REFERENCES songs (id),
    prev_item_id UUID UNIQUE REFERENCES queue_items (id),
    next_item_id UUID UNIQUE REFERENCES queue_items (id)
);

CREATE TABLE scores
(
    song_id UUID     NOT NULL REFERENCES songs (id),
    user_id UUID     NOT NULL REFERENCES users (id),
    score   SMALLINT NOT NULL,
    PRIMARY KEY (song_id, user_id)
);
CREATE INDEX idx_scores_user_id
    ON scores (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE scores;

DROP TABLE queue_items;

DROP TABLE devices;

DROP TABLE roommates;

DROP TABLE rooms;
DROP TYPE playback_type;

DROP TABLE saved_playlists;

DROP TABLE playlists_songs;

DROP TABLE playlists;

DROP TABLE songs;

DROP TABLE artists;

DROP TABLE albums;

DROP TABLE genres;

DROP TABLE covers;

DROP TABLE files;

DROP TABLE file_services;

DROP TABLE users;
-- +goose StatementEnd
