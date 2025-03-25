-- Initialize database
box.cfg {
    listen = 3301,
    wal_mode = 'write',  -- Enable persistence
}

-- Create a space (table) for users
box.schema.space.create('database', { if_not_exists = true })
-- Create admin user with password
box.schema.user.create('my_user', {
    password = '123456',
    if_not_exists = true
})

-- Grant permissions
box.schema.user.grant('my_user', 'read,write,execute', 'universe', nil, {if_not_exists = true})

-- Define field formats (optional but recommended)
box.space.database:format({
    { name = 'key', type = 'string' },
    { name = 'value', type = 'map' },
})


-- Create a primary index
box.space.database:create_index('primary', {
    parts = { 'key' },
    if_not_exists = true,
})

print('Tarantool database initialized!')