package securecookie

import (
	"bytes"
	"errors"
	"testing"
)

var key32 = []byte("0123456789abcdef0123456789abcdef")

func TestSealUnseal_RoundTrip(t *testing.T) {
	plaintext := []byte(`{"access_token":"a","refresh_token":"r"}`)
	sealed, err := Seal(key32, plaintext)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	got, err := Unseal(key32, sealed)
	if err != nil {
		t.Fatalf("unseal: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("round trip mismatch: got %q want %q", got, plaintext)
	}
}

func TestUnseal_WrongKey(t *testing.T) {
	sealed, err := Seal(key32, []byte("secret"))
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	otherKey := []byte("ffffffffffffffffffffffffffffffff")
	if _, err := Unseal(otherKey, sealed); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUnseal_Tampered(t *testing.T) {
	sealed, err := Seal(key32, []byte("secret"))
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	tampered := sealed[:len(sealed)-1] + "x"
	if _, err := Unseal(key32, tampered); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestSeal_ShortKeyDerived(t *testing.T) {
	// Any non-empty key works (it is hashed to 32 bytes).
	sealed, err := Seal([]byte("short"), []byte("x"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := Unseal([]byte("short"), sealed)
	if err != nil || string(got) != "x" {
		t.Fatalf("round trip failed: got %q err %v", got, err)
	}
}

func TestSeal_EmptyKey(t *testing.T) {
	if _, err := Seal([]byte(""), []byte("x")); err == nil {
		t.Fatal("expected error for empty key")
	}
}
