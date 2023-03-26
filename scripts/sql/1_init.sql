CREATE SCHEMA mail;

CREATE TABLE mail.users
(
    user_id    bigserial                NOT NULL,
    here_since timestamp with time zone NOT NULL DEFAULT current_timestamp,
    is_deleted boolean                  NOT NULL DEFAULT false,
    email      text                     NOT NULL,
    password   text                     NOT NULL,
    first_name text                     NOT NULL,
    last_name  text                     NOT NULL,
    avatar     text                     NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY (user_id)
);

CREATE TABLE mail.folders
(
    folder_id       bigserial NOT NULL,
    user_id         bigint    NOT NULL,
    local_name      text      NOT NULL,
    name            text      NOT NULL,
    messages_unseen integer   NOT NULL DEFAULT 0,
    messages_count  integer   NOT NULL DEFAULT 0,

    CONSTRAINT check_message_count CHECK (
                messages_count >= 0 AND
                messages_unseen >= 0 AND
                messages_count >= messages_unseen
        ),

    CONSTRAINT check_non_empty_name CHECK (
        name != ''
        ),
    CONSTRAINT check_folder_id_natural CHECK (folder_id > 0),
    CONSTRAINT pk_folders PRIMARY KEY (folder_id),
    CONSTRAINT fk_folders_user_id_users FOREIGN KEY (user_id)
        REFERENCES mail.users ON DELETE RESTRICT
);

CREATE TYPE mail.recipient AS
    (
    name text,
    email text
    );

CREATE TYPE mail.attach AS
    (
    type text,
    filename text,
    size integer,
    raw_data bytea
    );

CREATE TABLE mail.messages
(
    message_id   bigserial NOT NULL,
    from_user_id bigint    NOT NULL,
    size         integer   NOT NULL,

    title        text,
    reply_to     bigint default null,

    created_at   timestamp with time zone, -- дата редактирования или отправки (финально)
    text         text,


    CONSTRAINT pk_messages PRIMARY KEY (message_id),
    CONSTRAINT check_size CHECK (
        size >= 0
        ),
    CONSTRAINT fk_messages_users_user_id FOREIGN KEY (from_user_id)
        REFERENCES mail.users ON DELETE restrict,
    constraint fk_reply_to_message_message_id FOREIGN KEY (reply_to)
        REFERENCES mail.messages ON DELETE restrict
);


CREATE TABLE mail.box
(
    user_id    bigint  NOT NULL,
    message_id bigint  NOT NULL,
    folder_id  bigint  NOT NULL,
    seen       boolean NOT null,
    favourite  boolean NOT null default false,
    deleted    boolean NOT null default false,

    CONSTRAINT pk_box PRIMARY KEY (user_id, message_id),
--    CONSTRAINT uk_box_id UNIQUE (user_id, folder_id),
    CONSTRAINT fk_box_messages_user_id FOREIGN KEY (user_id)
        REFERENCES mail.users ON DELETE RESTRICT INITIALLY deferred,
    CONSTRAINT fk_box_messages_message_id FOREIGN KEY (message_id)
        REFERENCES mail.messages ON DELETE RESTRICT INITIALLY deferred,
    CONSTRAINT fk_box_messages_folder_id FOREIGN KEY (folder_id)
        REFERENCES mail.folders ON DELETE RESTRICT INITIALLY DEFERRED
);
