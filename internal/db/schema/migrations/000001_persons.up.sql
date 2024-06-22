create table if not exists persons(
    id uuid default gen_random_uuid() not null,
    full_name varchar not null,
    email varchar not null unique,
    avatar varchar not null,
    location varchar not null,
    metadata jsonb,
    created_at timestamp without time zone default timezone('utc'::text, now()) not null
);

create table if not exists persons_ranking(
    id uuid default gen_random_uuid() not null,
    person_id uuid not null unique,
    search_count int not null default 0,
    created_at timestamp without time zone default timezone('utc'::text, now()) not null
);

create or replace function update_person_ranking(pid uuid)
returns int as $$
declare
    current_search_count int;
begin
    select search_count
    into current_search_count
    from persons_ranking
    where person_id = pid
    for update skip locked;

    if not found then
        insert into persons_ranking(person_id, search_count)
        values (pid, 1)
        on conflict (person_id) do update
        set search_count = persons_ranking.search_count + 1
        returning search_count into current_search_count;
    else
        update persons_ranking
        set search_count = current_search_count + 1
        where person_id = pid
        returning search_count into current_search_count;
    end if;

    return current_search_count;
end;
$$ language plpgsql;
