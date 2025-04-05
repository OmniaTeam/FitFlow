CREATE TYPE gender_enum AS ENUM ('male', 'female');
CREATE TYPE difficulty_enum AS ENUM ('easy', 'medium', 'hard');
CREATE TYPE workout_status_enum AS ENUM ('planned', 'completed', 'skipped');

create table public.users
(
    id         serial  primary key,
    last_name  varchar(100)                                       not null,
    first_name varchar(100)                                       not null,
    gender     gender_enum,
    birthday   date,
    weight     numeric(5, 2),
    height     numeric(5, 2),
    google_id  varchar(100),
    vk_id      varchar(100),
    created_at timestamp default CURRENT_TIMESTAMP,
    email      varchar(255)                                       not null,
    roles      text[]
);

alter table public.users
    owner to omnia;

create table public.user_stats
(
    id      serial
        primary key,
    user_id integer not null
        references public.users,
    date    date    not null,
    weight  numeric(5, 2),
    height  numeric(5, 2),
    unique (user_id, date)
);

alter table public.user_stats
    owner to omnia;

create index idx_user_stats_user_id
    on public.user_stats (user_id);

create table public.programs
(
    id          serial               primary key,
    user_id     integer  unique      not null
        constraint program_user_id_fkey
            references public.users,
    name        varchar(100)                                          not null,
    description text,
    created_at  timestamp default CURRENT_TIMESTAMP
);

alter table public.programs
    owner to omnia;

create table public.exercises
(
    id          serial               primary key,
    name        varchar(100)                                           not null,
    description text,
    type        varchar(50) not null ,
    difficulty  difficulty_enum not null ,
    equipment   varchar(100),
    video_url   varchar(255),
    photo_urls  text[],
    created_at  timestamp default CURRENT_TIMESTAMP
);

alter table public.exercises
    owner to omnia;

create table public.muscles
(
    id          serial               primary key,
    name        varchar(100)                                         not null,
    description text,
    photo       varchar(255),
    created_at  timestamp default CURRENT_TIMESTAMP
);

alter table public.muscles
    owner to omnia;

create table public.workouts
(
    id          serial                        primary key,
    program_id  integer
        constraint workout_program_id_fkey
            references public.programs,
    user_id     integer                                                         not null
        constraint workout_user_id_fkey
            references public.users,
    name        varchar(100)                                                    not null,
    description text,
    date_time   timestamp,
    status      workout_status_enum default 'planned'::workout_status_enum,
    notes       text,
    created_at  timestamp           default CURRENT_TIMESTAMP
);

alter table public.workouts
    owner to omnia;

create index idx_workout_user_id
    on public.workouts (user_id);

create index idx_workout_program_id
    on public.workouts (program_id);

create table public.workout_exercise
(
    id                 serial
        primary key,
    workout_id         integer not null
        references public.workouts,
    exercise_id        integer not null
        references public.exercises,
    order_number       integer not null,
    unique (workout_id, exercise_id, order_number)
);

alter table public.workout_exercise
    owner to omnia;

create index idx_workout_exercise_workout_id
    on public.workout_exercise (workout_id);

create table public.exercise_set
(
    id                  serial
        primary key,
    workout_exercise_id integer not null
        references public.workout_exercise,
    set_number          integer not null,
    planned_sets        integer,
    planned_reps        integer,
    planned_weight      numeric(5, 2),
    completed_sets      integer,
    completed_reps      integer,
    completed_weight    numeric(5, 2),
    is_completed        boolean default false,
    unique (workout_exercise_id, set_number)
);

alter table public.exercise_set
    owner to omnia;

create index idx_exercise_set_workout_exercise_id
    on public.exercise_set (workout_exercise_id);

create table public.exercise_muscle
(
    exercise_id      integer       not null
        constraint exercise_muscle_groups_exercise_id_fkey
            references public.exercises,
    muscle_id        integer       not null
        constraint exercise_muscle_groups_muscle_id_fkey
            references public.muscles,
    muscles_involved numeric(3, 2) not null,
    constraint exercise_muscle_groups_pkey
        primary key (exercise_id, muscle_id)
);

alter table public.exercise_muscle
    owner to omnia;

create index idx_exercise_muscle_exercise_id
    on public.exercise_muscle (exercise_id);

