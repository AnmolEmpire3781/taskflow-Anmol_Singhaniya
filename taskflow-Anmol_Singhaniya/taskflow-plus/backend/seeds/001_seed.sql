INSERT INTO users (id, name, email, password)
VALUES (
        '11111111-1111-1111-1111-111111111111',
        'Test User',
        'test@example.com',
        '$2a$12$9NVoOKO9y1cwlnGY2xiio.UZrhBScs5EdNJp7l6QyUcJjW5fi2WYK'
    ) ON CONFLICT (email) DO NOTHING;
INSERT INTO projects (id, name, description, owner_id)
VALUES (
        '22222222-2222-2222-2222-222222222222',
        'TaskFlow Demo Project',
        'Seed project for reviewer testing',
        '11111111-1111-1111-1111-111111111111'
    ) ON CONFLICT (id) DO NOTHING;
INSERT INTO tasks (
        id,
        title,
        description,
        status,
        priority,
        project_id,
        assignee_id,
        creator_id,
        due_date
    )
VALUES (
        '33333333-3333-3333-3333-333333333331',
        'Create API spec',
        'Prepare backend API reference',
        'todo',
        'high',
        '22222222-2222-2222-2222-222222222222',
        '11111111-1111-1111-1111-111111111111',
        '11111111-1111-1111-1111-111111111111',
        CURRENT_DATE + INTERVAL '3 day'
    ),
    (
        '33333333-3333-3333-3333-333333333332',
        'Implement auth',
        'JWT and bcrypt setup',
        'in_progress',
        'medium',
        '22222222-2222-2222-2222-222222222222',
        '11111111-1111-1111-1111-111111111111',
        '11111111-1111-1111-1111-111111111111',
        CURRENT_DATE + INTERVAL '5 day'
    ),
    (
        '33333333-3333-3333-3333-333333333333',
        'Finish review notes',
        'Prepare for code review call',
        'done',
        'low',
        '22222222-2222-2222-2222-222222222222',
        NULL,
        '11111111-1111-1111-1111-111111111111',
        CURRENT_DATE + INTERVAL '7 day'
    ) ON CONFLICT (id) DO NOTHING;