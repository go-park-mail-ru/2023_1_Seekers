SET search_path TO mail, public;

CREATE USER mail_service_user WITH PASSWORD :db_mail_service_user_pw;
CREATE USER user_service_user WITH PASSWORD :db_user_service_user_pw;

GRANT USAGE ON SCHEMA mail TO mail_service_user;
GRANT USAGE ON SCHEMA mail TO user_service_user;

GRANT SELECT, UPDATE, INSERT, DELETE, TRIGGER ON folders, boxes TO mail_service_user;
GRANT SELECT, UPDATE, INSERT ON messages, attaches TO mail_service_user;

GRANT SELECT, UPDATE, INSERT, DELETE  ON users TO user_service_user;
