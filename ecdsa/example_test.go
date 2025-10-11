package ecdsa_test

import (
	"fmt"

	myecdsa "github.com/masatomo57/golang-oreore-comparable/ecdsa"
)

func Example() {
	priv, err := myecdsa.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	message := []byte("hello ECDSA")
	sig, err := myecdsa.Sign(priv, message)
	if err != nil {
		panic(err)
	}

	pub, err := myecdsa.PublicKey(priv)
	if err != nil {
		panic(err)
	}

	fmt.Println(myecdsa.Verify(pub, message, sig))
	// Output: true
}
