package caches

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/adamluzsi/frameless/reflects"
	"github.com/adamluzsi/frameless/resources"
	"github.com/toggler-io/toggler/usecases"
	"github.com/go-redis/redis"
	"time"
)

func NewRedis(connstr string, storage usecases.Storage) (*Redis, error) {
	redisClientOpt, err := redis.ParseURL(connstr)
	if err != nil {
		return nil, err
	}
	return &Redis{Storage: storage, client: redis.NewClient(redisClientOpt)}, nil
}

// TODO provide caching for every Storage contract (function)
type Redis struct {
	usecases.Storage
	client *redis.Client
	ttl    time.Duration
}

func (r *Redis) SetTimeToLiveForValuesToCache(d time.Duration) error {
	r.ttl = d
	return nil
}

func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		return err
	}
	return r.Storage.Close()
}

func (r *Redis) FindByID(ctx context.Context, ptr interface{}, ID string) (bool, error) {
	key := fmt.Sprintf(`%s#%s`, reflects.FullyQualifiedName(ptr), ID)
	reply := r.client.Get(key)

	err := reply.Err()

	if err != nil && err != redis.Nil {
		return false, err
	}

	if err == redis.Nil {
		found, err := r.Storage.FindByID(ctx, ptr, ID);
		if err != nil {
			return false, err
		}

		bs, err := r.marshal(ptr)

		if err != nil {
			return false, err
		}

		if err := r.client.Set(key, bs, r.ttl).Err(); err != nil {
			return false, err
		}

		return found, err
	} else {
		bs, err := reply.Bytes()

		if err != nil {
			return false, err
		}

		return true, r.unmarshal(ptr, bs)
	}
}

func (r *Redis) DeleteByID(ctx context.Context, Type interface{}, ID string) error {
	if err := r.invalidate(Type, ID); err != nil {
		return err
	}
	return r.Storage.DeleteByID(ctx, Type, ID)
}

func (r *Redis) DeleteAll(ctx context.Context, Type interface{}) error {
	keysRes := r.client.Keys(reflects.FullyQualifiedName(Type) + `*`)

	if err := keysRes.Err(); err != nil && err != redis.Nil {
		return err
	}

	keys, err := keysRes.Result()
	if err != nil && err != redis.Nil {
		return err
	}

	for _, key := range keys {
		if err := r.invalidateKey(key); err != nil {
			return err
		}
	}

	return r.Storage.DeleteAll(ctx, Type)
}

func (r *Redis) Update(ctx context.Context, ptr interface{}) error {
	id, found := resources.LookupID(ptr)
	if found {
		if err := r.invalidateKey(r.cacheKeyOfObject(ptr, id)); err != nil {
			return err
		}
	}
	return r.Storage.Update(ctx, ptr)
}

//
//func (*Redis) FindAll(ctx context.Context, Type interface{}) frameless.Iterator {
//	panic("implement me")
//}
//
//func (*Redis) FindReleaseFlagByName(ctx context.Context, name string) (*rollouts.FeatureFlag, error) {
//	panic("implement me")
//}
//
//func (*Redis) FindReleasePilotByReleaseFlagAndDeploymentEnvironmentAndExternalID(ctx context.Context, FlagID, ExternalPilotID string) (*rollouts.Pilot, error) {
//	panic("implement me")
//}
//
//func (*Redis) FindReleasePilotsByReleaseFlag(ctx context.Context, ff *rollouts.FeatureFlag) frameless.Iterator {
//	panic("implement me")
//}
//
//func (*Redis) FindTokenBySHA512Hex(ctx context.Context, sha512hex string) (*security.Token, error) {
//	panic("implement me")
//}
//

//--------------------------------------------------------------------------------------------------------------------//

func (r *Redis) lookup(ptr interface{}, ID string) (bool, error) {
	return false, nil
}

func (r *Redis) marshal(ptr interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(ptr)
	return buf.Bytes(), err
}

func (r *Redis) unmarshal(ptr interface{}, bs []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(bs))
	return decoder.Decode(ptr)
}

func (r *Redis) invalidate(t interface{}, id string) error {
	return r.invalidateKey(r.cacheKeyOfObject(t, id))
}

func (r *Redis) cacheKeyOfObject(entity interface{}, id string) string {
	return fmt.Sprintf(`%s#%s`, reflects.FullyQualifiedName(entity), id)
}

func (r *Redis) invalidateKey(key string) error {
	reply := r.client.Del(key)
	if err := reply.Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}
