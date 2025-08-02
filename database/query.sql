-- name: InsertSleep :one
insert into sleep (duration, efficiency, deep_pct, latency, num_awakenings)
values ($1, $2, $3, $4, $5)
returning *;

-- name: InsertDiet :one
insert into diet (meal, time, items)
values ($1, $2, $3)
returning *;

-- name: InsertMenstrual :one
insert into menstrual (cycle_day, pain_rating, stress_level, medication)
values ($1, $2, $3, $4)
returning *;