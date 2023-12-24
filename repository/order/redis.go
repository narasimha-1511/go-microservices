package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/narasimha-1511/go-microservices/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client;
}

func OrderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {

	data,err := json.Marshal(order);

	if err !=nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key:= OrderIDKey(order.OrderID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)

	if res.Err() != nil {
		txn.Discard()
		return fmt.Errorf("failed to insert order: %w", res.Err())
	}

	if err = txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if _, err = txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

func (r *RedisRepo) FindById(ctx context.Context, id uint64) (model.Order, error) {
	key:= OrderIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return model.Order{}, fmt.Errorf("order not found: %w", err)
	}else if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	var order model.Order

	err = json.Unmarshal([]byte(value), &order)

	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key:= OrderIDKey(id)

	txn:= r.Client.TxPipeline()

	res := txn.Del(ctx, key).Err()

	if errors.Is(res, redis.Nil) {
		txn.Discard()
		return fmt.Errorf("order not found: %w", res)
	}else if res != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order: %w", res)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

func (r *RedisRepo) UpdateById(ctx context.Context, order model.Order) error {

	data,err := json.Marshal(order);

	if err !=nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key:= OrderIDKey(order.OrderID)

	res := r.Client.SetXX(ctx, key, string(data), 0).Err()

	if errors.Is(res, redis.Nil) {
		return fmt.Errorf("order not found: %w", res)
	}else if res != nil {
		return fmt.Errorf("failed to update order: %w", res)
	}

	return nil
}

type FindAllPage struct {
	Size uint64
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context,page FindAllPage) (FindResult, error) {

	res:= r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()

	if res.Err() != nil {
		return FindResult{}, fmt.Errorf("failed to find orders: %w", res.Err())
	}

	if len(keys)==0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil;
	}

	xs , err := r.Client.MGet(ctx, keys...).Result()

	if err!=nil {
		return FindResult{}, fmt.Errorf("failed to find orders: %w", err)
	}

	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x:= x.(string)
		var order model.Order
		err = json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order: %w", err)
		}
		orders[i] = order
	}

	return FindResult{Orders: orders, Cursor: cursor}, nil;
}