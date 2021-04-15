create table if not exists reg_users(
	user_id serial primary key,
	username varchar(100) not null,
	password varchar(100) not null,
	email varchar(50) unique not null,
	registered boolean default false
);

create table if not exists classes(
	cls_id serial primary key,
	user_id int,
	constraint fk_reg_user
		foreign key(user_id)
			references reg_users(user_id)
				on delete cascade
);

create table if not exists files(
	files_id serial primary key,
	storage_url varchar(200) not null,
	user_id int,
	cls_id int,
	constraint fk_reg_user
		foreign key(user_id)
			references reg_users(user_id)
				on delete cascade,
	constraint fk_class
		foreign key(cls_id)
			references classes(cls_id)
				on delete cascade
);

create table if not exists environments(
	classes_id serial primary key,
	user_id int,
	cls_id int,
	constraint fk_reg_user
		foreign key(user_id)
			references reg_users(user_id)
				on delete cascade,
	constraint fk_class
		foreign key(cls_id)
			references classes(cls_id)
				on delete cascade
);
create table if not exists oauth2_clients(
	id text not null primary key,
	secret text not null,
	domain text not null,
	data jsonb not null
	);
