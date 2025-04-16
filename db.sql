-- create role stockinos createdb login encrypted password 'stockinos';

CREATEUSER --createdb --login -h localhost -p 5432 -U postgres -P stockinos

CREATEDB -h localhost -p 5432 -U stockinos stockinos;
