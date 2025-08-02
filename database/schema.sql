create table sleep (
    id serial primary key,
    date date not null,
    duration double precision, -- hours
    quality integer, -- 1 to 10 scale
    disruptions text,
    notes text
);

create table diet (
    id serial primary key,
    meal text, -- breakfast, lunch, dinner, snack etc.
    date date not null,
    items text[], -- also mention ingredients
    notes text
);

create table menstrual (
    id serial primary key,
    period_event text, -- start, end, ovulation, etc.
    date date not null,
    flow_level text, -- light, medium, heavy
    notes text
);


create table symptoms (
    id serial primary key,
    date date not null,
    nausea integer, -- 1 to 10 scale
    fatigue integer, -- 1 to 10 scale
    pain integer, -- 1 to 10 scale
    notes text
);