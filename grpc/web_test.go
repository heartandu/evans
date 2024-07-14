package grpc_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/ktr0731/evans/grpc"
)

func TestWebClient(t *testing.T) {
	client, err := grpc.NewWebClient("", false, false, "", "", "", nil)
	if err != nil {
		t.Errorf("NewWebClient should not return errors: %v", err)
		return
	}

	t.Run("Invoke returns an error if FQRN is invalid", func(t *testing.T) {
		_, _, err := client.Invoke(context.Background(), "invalid-fqrn", nil, nil)
		if err == nil {
			t.Errorf("expected an error, but got nil")
		}
	})
	t.Run("NewClientStream returns an error if FQRN is invalid", func(t *testing.T) {
		_, err := client.NewClientStream(context.Background(), nil, "invalid-fqrn")
		if err == nil {
			t.Errorf("expected an error, but got nil")
		}
	})
	t.Run("NewServerStream returns an error if FQRN is invalid", func(t *testing.T) {
		_, err := client.NewServerStream(context.Background(), nil, "invalid-fqrn")
		if err == nil {
			t.Errorf("expected an error, but got nil")
		}
	})
	t.Run("NewBidiStream returns an error if FQRN is invalid", func(t *testing.T) {
		_, err := client.NewBidiStream(context.Background(), nil, "invalid-fqrn")
		if err == nil {
			t.Errorf("expected an error, but got nil")
		}
	})
}

func TestNewClient(t *testing.T) {
	certPath := func(s ...string) string {
		return filepath.Join(append([]string{"testdata", "cert"}, s...)...)
	}

	cases := map[string]struct {
		addr          string
		useReflection bool
		useTLS        bool
		cacert        string
		cert          string
		certKey       string

		hasErr bool
		err    error
	}{
		"certKey is missing":                      {useTLS: true, cert: "foo", err: grpc.ErrMutualAuthParamsAreNotEnough},
		"cert is missing":                         {useTLS: true, certKey: "bar", err: grpc.ErrMutualAuthParamsAreNotEnough},
		"certKey is missing, but useTLS is false": {cert: "foo"},
		"cert is missing, but useTLS is false":    {certKey: "foo"},
		"enable server TLS":                       {useTLS: true},
		"enable server TLS with a trusted CA":     {useTLS: true, cacert: certPath("rootCA.pem")},
		"enable mutual TLS":                       {useTLS: true, cert: certPath("localhost.pem"), certKey: certPath("localhost-key.pem")},
		"enable mutual TLS with a trusted CA":     {useTLS: true, cacert: certPath("rootCA.pem"), cert: certPath("localhost.pem"), certKey: certPath("localhost-key.pem")},
		"invalid cacert file path":                {useTLS: true, cacert: "fooCA.pem", hasErr: true},
		"invalid cert and key file path":          {useTLS: true, cert: "foo.pem", certKey: "foo-key.pem", hasErr: true},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			_, err := grpc.NewWebClient(c.addr, c.useReflection, c.useTLS, c.cacert, c.cert, c.certKey, nil)
			if c.err != nil {
				if err == nil {
					t.Fatalf("NewClient must return an error, but got nil")
				}
				if !errors.Is(err, c.err) {
					t.Errorf("expected: '%s', but got '%s'", c.err, err)
				}

				return
			} else if c.hasErr {
				if err == nil {
					t.Fatalf("NewClient must return an error, but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewClient must not return an error, but got '%s'", err)
			}
		})
	}
}
