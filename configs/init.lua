-- Initialize database
box.cfg {
    listen = 3301,
    wal_mode = 'write',  -- Enable persistence
    memtx_memory = 128 * 1024 * 1024,  -- 128MB memory limit
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


print('Tarantool database initialized!')