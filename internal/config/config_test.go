package config

import (
	"errors"
	"testing"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		c, err := LoadToml("test_resources/valid.toml")
		if err != nil {
			t.Fatal(err)
		}

		got := c.Verify()

		if err != nil {
			t.Fatalf("want nil, got %v", got)
		}
	})

	t.Run("Invalid level", func(t *testing.T) {
		t.Parallel()

		c, err := LoadToml("test_resources/invalid_level.toml")
		if err != nil {
			t.Fatal(err)
		}

		got := c.Verify()
		want := errLogLevel

		if !errors.Is(got, want) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}
