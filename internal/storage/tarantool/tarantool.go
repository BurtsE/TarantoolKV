package tarantool

import (
	"TarantoolKV/internal/application/core/domain"
	"context"
	"github.com/tarantool/go-iproto"
	"github.com/tarantool/go-tarantool/v2"
	"log"
	"time"
)

type TarantoolDB struct {
	conn *tarantool.Connection
}

func NewStorage() TarantoolDB {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  "127.0.0.1:3301",
		User:     "guest",
		Password: "",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}
	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		log.Fatal(err)
	}
	return TarantoolDB{
		conn: conn,
	}
}

func (t *TarantoolDB) Create(ctx context.Context, entity domain.Entity) error {
	data := []interface{}{
		entity.Key,
		entity.Value,
	}
	data, err := t.conn.Do(
		tarantool.NewInsertRequest("database").Tuple(data),
	).Get()
	if err != nil {
		if tarantoolError, ok := err.(tarantool.Error); ok && tarantoolError.Code == iproto.ER_TUPLE_FOUND {
			return domain.ErrKeyExists
		}
	}
	return err
}

func (t *TarantoolDB) Update(ctx context.Context, entity domain.Entity) error {
	data, err := t.conn.Do(
		tarantool.NewUpdateRequest("database").
			Key(tarantool.StringKey{entity.Key}).
			Operations(tarantool.NewOperations().Assign(1, entity.Value)),
	).Get()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return domain.ErrKeyNotFound
	}
	return err
}
func (t *TarantoolDB) Get(ctx context.Context, key string) (domain.Entity, error) {
	var result [][]interface{}
	err := t.conn.Do(
		tarantool.NewSelectRequest("database").
			Iterator(tarantool.IterEq).
			Key(tarantool.StringKey{key}),
	).GetTyped(&result)
	if err != nil {
		return domain.Entity{}, err
	}
	if len(result) == 0 {
		return domain.Entity{}, domain.ErrKeyNotFound
	}
	entity := convertToDomain(result)
	return entity, nil
}

func (t *TarantoolDB) Delete(ctx context.Context, key string) error {
	data, err := t.conn.Do(
		tarantool.NewDeleteRequest("database").Key(tarantool.StringKey{key}),
	).Get()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return domain.ErrKeyNotFound
	}
	return err
}

func (t *TarantoolDB) Shutdown() error {
	return t.conn.CloseGraceful()
}
