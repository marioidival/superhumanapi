-- name: PersonByEmail :one
select id, full_name, avatar, "location", metadata from persons where email = $1;

-- name: CreatePerson :one
insert into persons (full_name, email, avatar, location, metadata)
values ($1, $2, $3, $4, $5) returning *;

-- name: UpdatePersonRanking :one
select update_person_ranking($1);

-- name: PopularPersons :many
select persons.full_name, persons.avatar, persons.location, persons_ranking.search_count
from persons_ranking
join persons on persons_ranking.person_id = persons.id
order by persons_ranking.search_count desc
limit 10;

-- name: SomoeShow :many
select * from persons;