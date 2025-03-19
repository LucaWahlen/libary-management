package main

import (
	"libary-service/internal/direct-service/repository"
	"libary-service/internal/direct-service/router"
	"log"
)

// Dein "direct-service" ist ein "besonderer" Fall, wo die ganzen Objekte im Prinzip Singletons sind, die ueber globale Package-Variablen abgebildet sind,
// und wo die Dependencies nur durch die Abhaengigkeit zu den anderen Packages sichbar werden (wo die "singleton" globalen Variablen liegen)
//
// Ich haette jetzt eher ein Beispiel erwartet, wo die "app" Schicht sich z.B. die Datenbank selber instanziiert und verwendet. Ich glaube, das ist das klassichere Beispiel
// von "keine dependency injection".
//
// Siehe auch https://en.wikipedia.org/wiki/Dependency_injection
// In software engineering, dependency injection is a programming technique in which an object or function receives other objects or functions that it requires,
// as opposed to creating them internally.
//
// Das "creating them internally" fehlt jetzt in deinem Beispiel.
//
// Bei deinem Beispiel ist die loose Koppelung und die fehlende Unit-Testbarkeit natuerlich auch gegeben, aber man sieht hier nicht in den packages die Erzeugung
// ihrer Abhaengigkeiten, die typisch fuer "keine dependency injection" ist.

func main() {

	if err := repository.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.Disconnect()
	router.NewGinRouter()
	if err := router.Serve(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
