package testing

import (
	"fmt"
	"testing"
)

// TestNewAccount tests the NewAccount function.
func TestNewAccount(t *testing.T) {
	add, pubK := creatNewAddress()

	fmt.Println("add: ", add)
	fmt.Println("pubK: ", pubK)

	// address := "a04467fb2ef4dcb8d067cc39ae002cfc9423b6a1"
	// publicKeyHex := "cd56b13df64a1f5f4f6e5bcd13dd4109b63bf1f303634d90099a958afcdf4c7677c2e15e21ceb5e8b295f8887049e4efd4e343a8251d8ff3dfa9f63cc8a6c1f3"

	r, s, privateKey, sig := signTrx(add)

	pubKey := privateKey.PublicKey

	// fmt.Println("pubkey: ", pubKey)
	// fmt.Println(pubKey.Curve)
	// fmt.Println(pubKey.X)
	// fmt.Println(pubKey.Y)
	// fmt.Println("")

	// pubKeyA := TestPubkey(pubK)

	// fmt.Println("pubkeyA: ", pubKeyA)
	// fmt.Println(pubKeyA.Curve)
	// fmt.Println(pubKeyA.X)
	// fmt.Println(pubKeyA.Y)
	// fmt.Println("")

	verifyTrx(r, s, pubKey, sig, pubK)

}
