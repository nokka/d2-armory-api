package env

import "testing"

func TestBool(t *testing.T) {
	env := &Client{func(key string) (string, bool) {
		value, found := map[string]string{
			"foo": "true",
			"bar": "invalid",
		}[key]

		return value, found
	}}

	for _, tt := range []struct {
		key         string
		fallback    bool
		expectError bool
		want        bool
	}{
		{"foo", false, false, true},
		{"bar", false, true, false},
		{"unknown", false, false, false},
	} {
		got, err := env.Bool(tt.key, tt.fallback)

		if tt.expectError && err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !tt.expectError && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != tt.want {
			t.Fatalf(`env.Bool(%q, %v) = %v, want %v`, tt.key, tt.fallback, got, tt.want)
		}
	}
}

func TestInt(t *testing.T) {
	env := &Client{func(key string) (string, bool) {
		value, found := map[string]string{
			"foo":  "123",
			"fail": "not_integer",
		}[key]

		return value, found
	}}

	for _, tt := range []struct {
		key         string
		fallback    int
		expectError bool
		want        int
	}{
		{"foo", 1, false, 123},
		{"fail", 2, true, 0},
		{"unknown", 3, false, 3},
	} {
		got, err := env.Int(tt.key, tt.fallback)

		if tt.expectError && err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !tt.expectError && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != tt.want {
			t.Fatalf(`env.Int(%q, %d) = %d, want %d`, tt.key, tt.fallback, got, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	env := &Client{func(key string) (string, bool) {
		value, found := map[string]string{
			"foo": "abc",
		}[key]

		return value, found
	}}

	for _, tt := range []struct {
		key      string
		fallback string
		want     string
	}{
		{"foo", "fallback", "abc"},
		{"unknown", "fallback", "fallback"},
	} {
		if got := env.String(tt.key, tt.fallback); got != tt.want {
			t.Fatalf(`env.String(%q, %q) = %q, want %q`, tt.key, tt.fallback, got, tt.want)
		}
	}
}
