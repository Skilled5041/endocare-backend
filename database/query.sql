-- name: InsertSleep :one
insert into sleep (date, duration, quality, disruptions, notes)
values ($1, $2, $3, $4, $5)
returning *;

-- name: InsertDiet :one
insert into diet (meal, date, items, notes)
values ($1, $2, $3, $4)
returning *;

-- name: InsertMenstrual :one
insert into menstrual (period_event, date, flow_level, notes)
values ($1, $2, $3, $4)
returning *;

-- name: InsertSymptoms :one
insert into symptoms (date, nausea, fatigue, pain, notes)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetAllSleep :many
select * from sleep;

-- name: GetAllDiet :many
select * from diet;

-- name: GetAllMenstrual :many
select * from menstrual;

-- name: GetAllSymptoms :many
select * from symptoms;
