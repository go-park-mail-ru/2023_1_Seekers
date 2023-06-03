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
    is_fake      bool                         NOT NULL DEFAULT false,
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
    name  text,
    email text
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

CREATE TABLE mail.attaches
(
    attach_id  bigserial NOT NULL,
    message_id bigint    NOT NULL,
    type       text,
    filename   text,
    s3_fname   text,
    size_str   text,
    size_count bigint,


    CONSTRAINT pk_attaches PRIMARY KEY (attach_id),
    CONSTRAINT fk_attaches_messages_message_id FOREIGN KEY (message_id)
        REFERENCES mail.messages ON DELETE cascade
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

CREATE TABLE mail.user2fake
(
    user_id bigint NOT NULL,
    fake_id bigint NOT NULL,

    CONSTRAINT user2fake_user FOREIGN KEY (user_id)
        REFERENCES mail.users ON DELETE cascade,
    CONSTRAINT user2fake_fake FOREIGN KEY (fake_id)
        REFERENCES mail.users ON DELETE cascade
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
                WHEN local_name != 'outbox' AND (local_name = 'drafts' OR NEW.is_draft = false) AND NEW.seen = false
                    THEN messages_unseen + 1
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

-- срабатывает после пометки строки из boxes как deleted
CREATE TRIGGER update_cnt_after_mark_deleted
    AFTER UPDATE
    ON mail.boxes
    FOR EACH ROW
    WHEN (OLD.deleted IS DISTINCT FROM NEW.deleted AND NEW.deleted = true)
EXECUTE PROCEDURE update_count_messages_after_delete();

-- для поиска сообщений по тексту, получателю и отправителю
CREATE
    OR REPLACE FUNCTION get_messages(in_folder_id bigint, from_email text, to_email text, filter_text TEXT, in_is_draft bool)
    RETURNS TABLE
            (
                id bigint
            )
AS
$$
BEGIN
    RETURN QUERY
        (
            SELECT DISTINCT ON (m.message_id, m.created_at) m.message_id
            FROM mail.folders AS f
                     JOIN mail.boxes AS b ON f.folder_id = in_folder_id AND
                                             f.folder_id = b.folder_id AND
                                             is_draft = in_is_draft
                     JOIN mail.messages AS m ON m.message_id = b.message_id
                     JOIN mail.users AS u_from ON m.from_user_id = u_from.user_id
                     LEFT JOIN mail.boxes AS b2 ON b.message_id = b2.message_id AND
                                                   (b2.user_id != m.from_user_id OR
                                                    b2.folder_id != in_folder_id)
                     LEFT JOIN mail.users AS u_to ON b2.user_id = u_to.user_id
            WHERE (filter_text = '' AND from_email = '' AND to_email = '') or (filter_text != '' AND
                                                                               (m.title ILIKE '%' || filter_text || '%' OR
                                                                                m.text ILIKE '%' || filter_text || '%')) OR (from_email != '' AND
                                                                                                                             u_from.email ILIKE '%' || from_email || '%') OR (to_email != '' AND
                                                                                                                                                                              u_to.email ILIKE '%' || to_email || '%')
            ORDER BY
                m.created_at,
                m.message_id
        );
end
$$
    language 'plpgsql';

-- для поиска получателей для конкретного пользователя, сначала недавние
CREATE
    OR REPLACE FUNCTION get_recipes(from_id bigint[])
    RETURNS TABLE
            (
                user_id    bigint,
                first_name text,
                last_name  text,
                email      text
            )
AS
$$
    # variable_conflict use_column
BEGIN
    RETURN QUERY (SELECT user_id, first_name, last_name, email
                  FROM (select DISTINCT on (users.user_id) users.user_id,
                                                           users.first_name,
                                                           users.last_name,
                                                           users.email,
                                                           messages.message_id
                        from mail.messages
                                 join mail.boxes on boxes.message_id = messages.message_id
                                 join mail.users on boxes.user_id = users.user_id
                        where from_user_id = ANY(from_id)
                          and users.user_id != ALL(from_id)) as all_recipes
                  order by all_recipes.message_id desc);
end
$$
    language 'plpgsql';