package gen

func init(){
	registerTemplate("redis_lock", `{{ $Name := .Name}} 
func (x *x{{$Name}}) Lock(ctx context.Context, expiration time.Duration) (lockID string, err error) {
	var cmd *redis.BoolCmd
	lockID = uuid.NewString() {{Import "github.com/google/uuid" "NewString"}}
	switch expiration {
	case 0:
		// Use old 'SETNX' to support old Redis versions.
		cmd = redis.NewBoolCmd(ctx, "setnx", x.key, lockID)
	case {{.SvcPkg}}.KeepTTL:
		cmd = redis.NewBoolCmd(ctx, "set", x.key, lockID, "keepttl", "nx")
	default:
		if {{.SvcPkg}}.UsePrecise(expiration) {
			cmd = redis.NewBoolCmd(ctx, "set", x.key, lockID, "px", {{.SvcPkg}}.FormatMs(ctx, expiration), "nx")
		} else {
			cmd = redis.NewBoolCmd(ctx, "set", x.key, lockID, "ex", {{.SvcPkg}}.FormatSec(ctx, expiration), "nx")
		}
	}
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		return
	}
	if !cmd.Val() {
		err = {{.SvcPkg}}.ErrLockFailed
	}
	return
}

func (x *x{{$Name}}) UnLock(ctx context.Context, lockID string) (ok bool, err error) {
	cmd := redis.NewIntCmd(ctx, "evalsha", {{.SvcPkg}}.LockerScriptUnlock.Hash, "1", x.key, lockID)
	err = x.rds.Process(ctx, cmd)
	if err != nil {
		if !redis.HasErrorPrefix(err, "NOSCRIPT") {
			return
		}
		cmd = redis.NewIntCmd(ctx, "eval", {{.SvcPkg}}.LockerScriptUnlock.Script, "1", x.key, lockID)
		err = x.rds.Process(ctx, cmd)
		if err != nil {
			return
		}
	}
	ok = cmd.Val() > 0
	return
}

func (x *x{{$Name}}) LockFunc(ctx context.Context, expiration time.Duration) (unlock func(ctx context.Context), err error) {
	lockID, err := x.Lock(ctx, expiration)
	if err != nil {
		return
	}
	unlock = func(ctx context.Context) {
		x.UnLock(ctx, lockID)
	}
	return
}
`)
}
