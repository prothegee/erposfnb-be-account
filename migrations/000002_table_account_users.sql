create table if not exists account.users(
    id              uuid        unique primary key default uuidv7(),
    id_ref          text        unique not null
                                generated always as (replace(id::text, '-', '')) stored,
    email           text        unique not null,
    phone_number    text        null,
    username        text        unique not null
                                generated always
                                as (regexp_replace(email, '[.@]', '', 'g')) stored,
    password        text        not null,
    banned          boolean     not null default false,
    banned_reason   text        null,
    pub_path        text        unique not null generated always
                                as ('t/' || replace(id::text, '-', '')) stored,
    pub_pict        text        null generated always
                                as (
                                    'https://dummyjson.com/icon/'
                                    ||
                                    replace(id::text, '-', '')
                                    ||
                                    '/150'
                                ) stored,
    date_registered timestamp   not null default now(),
    last_update     timestamp   null
);

create index if not exists idx_account_users_id_ref on account.users(id_ref);
create index if not exists idx_account_users_email_and_username
    on account.users(email, username);
