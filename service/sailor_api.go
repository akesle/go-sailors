package service

import (
  "context"
  "database/sql"
  "github.com/akesle/sailors/controllers"
  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
  "github.com/golang-sql/sqlexp"
  circuit "github.com/rubyist/circuitbreaker"
  "log"
  "time"
)

type SailorAPI struct {
  SailorDBPath          string
  SailorVirtualPath     string
  BindAddress           string
  SailorArtificialDelay time.Duration
  SailorBreakerRate     float64
  SailorBreakerSamples  int64
}

type CircuitBreakerQuerier struct {
  CB             *circuit.Breaker
  Timeout        time.Duration
  Querier        sqlexp.Querier
  SimulatedDelay time.Duration
}

func (q *CircuitBreakerQuerier) ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
  var (
    cancellableCtx context.Context
    cancel         context.CancelFunc
  )
  if err = q.CB.Call(func() error {
    cancellableCtx, cancel = context.WithCancel(ctx)
    time.Sleep(q.SimulatedDelay)
    result, err = q.Querier.ExecContext(cancellableCtx, query, args...)
    return err
  }, q.Timeout); err != nil {
    // Open breaker will prevent setting of CancelFunc
    if cancel != nil {
      cancel()
    }
  }
  return
}

func (q *CircuitBreakerQuerier) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
  var (
    cancellableCtx context.Context
    cancel         context.CancelFunc
  )
  if err = q.CB.Call(func() error {
    cancellableCtx, cancel = context.WithCancel(ctx)
    time.Sleep(q.SimulatedDelay)
    rows, err = q.Querier.QueryContext(ctx, query, args...)
    return err
  }, q.Timeout); err != nil {
    // Open breaker will prevent setting of CancelFunc
    if cancel != nil {
      cancel()
    }
  }
  return
}

func (q *CircuitBreakerQuerier) QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
  var (
    cancellableCtx context.Context
    cancel         context.CancelFunc
  )
  if err := q.CB.Call(func() error {
    cancellableCtx, cancel = context.WithCancel(ctx)
    time.Sleep(q.SimulatedDelay)
    row = q.Querier.QueryRowContext(ctx, query, args...)
    return nil
  }, q.Timeout); err != nil {
    // Open breaker will prevent setting of CancelFunc
    if cancel != nil {
      cancel()
    }
  }
  return
}

func (s *SailorAPI) Run() error {
  router := gin.Default()

  db, dbErr := sql.Open(sqlexp.DialectMySQL, s.SailorDBPath)
  if dbErr != nil {
    log.Printf("Failed to open connection to Sailors database: %v", dbErr)
    return dbErr
  }
  defer func() {
    if err := db.Close(); err != nil {
      log.Printf("Failed closing database: %v", err)
    }
  }()
  cb := circuit.NewRateBreaker(s.SailorBreakerRate, s.SailorBreakerSamples)
  cbQuery := &CircuitBreakerQuerier{
    CB:             cb,
    Timeout:        time.Second * 2,
    Querier:        db,
    SimulatedDelay: s.SailorArtificialDelay,
  }

  // Wire-up the Sailor controller functionality
  sc := &controllers.SailorController{
    DBSrc: func() sqlexp.Querier {
      return cbQuery
    },
  }

  // Wire-up the routes
  router.POST(s.SailorVirtualPath, sc.AddSailor)
  router.GET(s.SailorVirtualPath, sc.FindSailor)
  router.DELETE(s.SailorVirtualPath, sc.RemoveSailor)
  router.PUT(s.SailorVirtualPath, sc.ModifySailor)

  return router.Run(s.BindAddress)
}
