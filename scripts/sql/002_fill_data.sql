insert into mail.users(is_deleted, email, password, first_name, last_name, avatar) values
(false, 'support@mailbox.ru', '\x344b4f7968757548475a905c4afcff00a583d73bd6e98e612225a74563da5478eb10b08dcfff4280bf13'::bytea, 'Support', 'Testov', 'default_avatar.png'),
(false, 'test@mailbox.ru', '\x344b4f7968757548475af7b7075ce8720f3ddb16b80c6603b8c46d91d410479219489f516038b514bdf7'::bytea, 'Michail', 'Testov', 'default_avatar.png'),
(false, 'gena@mailbox.ru', '\x344b4f7968757548475a16d5a64e9ef192f45a78239b7429fb886e468a6c01e94f238f884d29986f8fa5'::bytea, 'Ivan', 'Ivanov', 'default_avatar.png'),
(false, 'max@mailbox.ru', '\x344b4f7968757548475afb7a9c8404720c9f041fe049f1cf6db37d269d06e15be8c3fee26e60ec1fd204'::bytea, 'Michail', 'Sidorov', 'default_avatar.png'),
(false, 'valera@mailbox.ru', '\x344b4f7968757548475af7b7075ce8720f3ddb16b80c6603b8c46d91d410479219489f516038b514bdf7'::bytea, 'Valera', 'Testov', 'default_avatar.png');

select * from mail.users;

insert into mail.folders(user_id, local_name, name, messages_unseen, messages_count) values
(1, 'inbox', 'Входящие', 0, 0), --1
(1, 'outbox', 'Отправленные', 0, 0),
(1, 'trash', 'Корзина', 0, 0),
(1, 'drafts', 'Черновики', 0, 0),
(1, 'spam', 'Спам', 0, 0),
(2, 'inbox', 'Входящие', 8, 8), --6
(2, 'outbox', 'Отправленные', 0, 9),
(2, 'trash', 'Корзина', 0, 0),
(2, 'drafts', 'Черновики', 0, 0),
(2, 'spam', 'Спам', 0, 0),
(3, 'inbox', 'Входящие', 4, 4), --11
(3, 'outbox', 'Отправленные', 0, 4),
(3, 'trash', 'Корзина', 0, 0),
(3, 'drafts', 'Черновики', 0, 0),
(3, 'spam', 'Спам', 0, 0),
(4, 'inbox', 'Входящие', 5, 5), --16
(4, 'outbox', 'Отправленные', 0, 4),
(4, 'trash', 'Корзина', 0, 0),
(4, 'drafts', 'Черновики', 0, 0),
(4, 'spam', 'Спам', 0, 0),
(5, 'inbox', 'Входящие', 0, 0), --21
(5, 'outbox', 'Отправленные', 0, 0),
(5, 'trash', 'Корзина', 0, 0),
(5, 'drafts', 'Черновики', 0, 0),
(5, 'spam', 'Спам', 0, 0),
(3, '1', 'My', 0, 0),
(4, '1', 'Empty', 0, 0);

select * from mail.folders;

insert into mail.messages(from_user_id, size, title, created_at, text) values 
(3, 100, 'Invitation', '2023-01-01', 'Hello, we decided to invite you to our party, lets go it will be fine!'),
(4, 100, 'Spam letter', '2023-01-02', 'Nunc non velit commodo, vestibulum enim ullamcorper, lobortis mi. Integer eu elit nibh. Integer bibendum semper arcu, eget consectetur nisi gravida eu. Suspendisse maximus id urna a volutpat. Quisque nec iaculis purus, non facilisis massa. Maecenas finibus dui ipsum, ut tempor sapien tincidunt blandit. Ut at iaculis eros, ultrices iaculis nibh. Mauris fermentum elit erat, at cursus urna euismod vel. In congue, ipsum a fermentum semper, dolor sem scelerisque leo, a tempus risus orci eu leo. Fusce vulputate venenatis imperdiet. Vestibulum interdum pellentesque facilisis'),
(3, 100, 'Lorem', '2023-01-04', 'Mauris imperdiet massa ante. Pellentesque feugiat nisl nec ultrices laoreet. Aenean a mauris mi. Sed auctor egestas nulla et vulputate. Praesent lobortis nulla ante, vel dignissim odio aliquet et. Suspendisse potenti. Donec venenatis nibh a sem consectetur, bibendum consectetur metus venenatis. Mauris lorem tellus, finibus id dui sit amet, facilisis fermentum orci. Mauris arcu ante, lacinia vitae orci in, tempus elementum lacus. Donec eu augue vulputate, tempor neque nec, efficitur purus. Mauris ut lorem non sapien placerat mattis. In in lacus a lorem viverra laoreet ut et orci. Maecenas auctor, justo nec hendrerit interdum, nibh nisi consectetur sapien, id ultrices lacus mi sed risus.'),
(4, 100, 'Very interesting letter', '2023-01-05', 'Morbi sit amet porttitor sapien, eget venenatis est. Suspendisse sollicitudin elit velit, quis sodales dolor maximus id. Vestibulum gravida scelerisque nibh, sit amet tincidunt augue gravida nec. Maecenas non placerat justo, at feugiat nulla. Phasellus dapibus a mi ut interdum. Aliquam nec quam feugiat, rutrum urna ut, cursus purus. Lorem ipsum dolor sit amet, consectetur adipiscing elit.'),
(3, 100, 'Small text letter', '2023-01-06', 'Hi! how are you?'),
(4, 100, 'Do you like to read books?', '2023-01-06', 'We have a lot of new books that may interest you'),
(3, 100, 'Advertisement', '2023-01-07', 'Hi, visit our shop!'),
(4, 100, 'Let’s get acquainted', '2023-01-29', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout'),
(2, 100, 'Not spam', '2023-01-29', 'Open this letter please'),
(2, 100, 'Vacancy', '2023-01-01', 'We will be glad to offer you this job'),
(2, 100, 'Invitation', '2023-01-01', 'Hello, we decided to invite you to our party, lets go it will be fine!'),
(2, 100, 'Spam letter', '2023-01-02', 'Nunc non velit commodo, vestibulum enim ullamcorper, lobortis mi. Integer eu elit nibh. Integer bibendum semper arcu, eget consectetur nisi gravida eu. Suspendisse maximus id urna a volutpat. Quisque nec iaculis purus, non facilisis massa. Maecenas finibus dui ipsum, ut tempor sapien tincidunt blandit. Ut at iaculis eros, ultrices iaculis nibh. Mauris fermentum elit erat, at cursus urna euismod vel. In congue, ipsum a fermentum semper, dolor sem scelerisque leo, a tempus risus orci eu leo. Fusce vulputate venenatis imperdiet. Vestibulum interdum pellentesque facilisis'),
(2, 100, 'Lorem', '2023-01-04', 'Mauris imperdiet massa ante. Pellentesque feugiat nisl nec ultrices laoreet. Aenean a mauris mi. Sed auctor egestas nulla et vulputate. Praesent lobortis nulla ante, vel dignissim odio aliquet et. Suspendisse potenti. Donec venenatis nibh a sem consectetur, bibendum consectetur metus venenatis. Mauris lorem tellus, finibus id dui sit amet, facilisis fermentum orci. Mauris arcu ante, lacinia vitae orci in, tempus elementum lacus. Donec eu augue vulputate, tempor neque nec, efficitur purus. Mauris ut lorem non sapien placerat mattis. In in lacus a lorem viverra laoreet ut et orci. Maecenas auctor, justo nec hendrerit interdum, nibh nisi consectetur sapien, id ultrices lacus mi sed risus.'),
(2, 100, 'Very interesting letter', '2023-01-05', 'Morbi sit amet porttitor sapien, eget venenatis est. Suspendisse sollicitudin elit velit, quis sodales dolor maximus id. Vestibulum gravida scelerisque nibh, sit amet tincidunt augue gravida nec. Maecenas non placerat justo, at feugiat nulla. Phasellus dapibus a mi ut interdum. Aliquam nec quam feugiat, rutrum urna ut, cursus purus. Lorem ipsum dolor sit amet, consectetur adipiscing elit.'),
(2, 100, 'Small text letter', '2023-01-06', 'Hi! how are you?'),
(2, 100, 'Not spam', '2023-01-06', 'Open this letter please'),
(2, 100, 'Advertisement', '2023-01-07', 'Hi, visit our shop!');

select * from mail.messages;

insert into mail.boxes(user_id, message_id, folder_id, seen) values
(3, 1, 12, true),
(4, 2, 17, true),
(3, 3, 12, true),
(4, 4, 17, true),
(3, 5, 12, true),
(4, 6, 17, true),
(3, 7, 12, true),
(4, 8, 17, true),
(2, 9, 7, true),
(2, 10, 7, true),
(2, 11, 7, true),
(2, 12, 7, true),
(2, 13, 7, true),
(2, 14, 7, true),
(2, 15, 7, true),
(2, 16, 7, true),
(2, 17, 7, true),
(2, 1, 6, false),
(2, 2, 6, false),
(2, 3, 6, false),
(2, 4, 6, false),
(2, 5, 6, false),
(2, 6, 6, false),
(2, 7, 6, false),
(2, 8, 6, false),
(4, 9, 16, false),
(3, 10, 11, false),
(4, 11, 16, false),
(3, 12, 11, false),
(4, 13, 16, false),
(3, 14, 11, false),
(4, 15, 16, false),
(3, 16, 11, false),
(4, 17, 16, false);

select * from mail.boxes;
