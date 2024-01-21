package testing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

func decompressPublicKey(compressedKey []byte) ([]byte, error) {
	// Check if the compressed key is valid
	if len(compressedKey) != 33 {
		return nil, errors.New("invalid compressed public key length")
	}

	curve := elliptic.P256() // Replace with the appropriate elliptic curve

	// Determine the sign byte (0x02 or 0x03)
	signByte := compressedKey[0]
	if signByte != 0x02 && signByte != 0x03 {
		return nil, errors.New("invalid sign byte")
	}

	// Parse the X-coordinate (skip the sign byte)
	xBytes := compressedKey[1:]
	x := new(big.Int).SetBytes(xBytes)

	// Calculate the Y-coordinate based on the X-coordinate and the curve equation
	y := new(big.Int)
	y.ModSqrt(y.Mod(y.Mul(x, x), curve.Params().P), curve.Params().P)

	// Construct the full uncompressed public key
	uncompressedKey := append([]byte{0x04}, x.Bytes()...)
	uncompressedKey = append(uncompressedKey, y.Bytes()...)

	return uncompressedKey, nil
}

// VerifyTransactionSignature verifies the authenticity and integrity of a transaction using the public key.
func VerifyTransactionSignature(pubKeyBytes []byte, signatureHex string, r *big.Int, s *big.Int, pubKey ecdsa.PublicKey, transactionDataHex string) (bool, error) {
	// Parse the public key
	// CpublicKeyBytes, err := hex.DecodeString(publicKeyHex)
	// fmt.Println("CpublicKeyBytes: ", CpublicKeyBytes)

	// publicKeyBytes, err := decompressPublicKey(CpublicKeyBytes)
	// fmt.Println("publicKeyBytes: ", publicKeyBytes)

	// if err != nil {
	// 	return false, err
	// }

	// // publicKey := &ecdsa.PublicKey{
	// // 	Curve: elliptic.P256(), // Use the appropriate elliptic curve
	// // 	X:     new(big.Int),
	// // 	Y:     new(big.Int),
	// // }
	// // publicKey.X.SetBytes(publicKeyBytes[1:33])
	// // publicKey.Y.SetBytes(publicKeyBytes[33:])

	// publicKey := new(ecdsa.PublicKey)
	// publicKey.Curve = elliptic.P256() // Replace with the appropriate elliptic curve
	// publicKey.X = new(big.Int).SetBytes(publicKeyBytes[1:33])
	// publicKey.Y = new(big.Int).SetBytes(publicKeyBytes[33:65])

	// Parse the signature
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, err
	}

	// Calculate the hash of the transaction data
	transactionData, err := hex.DecodeString(transactionDataHex)
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256(transactionData)

	// transactionData = []byte("data to sign")
	// Verify the signature
	r_ := new(big.Int).SetBytes(signatureBytes[:32])
	s_ := new(big.Int).SetBytes(signatureBytes[32:])
	// r and s components as decimal values
	// rDecimal := new(big.Int).SetBytes([]byte{105, 173, 65, 93, 210, 78, 245, 62, 247, 95, 131, 50, 100, 83, 246, 160, 81, 103, 230, 223, 244, 229, 128, 167, 237, 77, 30, 109, 158, 125, 100, 237})
	// sDecimal := new(big.Int).SetBytes([]byte{181, 218, 90, 6, 12, 35, 254, 214, 189, 221, 100, 35, 41, 179, 238, 5, 63, 13, 0, 89, 2, 48, 119, 102, 168, 151, 235, 176, 125, 93, 31, 239})

	// fmt.Println("r: ", r)
	// fmt.Println("s: ", s)

	// transactionData1 := []byte("data to sign")
	// hash1 := sha256.Sum256(transactionData1)

	if err != nil {
		fmt.Println("Error decoding public key hex:", err)
		return false, err
	}

	// Create big integers from the bytes

	xBytes := pubKeyBytes[:32] // Assuming X component is 32 bytes
	yBytes := pubKeyBytes[32:] // Assuming Y component is 32 bytes

	// Create an ecdsa.PublicKey
	publicKey__ := ecdsa.PublicKey{
		Curve: elliptic.P256(), // Use the appropriate elliptic curve
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}
	fmt.Println("publicKey__: ", publicKey__)

	tF := ecdsa.Verify(&publicKey__, hash[:], r_, s_)
	if tF {
		return true, nil
	}

	return false, errors.New("signature verification failed")
}

func verifyTrx(r *big.Int, s *big.Int, pubKey ecdsa.PublicKey, sig string, pubKeyBytes []byte) {
	// Example usage
	signatureHex := hex.EncodeToString([]byte(sig))

	transactionData := []byte("data to sign1")
	transactionDataHex := hex.EncodeToString(transactionData)

	isValid, err := VerifyTransactionSignature(pubKeyBytes, signatureHex, r, s, pubKey, transactionDataHex)
	if err != nil {
		fmt.Println(err)
	} else if isValid {
		fmt.Println("Transaction signature is valid.")
	} else {
		fmt.Println("Transaction signature is invalid.")
	}
}

func TestPubkey(pubKeyBytes []byte) ecdsa.PublicKey {
	// publicKeyHex_ := "cd56b13df64a1f5f4f6e5bcd13dd4109b63bf1f303634d90099a958afcdf4c7677c2e15e21ceb5e8b295f8887049e4efd4e343a8251d8ff3dfa9f63cc8a6c1f3"
	xBytes := pubKeyBytes[:32] // Assuming X component is 32 bytes
	yBytes := pubKeyBytes[32:] // Assuming Y component is 32 bytes

	// Create an ecdsa.PublicKey
	publicKey__ := ecdsa.PublicKey{
		Curve: elliptic.P256(), // Use the appropriate elliptic curve
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}
	fmt.Println("publicKey__: ", publicKey__)

	return publicKey__
}
