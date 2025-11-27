package main

import (
	"fmt"

	encrypt "github.com/tangthinker/encrypt-conn-tools/pkg"
)

func main() {
	pubKey, privKey := encrypt.GenerateKeyPairECDSA()
	fmt.Println("pubKey:", pubKey)
	fmt.Println("privKey:", privKey)
}
