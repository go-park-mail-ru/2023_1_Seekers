CREATE SCHEMA mail;

CREATE TABLE mail.users
(
    user_id          bigserial                NOT NULL,
    here_since       timestamp with time zone NOT NULL DEFAULT current_timestamp,
    is_deleted       boolean                  NOT NULL DEFAULT false,
    email            text                     NOT NULL,
    password         bytea                    NOT NULL,
    first_name       text,
    last_name        text,
    avatar           text                     NOT NULL,
    is_custom_avatar bool                     NOT NULL DEFAULT false,
    is_external      bool                     NOT NULL DEFAULT false,
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
) ,
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
    message_id          bigserial NOT NULL,
    from_user_id        bigint    NOT NULL,
    size                integer,                  --NOT NULL,

    title               text,
    reply_to_message_id bigint default null,

    created_at          timestamp with time zone, -- дата редактирования или отправки (финально)
    text                text,


    CONSTRAINT pk_messages PRIMARY KEY (message_id),
--     CONSTRAINT check_size CHECK (
--         size >= 0
--         ),
    CONSTRAINT fk_messages_users_user_id FOREIGN KEY (from_user_id)
        REFERENCES mail.users ON DELETE restrict,
    constraint fk_reply_to_message_message_id FOREIGN KEY (reply_to_message_id)
        REFERENCES mail.messages ON DELETE restrict
);


CREATE TABLE mail.boxes
(
    user_id    bigint  NOT NULL,
    message_id bigint  NOT NULL,
    folder_id  bigint  NOT NULL,
    seen       boolean NOT null,
    favorite   boolean NOT null default false,
    deleted    boolean NOT null default false,
    is_draft   boolean NOT NULL default false,

--     CONSTRAINT pk_box PRIMARY KEY (user_id, message_id),
--    CONSTRAINT uk_box_id UNIQUE (user_id, folder_id),
    CONSTRAINT fk_box_messages_user_id FOREIGN KEY (user_id)
        REFERENCES mail.users ON DELETE RESTRICT INITIALLY deferred,
    CONSTRAINT fk_box_messages_message_id FOREIGN KEY (message_id)
        REFERENCES mail.messages ON DELETE cascade,
    CONSTRAINT fk_box_messages_folder_id FOREIGN KEY (folder_id)
        REFERENCES mail.folders ON DELETE RESTRICT INITIALLY DEFERRED
);

-- триггер на увеличение непрочитанного и общего числа сообщений
CREATE
OR REPLACE FUNCTION increment_count_messages()
    RETURNS TRIGGER
AS
$BODY$
BEGIN
UPDATE mail.folders
SET messages_count  =
        CASE
            WHEN local_name = 'drafts' OR NEW.is_draft = false THEN messages_count + 1
            ELSE messages_count
            END,
    messages_unseen =
            CASE
                WHEN local_name != 'outbox' AND (local_name = 'drafts' OR NEW.is_draft = false) AND NEW.seen = false THEN messages_unseen + 1
                ELSE messages_unseen
END
WHERE folders.folder_id = NEW.folder_id;
RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;


-- срабатывает после вставки записей в boxes
CREATE TRIGGER inc_cnt_after_ins_box
    AFTER INSERT
    ON mail.boxes
    FOR EACH ROW
    EXECUTE PROCEDURE increment_count_messages();


-- триггер на изменение (+1 или -1) количества непрочитанных сообщений при прочитывании (или наоборот) сообщения
CREATE
OR REPLACE FUNCTION update_count_messages_after_seen()
    RETURNS TRIGGER
AS
$BODY$
BEGIN
UPDATE mail.folders
SET messages_unseen =
        CASE
            WHEN NEW.seen = true THEN messages_unseen - 1
            ELSE messages_unseen + 1
            END
WHERE folders.folder_id = NEW.folder_id;
RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

-- срабатывает после обновления столбца seen в boxes
CREATE TRIGGER update_cnt_after_update_seen
    AFTER UPDATE
    ON mail.boxes
    FOR EACH ROW
    WHEN (OLD.seen IS DISTINCT FROM NEW.seen)
EXECUTE PROCEDURE update_count_messages_after_seen();

-- триггер на изменение (+1 или -1) количества непрочитанных сообщений при переносе сообщения в другую папку
CREATE
OR REPLACE FUNCTION update_count_messages_after_move()
    RETURNS TRIGGER
AS
$BODY$
BEGIN
UPDATE mail.folders
SET messages_count  = messages_count + 1,
    messages_unseen =
        CASE
            WHEN NEW.seen = false THEN messages_unseen + 1
            ELSE messages_unseen
            END
WHERE folders.folder_id = NEW.folder_id;

UPDATE mail.folders
SET messages_count  = messages_count - 1,
    messages_unseen =
        CASE
            WHEN OLD.seen = false THEN messages_unseen - 1
            ELSE messages_unseen
            END
WHERE folders.folder_id = OLD.folder_id;
RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

-- срабатывает после обновления столбца folder_id в boxes
CREATE TRIGGER update_cnt_after_update_move
    AFTER UPDATE
    ON mail.boxes
    FOR EACH ROW
    WHEN (OLD.folder_id IS DISTINCT FROM NEW.folder_id)
EXECUTE PROCEDURE update_count_messages_after_move();

-- триггер по уменьшению количества сообщение после удаления
CREATE
OR REPLACE FUNCTION update_count_messages_after_delete()
    RETURNS TRIGGER
AS
$BODY$
BEGIN
UPDATE mail.folders
SET messages_count  =
        CASE
            WHEN local_name = 'drafts' OR OLD.is_draft = false THEN messages_count - 1
            ELSE messages_count
            END,
    messages_unseen =
        CASE
            WHEN OLD.seen = false AND (local_name = 'drafts' OR OLD.is_draft = false) THEN messages_unseen - 1
            ELSE messages_unseen
            END
WHERE folders.folder_id = OLD.folder_id;
RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

-- срабатывает после удаления столбца записи из boxes
CREATE TRIGGER update_cnt_after_delete
    AFTER DELETE
    ON mail.boxes
    FOR EACH ROW
    EXECUTE PROCEDURE update_count_messages_after_delete();

-- для поиска сообщений по тексту, получателю и отправителю
CREATE
OR REPLACE FUNCTION get_messages(from_id bigint, from_email text, to_email text, folder text, filter_text text)
    RETURNS TABLE
            (
                id bigint
            )
AS
$$
BEGIN
RETURN QUERY(SELECT messages.message_id --, messages.text, messages.title, boxes.user_id, users.email
                  from mail.messages
                           JOIN mail.folders on folders.user_id = from_id
                           JOIN mail.boxes
                                on boxes.message_id = messages.message_id AND boxes.folder_id = folders.folder_id AND
                                   boxes.user_id = from_id
                           JOIN mail.users on users.user_id = messages.from_user_id
                  WHERE (local_name ilike '%' || folder || '%' OR
                         name ilike '%' || folder || '%')
                    AND (email ilike '%' || from_email || '%'
                      AND email ilike '%' || to_email || '%')
                    AND (email ilike '%' || filter_text || '%'
                      OR title ilike '%' || filter_text || '%'
                      OR text ilike '%' || filter_text || '%')
    ORDER BY messages.message_id DESC);
);
end
$$
language 'plpgsql';