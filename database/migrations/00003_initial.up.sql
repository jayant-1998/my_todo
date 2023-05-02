begin;

create table if not exists user_session (
                                            session_id serial primary key,
                                            user_id int references users(id) not null,
                                            session_token TEXT NOT NULL,
                                            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

alter table users add column password text ;

commit ;