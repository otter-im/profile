CREATE TABLE profile
(
    id                  uuid         NOT NULL PRIMARY KEY,
    display_name        varchar(255) NULL,
    description         TEXT         NULL,
    location            varchar(128) NULL,
    website             varchar(255) NULL,
    pronouns            varchar(255) NULL,
    date_of_birth       timestamp    NULL,
    avatar_url          TEXT         NOT NULL DEFAULT 'https://s3.otter.im/profile/default.png'
);

INSERT INTO profile (id, display_name, description, location, website, pronouns, date_of_birth)
VALUES ('d6ef6dc7-ce36-449c-8265-07f60ca3b2ff', 'Pico! ✨🇸🇪', 'Founder of Otter Social', 'Stockholm, Sweden', NULL,
        'she/her', '1995-01-22 00:00:00');
