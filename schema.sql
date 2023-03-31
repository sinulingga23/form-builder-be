create database form_builder;

create schema partner;

create table partner.m_partner (
    id uuid not null primary key,
    name varchar(150) not null,
    description varchar(150) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp null
);

create table partner.m_form (
    id uuid not null primary key,
    code varchar(30) not null unique,
    name varchar(150) not null,
    m_partner_id uuid not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp null,
    foreign key (m_partner_id) references partner.m_partner (id)
    on update cascade on delete restrict
);

create table partner.m_field_type (
    id uuid not null primary key,
    name varchar(150) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp null
);