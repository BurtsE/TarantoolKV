-- Initialize database
box.cfg {
    listen = 3301,
    wal_mode = 'write',  -- Enable persistence
}

-- Create a space (table) for users
box.schema.space.create('database', { if_not_exists = true })

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

-- Insert sample data
box.space.database:insert({'Alice', {foo = "bar"}})

print('Tarantool database initialized!')