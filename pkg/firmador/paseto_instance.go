package firmador

import (
	"fmt"
	"os"

	"aidanwoods.dev/go-paseto"
)

var SecretKey paseto.V4SymmetricKey

// Inicializar la clave secreta para firmar y verificar tokens

func InitPaseto() {
	keyHex := os.Getenv("SECRET_KEY")

	if keyHex == "" {
		panic("SECRET_KEY is not set")
	}

	var err error
	SecretKey, err = paseto.V4SymmetricKeyFromHex(keyHex)

	if err != nil {
		panic(fmt.Sprintf("Error al inicializar la clave paseto: %v", err))
	}

}
