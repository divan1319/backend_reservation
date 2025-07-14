package firmador

import (
	"time"

	"aidanwoods.dev/go-paseto"
)

// Firmar un token con los datos proporcionados y una duracion
func FirmarToken(data map[string]string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	for key, value := range data {
		token.Set(key, value)
	}

	now := time.Now()

	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(duration))

	return token.V4Encrypt(SecretKey, nil), nil
}

// Verificar un token y retornar los datos del token
func VerificarToken(tokenStr string) (*paseto.Token, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())        // Verifica que el token no haya expirado
	parser.AddRule(paseto.ValidAt(time.Now())) // Verifica que el token sea valido en el momento actual

	return parser.ParseV4Local(SecretKey, tokenStr, nil) // Parseamos el token y lo retornamos

}
