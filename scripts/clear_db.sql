truncate ftp.file cascade;
truncate ftp.file_attr cascade;
truncate ftp.folder cascade;
truncate ftp.empty_folder cascade;
truncate ftp.server cascade;

ALTER SEQUENCE ftp.files_id_seq RESTART WITH 1;
UPDATE ftp.file SET id=nextval('ftp.files_id_seq');

ALTER SEQUENCE ftp.file_attrs_id_seq RESTART WITH 1;
UPDATE ftp.file_attr SET id=nextval('ftp.file_attrs_id_seq');

ALTER SEQUENCE ftp.folder_id_seq RESTART WITH 1;
UPDATE ftp.folder SET id=nextval('ftp.folder_id_seq');

ALTER SEQUENCE ftp.empty_folder_id_seq RESTART WITH 1;
UPDATE ftp.empty_folder SET id=nextval('ftp.empty_folder_id_seq');

ALTER SEQUENCE ftp.server_id_seq RESTART WITH 1;
UPDATE ftp.server SET id=nextval('ftp.server_id_seq');

