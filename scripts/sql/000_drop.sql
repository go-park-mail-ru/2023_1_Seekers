DROP USER IF EXISTS auth_service_user, dwh_service_user, mail_service_user, user_service_user;

DROP TABLE IF EXISTS  mail.users CASCADE;
DROP TABLE IF EXISTS  mail.accounts CASCADE;
DROP TABLE IF EXISTS  mail.folders CASCADE;
DROP TABLE IF EXISTS  mail.messages CASCADE;
DROP TABLE IF EXISTS  mail.box CASCADE;
DROP SCHEMA IF EXISTS mail CASCADE;