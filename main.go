package main

import (
  "flag"
  "github.com/akesle/sailors/service"
  "time"
)

func main() {

  var (
    sailorDBPath             string
    sailorAPIVirtualPath     string
    serviceAddr              string
    sailorAPIBreakerRate     float64
    sailorAPIBreakerSamples  int64
    sailorAPIArtificialDelay time.Duration
  )

  flag.StringVar(&sailorDBPath, "sailor-db-conn", "api_user:example@tcp(127.0.0.1:3306)/sailors",
    "Sailor database connection string")
  flag.StringVar(&sailorAPIVirtualPath, "sailor-api-path", "/sailors", "Sailor API virtual path")
  flag.StringVar(&serviceAddr, "service-addr", ":8080", "Local address for the service")
  flag.Float64Var(&sailorAPIBreakerRate, "sailor-api-breaker-rate", 0.95,
    "Sailor API circuit breaker rate percentage")
  flag.Int64Var(&sailorAPIBreakerSamples, "sailor-api-breaker-samples", 5,
    "Sailor API circuit breaker count of samples")
  flag.DurationVar(&sailorAPIArtificialDelay, "sailor-api-artificial-delay", time.Millisecond*250,
    "Sailor API artificial delay for request processing")

  flag.Parse()

  s := service.SailorAPI{
    SailorDBPath:          sailorDBPath,
    SailorVirtualPath:     sailorAPIVirtualPath,
    BindAddress:           serviceAddr,
    SailorArtificialDelay: sailorAPIArtificialDelay,
    SailorBreakerRate:     sailorAPIBreakerRate,
    SailorBreakerSamples:  sailorAPIBreakerSamples,
  }
  _ = s.Run()
}
