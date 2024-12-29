create table messages
(
    id         serial primary key,
    content    text,
    created_by int references users (id) on delete cascade,
    room_id     int references rooms (id) on delete cascade,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)