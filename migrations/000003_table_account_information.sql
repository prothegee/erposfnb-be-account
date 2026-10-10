create table if not exists account.informations(
    id_user         uuid        primary key,
    first_name      text        null,
    middle_name     text        null,
    last_name       text        null,
    street1         text        null,
    street2         text        null,
    district        text        null,
    city            text        null,
    province        text        null,
    postal_code     text        null,
    phone_numbers   text[]      null,
    backup_emails   text[]      null,
    last_update     timestamp   null,

    constraint  fk_account_informations_id_user
        foreign key (id_user)
        references account.users(id)
        on delete cascade
);
