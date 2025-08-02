create table sleep (
    id serial primary key,
    duration double precision,
    efficiency integer,
    deep_pct integer,
    latency integer,
    num_awakenings integer
);

create table diet (
    id serial primary key,
    meal text,
    items text[],
    time timestamp
);

create table menstrual (
    id serial primary key,
    cycle_day integer,
    pain_rating integer,
    stress_level integer,
    medication text[]
);