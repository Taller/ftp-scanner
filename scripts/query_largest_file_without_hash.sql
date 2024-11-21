SELECT f.id, f.folder_id, full_name, max(size) size
    FROM ftp.file f left join ftp.file_attr fa on f.id=fa.file_id and fa.hash is null
    group by f.id, full_name, path, folder_id
    order by  size desc
    limit 1 ;