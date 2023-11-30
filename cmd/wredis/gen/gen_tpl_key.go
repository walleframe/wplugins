package gen

func init() {
	registerTemplate("key", `
func (x *x{{.Name}}) Del(ctx context.Context) (ok bool, err error) {
	n, err := x.rds.Del(ctx, x.key).Result()
	if err != nil {
		return
	}
	ok = n == 1
	return
}

func (x *x{{.Name}}) Exists(ctx context.Context) (ok bool, err error) {
	n, err := x.rds.Exists(ctx, x.key).Result()
	if err != nil {
		return
	}
	ok = n == 1
	return
}

func (x *x{{.Name}}) Expire(ctx context.Context, expire time.Duration) (ok bool, err error) {
	return x.rds.Expire(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) ExpireNX(ctx context.Context, expire time.Duration) (ok bool, err error) {
	return x.rds.ExpireNX(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) ExpireXX(ctx context.Context, expire time.Duration) (ok bool, err error) {
	return x.rds.ExpireLT(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) ExpireGT(ctx context.Context, expire time.Duration) (ok bool, err error) {
	return x.rds.ExpireLT(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) ExpireLT(ctx context.Context, expire time.Duration) (ok bool, err error) {
	return x.rds.ExpireLT(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) ExpireAt(ctx context.Context, expire time.Time) (ok bool, err error) {
	return x.rds.ExpireAt(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) TTL(ctx context.Context) (time.Duration, error) {
	return x.rds.TTL(ctx, x.key).Result()
}

func (x *x{{.Name}}) PExpire(ctx context.Context, expire time.Duration) (ok bool, err error) {
	return x.rds.PExpire(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) PExpireAt(ctx context.Context, expire time.Time) (ok bool, err error) {
	return x.rds.PExpireAt(ctx, x.key, expire).Result()
}

func (x *x{{.Name}}) PExpireTime(ctx context.Context) (time.Duration, error) {
	return x.rds.PExpireTime(ctx, x.key).Result()
}

func (x *x{{.Name}}) PTTL(ctx context.Context) (time.Duration, error) {
	return x.rds.PTTL(ctx, x.key).Result()
}

func (x *x{{.Name}}) Persist(ctx context.Context) (ok bool, err error) {
	return x.rds.Persist(ctx, x.key).Result()
}

func (x *x{{.Name}}) Rename(ctx context.Context, newKey string) (err error) {
	return x.rds.Rename(ctx, x.key, newKey).Err()
}

func (x *x{{.Name}}) RenameNX(ctx context.Context, newKey string) (ok bool, err error) {
	return x.rds.RenameNX(ctx, x.key, newKey).Result()
}

func (x *x{{.Name}}) Type(ctx context.Context) (string, error) {
	return x.rds.Type(ctx, x.key).Result()
}
`)
}
