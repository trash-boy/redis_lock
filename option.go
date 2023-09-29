package learn_redis_lock

const (
	DefaultIdleTimeoutSeconds = 10

	DefaultMaxActive = 100

	DefaultMaxIdle = 20

	DefaultLockExpireSeconds = 30

	WatchDogWorkStepSeconds = 10
)

type ClientOptions struct {
	maxIdle int
	idleTimeoutSeconds int
	maxActive int
	wait bool

	network string
	address string
	password string
}

type ClientOption func(c *ClientOptions)

func WithMaxIdle(maxIdle int)ClientOption{
	return func(c *ClientOptions) {
		c.maxIdle = maxIdle
	}
}

func WithIdleTimeoutSeconds(idleTimeoutSeconds int)ClientOption{
	return func(c *ClientOptions) {
		c.idleTimeoutSeconds = idleTimeoutSeconds
	}
}

func WithMaxActive(maxActive int)ClientOption{
	return func(c *ClientOptions) {
		c.maxActive = maxActive
	}
}

func WithWaitMode()ClientOption{
	return func(c *ClientOptions) {
		c.wait = true
	}
}

func repairClient(c *ClientOptions){
	if c.maxIdle < 0{
		c.maxIdle = DefaultMaxIdle
	}

	if c.idleTimeoutSeconds < 0{
		c.idleTimeoutSeconds = DefaultIdleTimeoutSeconds
	}

	if c.maxActive < 0{
		c.maxActive = DefaultMaxActive
	}
}

type LockOptions struct {
	isBlock bool
	blockWaitingSeconds int64
	expireSeconds int64
	watchDogMode bool
}

type LockOption func(*LockOptions)

func WithBlock()LockOption{
	return func(options *LockOptions) {
		options.isBlock = true
	}
}

func WithBlockWaitingSeconds(waitingSeconds int64)LockOption{
	return func(options *LockOptions) {
		options.blockWaitingSeconds = waitingSeconds
	}
}

func WithExpireSeconds(expireSeconds int64)LockOption{
	return func(options *LockOptions) {
		options.expireSeconds = expireSeconds
	}
}

func repairLock(options *LockOptions){
	if options.isBlock && options.blockWaitingSeconds <= 0{
		options.blockWaitingSeconds = 5
	}
	if options.expireSeconds > 0{
		return
	}

	options.expireSeconds = DefaultLockExpireSeconds
	options.watchDogMode = true
}
