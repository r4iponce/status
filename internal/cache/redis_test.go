package cache

import (
	"testing"
)

func TestRedisConnection(t *testing.T) {
	t.Parallel()

	c := Config{
		Address: "localhost:6379",
		Db:      0,
	}

	c.Connect()
	err := Ping()
	if err != nil {
		t.Fatal(err)
	}
}
