CREATE SCHEMA mail;

-- В таблице users хранится информация о пользователях
-- Таблица соответствует 3НФ
CREATE TABLE mail.users
(
    user_id          BIGSERIAL                NOT NULL,
    here_since       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
    email            TEXT                     UNIQUE NOT NULL,
    password         BYTEA                    NOT NULL,
    first_name       TEXT,
    last_name        TEXT,
    avatar           TEXT                     NOT NULL,
    is_custom_avatar BOOLEAN                  NOT NULL DEFAULT false,
    is_external      BOOLEAN                  NOT NULL DEFAULT false,
    is_deleted       BOOLEAN                  NOT NULL DEFAULT false,

    CONSTRAINT pk_users PRIMARY KEY (user_id)
);

-- В таблице folders хранится информация о папках
-- Таблица соотвествует 2НФ
-- Аттрибуты messages_unseen и messages_count являются зависимостями от первичного ключа folder_id,
-- но не являются простыми атрибутами, так как их значения могут быть вычислены на основе других значений в таблице.
-- Это было сделано для оптимизации запроса по подсчету количества сообщений в папке (общего и непрочитанных), а точнее,
-- для того, чтобы полностью исключить этот запрос
CREATE TABLE mail.folders
(
    folder_id       BIGSERIAL   NOT NULL,
    user_id         BIGINT      NOT NULL,
    local_name      TEXT        NOT NULL, -- локальное имя папки (может быть 'inbox', 'outbox', 'trash', 'spam', 'drafts' или любым числом в виде строки)
    name            TEXT        NOT NULL, -- имя папки, которое отображается на странице
    messages_unseen INTEGER     NOT NULL DEFAULT 0,
    messages_count  INTEGER     NOT NULL DEFAULT 0,

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

-- В таблице messages хранится общая информация о сообщениях
-- Таблица соотвесвует 3НФ
CREATE TABLE mail.messages
(
    message_id          BIGSERIAL                   NOT NULL,
    from_user_id        BIGINT                      NOT NULL,
    title               TEXT,
    reply_to_message_id BIGINT                      DEFAULT NULL,
    created_at          TIMESTAMP WITH TIME ZONE,
    text                TEXT,

    CONSTRAINT pk_messages PRIMARY KEY (message_id),
    CONSTRAINT fk_messages_users_user_id FOREIGN KEY (from_user_id)
        REFERENCES mail.users ON DELETE restrict,
    constraint fk_reply_to_message_message_id FOREIGN KEY (reply_to_message_id)
        REFERENCES mail.messages ON DELETE restrict
);

-- В таблице attaches хранится информация о вложениях
-- Таблица соотвесвует 3НФ
CREATE TABLE mail.attaches
(
    attach_id  BIGSERIAL NOT NULL,
    message_id BIGINT    NOT NULL,
    type       TEXT,
    filename   TEXT,
    s3_fname   TEXT,
    size_str   TEXT,
    size_count BIGINT,

    CONSTRAINT pk_attaches PRIMARY KEY (attach_id),
    CONSTRAINT fk_attaches_messages_message_id FOREIGN KEY (message_id)
        REFERENCES mail.messages ON DELETE cascade
);

-- В таблице boxes хранится информация о распложении сообщений в папке пользователя
-- То есть, это некая связка трех сущностей (сообщение, пользователь и папка)
-- Таблица соотвесвует 2НФ, так как аттрибут user_id является избыточным (так как по папке можно определить пользователя).
-- Он был добавлен для оптимизации запроса по выборке сообщений из папки - этот запрос требует доставать информацию о пользователе,
-- к которому относится письмо ("от кого" или "кому"). Чтобы не делать дополнительный запрос в таблицу folders (или не делать джойны),
-- для получения информации о пользователе, было решено добавить аттрибут user_id
CREATE TABLE mail.boxes
(
    user_id    BIGINT  NOT NULL,
    message_id BIGINT  NOT NULL,
    folder_id  BIGINT  NOT NULL,
    seen       BOOLEAN NOT NULL,
    favorite   BOOLEAN NOT NULL DEFAULT false,
    deleted    BOOLEAN NOT NULL DEFAULT false,
    is_draft   BOOLEAN NOT NULL DEFAULT false,

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
RETURN QUERY (SELECT messages.message_id --, messages.text, messages.title, boxes.user_id, users.email
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
end
$$
language 'plpgsql';

-- для поиска получателей для конкретного пользователя, сначала недавние
CREATE
OR REPLACE FUNCTION get_recipes(from_id bigint)
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
                        where from_user_id = from_id
                          and users.user_id != from_id) as all_recipes
                  order by all_recipes.message_id desc);
end
$$
language 'plpgsql';

VACUUM ANALYZE;