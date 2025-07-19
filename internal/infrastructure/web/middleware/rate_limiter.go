package middleware

import (
	"backend_reservation/pkg/handler"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	visitors    map[string]int
	limit       int
	resetTime   time.Duration
	done        chan bool
	lastReset   time.Time
	maxVisitors int
}

// Stop detiene el goroutine de reseteo del RateLimiter enviando una señal al canal 'done'.
// Esto permite finalizar limpiamente el proceso de reseteo periódico de los contadores de visitantes.
// Debe llamarse, por ejemplo, al apagar el servidor para evitar fugas de goroutines.
func (rl *RateLimiter) Stop() {
	rl.done <- true
}

// resetVisitorCount es un método privado del RateLimiter que se encarga de reiniciar periódicamente
// el contador de solicitudes por IP almacenado en el mapa rl.visitors. Este método corre en un goroutine
// independiente y utiliza un ticker para ejecutar la lógica de reseteo cada rl.resetTime.
//
// Funcionamiento:
//   - Cada vez que el ticker dispara (cada rl.resetTime), se adquiere un candado (mutex) para evitar condiciones de carrera.
//   - Si la cantidad de IPs almacenadas supera rl.maxVisitors, se limpia completamente el mapa de visitantes
//     para evitar un consumo excesivo de memoria.
//   - Si la cantidad de IPs es aceptable, simplemente se resetean los contadores de solicitudes a 0,
//     pero se mantienen las IPs conocidas en el mapa para eficiencia.
//   - Se actualiza rl.lastReset con la hora actual para llevar registro del último reseteo.
//   - Si se recibe una señal en rl.done, el método retorna y termina el goroutine.
//
// Este mecanismo permite controlar el uso de memoria y mantener el control de la tasa de solicitudes
// de manera eficiente y segura en concurrencia.
func (rl *RateLimiter) resetVisitorCount() {
	ticker := time.NewTicker(rl.resetTime) // Crea un ticker que dispara cada rl.resetTime
	defer ticker.Stop()                    // Asegura que el ticker se detenga al salir de la función

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock() // Bloquea el acceso concurrente al mapa de visitantes

			// Si la cantidad de IPs supera el máximo permitido, limpia completamente el mapa
			if len(rl.visitors) > rl.maxVisitors {
				rl.visitors = make(map[string]int)
			} else {
				// Si no, resetea los contadores de cada IP a 0, manteniendo las IPs conocidas
				for ip := range rl.visitors {
					rl.visitors[ip] = 0
				}
			}

			rl.lastReset = time.Now() // Actualiza la marca de tiempo del último reseteo
			rl.mu.Unlock()            // Libera el candado
		case <-rl.done:
			// Si se recibe una señal de parada, termina el goroutine limpiamente
			return
		}
	}
}

// NewRateLimiter crea e inicializa una nueva instancia de RateLimiter.
// Recibe como parámetros el límite de solicitudes permitidas (limit) y el intervalo de tiempo para reiniciar el contador (resetTime).
// Inicia en segundo plano un goroutine que reinicia periódicamente el contador de visitantes.
// Devuelve un puntero a la estructura RateLimiter creada.
// NewRateLimiter crea e inicializa una nueva instancia de RateLimiter.
// Parámetros:
//   - limit: número máximo de solicitudes permitidas por IP en el intervalo especificado.
//   - resetTime: duración del intervalo tras el cual se reinician los contadores de solicitudes.
//   - maxVisitors (opcional): número máximo de IPs distintas a mantener en memoria antes de limpiar el mapa.
//
// Si no se especifica maxVisitors, se usa un valor por defecto de 1000.
// La función inicia un goroutine en segundo plano que se encarga de reiniciar los contadores periódicamente.
// Retorna un puntero a la estructura RateLimiter inicializada.
func NewRateLimiter(limit int, resetTime time.Duration, maxVisitors ...int) *RateLimiter {
	// Valor por defecto para el máximo de visitantes únicos
	maxVis := 1000

	// Si se proporciona un valor para maxVisitors, se utiliza ese valor
	if len(maxVisitors) > 0 {
		maxVis = maxVisitors[0]
	}

	// Inicializa la estructura RateLimiter con los parámetros dados
	rl := &RateLimiter{
		visitors:    make(map[string]int), // Mapa para llevar el conteo de solicitudes por IP
		limit:       limit,                // Límite de solicitudes por IP
		resetTime:   resetTime,            // Intervalo de reseteo de contadores
		done:        make(chan bool),      // Canal para detener el goroutine de reseteo
		lastReset:   time.Now(),           // Marca de tiempo del último reseteo
		maxVisitors: maxVis,               // Máximo de IPs a mantener en memoria
	}

	// Inicia el goroutine encargado de reiniciar los contadores periódicamente
	go rl.resetVisitorCount()

	// Retorna el puntero a la nueva instancia de RateLimiter
	return rl
}

// Throttle es un middleware que limita la cantidad de solicitudes permitidas por dirección IP
// en un intervalo de tiempo determinado. Si una IP excede el límite de solicitudes (rl.limit),
// responde con un error HTTP 429 (Too Many Requests) y no permite que la solicitud avance.
// Si no se ha alcanzado el límite, incrementa el contador de solicitudes para esa IP y
// permite que la solicitud continúe al siguiente handler.
//
// Parámetros:
//   - next: http.Handler que representa el siguiente handler en la cadena de middlewares.
//
// Retorna:
//   - http.Handler: un handler que aplica la lógica de limitación de tasa antes de invocar al siguiente handler.
func (rl *RateLimiter) Throttle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener la IP del cliente usando la función auxiliar, que maneja headers de proxy y conexión directa
		ip := getClientIP(r)

		// Bloqueamos el acceso concurrente al mapa de visitantes para evitar condiciones de carrera
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Obtenemos el número actual de solicitudes realizadas por esta IP
		count := rl.visitors[ip]

		// Calculamos el próximo momento en que se reiniciará el contador de rate limit
		nextReset := rl.lastReset.Add(rl.resetTime)

		// Establecemos headers informativos para el cliente sobre el rate limit
		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))         // Límite máximo permitido
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", nextReset.Unix())) // Timestamp UNIX del próximo reset

		// Si la IP ha excedido el límite de solicitudes permitidas
		if count >= rl.limit {
			w.Header().Set("X-RateLimit-Remaining", "0") // No quedan solicitudes disponibles
			// Header adicional que indica en cuántos segundos podrá volver a intentar
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", time.Until(nextReset).Seconds()))
			handler.Error(w, r, http.StatusTooManyRequests, "Rate limit exceeded") // Respuesta 429
			return
		}

		// Incrementamos el contador de solicitudes para esta IP
		rl.visitors[ip] = count + 1
		// Indicamos cuántas solicitudes le quedan antes de alcanzar el límite
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.limit-count-1))

		// Continuamos con el siguiente handler de la cadena
		next.ServeHTTP(w, r)
	})
}

// getClientIP extrae la IP del cliente de forma más confiable
func getClientIP(r *http.Request) string {
	// Primero intentamos obtener la IP del cliente desde el header "X-Forwarded-For" (usado por proxies y balanceadores)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// El header puede contener varias IPs separadas por comas, tomamos la primera (la del cliente original)
		if ips := parseXForwardedFor(xff); len(ips) > 0 {
			return ips[0]
		}
	}

	// Si no existe "X-Forwarded-For", intentamos con "X-Real-IP" (otro header común de proxies)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		// Validamos que sea una IP válida
		if ip := net.ParseIP(xri); ip != nil {
			return ip.String()
		}
	}

	// Si no hay headers de proxy, usamos la IP de la conexión directa (RemoteAddr)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Si ocurre un error al separar el host y el puerto, devolvemos el valor original
		return r.RemoteAddr
	}
	// Devolvemos solo la IP (sin el puerto)
	return ip
}

// parseXForwardedFor maneja correctamente el header X-Forwarded-For que puede tener múltiples IPs
// parseXForwardedFor toma el valor del header "X-Forwarded-For" (que puede contener una lista de IPs separadas por comas)
// y devuelve un slice con todas las IPs válidas encontradas.
// Este header es comúnmente utilizado por proxies y balanceadores de carga para indicar la cadena de IPs
// por las que ha pasado la petición, donde la primera IP suele ser la del cliente original.
//
// Parámetros:
//   - xff: string que representa el valor del header "X-Forwarded-For", por ejemplo: "203.0.113.1, 70.41.3.18, 150.172.238.178"
//
// Retorna:
//   - []string: un slice de strings con las IPs válidas extraídas y normalizadas.
//
// Ejemplo de uso:
//
//	ips := parseXForwardedFor("203.0.113.1, 70.41.3.18, 150.172.238.178")
//	// ips = ["203.0.113.1", "70.41.3.18", "150.172.238.178"]
func parseXForwardedFor(xff string) []string {
	var ips []string
	for ip := range strings.SplitSeq(xff, ",") {
		if trimmed := strings.TrimSpace(ip); trimmed != "" {
			if parsed := net.ParseIP(trimmed); parsed != nil {
				ips = append(ips, parsed.String())
			}
		}
	}
	return ips
}
