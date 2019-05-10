package redis

import ( 
	"github.com/go-redis/redis"
)

// Subscribe to redis channels
func Subscribe(client *redis.Client, channels ...string) <-chan *redis.Message {
	pubsub := client.Subscribe(channels...)
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	return pubsub.Channel()
}
