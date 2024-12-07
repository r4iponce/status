package config

import (
	"reflect"
	"testing"

	"gitlab.gnous.eu/ada/status/internal/cache"
	"gitlab.gnous.eu/ada/status/internal/log"
	"gitlab.gnous.eu/ada/status/internal/modules/http"
)

func TestToml(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		got, err := LoadToml("test_resources/valid.toml")
		if err != nil {
			t.Fatal(err)
		}

		want := Config{
			Listen: "localhost:3000",
			Log: log.Config{
				Level: "DEBUG",
				File:  "log",
			},
			Redis: cache.Config{
				Enabled:  true,
				Address:  "localhost:6379",
				Db:       0,
				User:     "test",
				Password: "Password123",
			},
			Check: Check{Interval: 1},
			Probe: []Target{{
				Name:        "example",
				Description: "Check https://example.org website",
				Module:      "http",
				Http: http.Config{
					Target: "https://example.org",
					Valid:  []int{200, 404},
				},
				Webhooks: Alerting{
					Enabled:  true,
					Username: "Status alert",
					Url:      "https://discord.com/api/webhooks/28357/verysecuretoken",
				},
			}},
		}

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}
