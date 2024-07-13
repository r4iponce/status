package config

import (
	"reflect"
	"testing"

	"go.ada.wf/status/internal/cache"
	"go.ada.wf/status/internal/log"
	"go.ada.wf/status/internal/modules/http"
	"go.ada.wf/status/internal/probe"
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
			Probe: []probe.Target{{
				Name:        "example",
				Description: "Check https://example.org website",
				Module:      "http",
				Http: http.Config{
					Target: "https://example.org",
					Valid:  []int{200, 404},
				},
			}},
		}

		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}
