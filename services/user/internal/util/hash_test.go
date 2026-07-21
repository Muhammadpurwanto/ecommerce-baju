package util

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword123"

	// 1. Uji proses Hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Gagal melakukan hashing password: %v", err)
	}

	if hash == "" {
		t.Error("Diharapkan hasil hash tidak kosong")
	}

	if hash == password {
		t.Error("Diharapkan hasil hash berbeda dengan password asli")
	}

	// 2. Uji verifikasi password yang benar
	if !CheckPasswordHash(password, hash) {
		t.Error("Diharapkan CheckPasswordHash mengembalikan true untuk password yang benar")
	}

	// 3. Uji verifikasi password yang salah
	if CheckPasswordHash("wrongpassword123", hash) {
		t.Error("Diharapkan CheckPasswordHash mengembalikan false untuk password yang salah")
	}
}
