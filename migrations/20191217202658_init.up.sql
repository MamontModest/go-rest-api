CREATE table users
(
    login varchar(63) primary key ,
    password bigint
);
CREATE table recipe
(
    recipeId serial primary key ,
    name varchar(31),
    title varchar(255)
);
create table ingredient
(
    recipeId integer REFERENCES recipe (recipeId) ON DELETE CASCADE,
    ingredient varchar(31),
    primary key (recipeId, ingredient)
);
create table step
(
    recipeId integer REFERENCES recipe (recipeId) ON DELETE CASCADE,
    stepNumber integer,
    description varchar(255),
    timeDuration integer,
    primary key (recipeId, stepNumber)
);