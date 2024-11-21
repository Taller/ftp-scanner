truncate scanner.file cascade;
truncate scanner.file_attr cascade;
truncate scanner.folder cascade;
truncate scanner.empty_folder cascade;
truncate scanner.server cascade;

ALTER SEQUENCE scanner.files_id_seq RESTART WITH 1;
UPDATE scanner.file SET id=nextval('scanner.files_id_seq');

ALTER SEQUENCE scanner.file_attrs_id_seq RESTART WITH 1;
UPDATE scanner.file_attr SET id=nextval('scanner.file_attrs_id_seq');

ALTER SEQUENCE scanner.folder_id_seq RESTART WITH 1;
UPDATE scanner.folder SET id=nextval('scanner.folder_id_seq');

ALTER SEQUENCE scanner.empty_folder_id_seq RESTART WITH 1;
UPDATE scanner.empty_folder SET id=nextval('scanner.empty_folder_id_seq');

ALTER SEQUENCE scanner.server_id_seq RESTART WITH 1;
UPDATE scanner.server SET id=nextval('scanner.server_id_seq');

