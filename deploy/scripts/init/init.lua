box.cfg {
    listen = 3301,
    wal_mode = 'write',
}

box.schema.space.create('database', { if_not_exists = true })

box.schema.user.create('my_user', {
    password = '123456',
    if_not_exists = true
})

box.schema.user.grant('my_user', 'read,write,execute', 'universe', nil, {if_not_exists = true})

box.space.database:format({
    { name = 'key', type = 'string' },
    { name = 'value', type = 'map' },
})

box.space.database:create_index('primary', {
    parts = { 'key' },
    if_not_exists = true,
})

