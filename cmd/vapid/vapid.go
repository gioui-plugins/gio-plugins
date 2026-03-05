package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	curve := elliptic.P256()

	private, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return
	}

	public := elliptic.Marshal(curve, x, y)

	// Convert to base64
	publicKey := base64.RawURLEncoding.EncodeToString(public)
	privateKey := base64.RawURLEncoding.EncodeToString(private)

	fmt.Printf("Public Key (use in Gio-Plugins): %s\n", publicKey)
	fmt.Printf("Private Key (use in your server): %s\n", privateKey)
}
