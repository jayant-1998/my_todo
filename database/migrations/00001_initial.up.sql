begin ;
create table if not exists users (
    id serial primary key,
    name TEXT not null ,
    email TEXT not null unique,
    created_at timestamp with time zone default now(),
    updated_at TIMESTAMP WITH TIME ZONE default now(),
    archived_at TIMESTAMP WITH TIME ZONE
);

create table if not exists todo(
    id  serial primary key ,
    user_id int references users(id) not null,
    task TEXT NOT NULL,
    description TEXT NOT NULL,
    is_completed boolean default false,
    due_date TIMESTAMP WITH TIME ZONE default  now() + interval '5 DAY',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
commit ;