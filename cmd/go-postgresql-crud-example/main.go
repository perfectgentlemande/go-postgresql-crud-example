package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/api"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/database/dbuser"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/logger"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/service"
	dbUserMigrations "github.com/perfectgentlemande/go-postgresql-crud-example/migrations/dbuser"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()
	log := logger.DefaultLogger()

	configPath := flag.String("c", "config.yaml", "path to your config")
	flag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		log.WithField("config_path", *configPath).WithError(err).Fatal("failed to read config")
	}

	dbUser, err := dbuser.NewDatabase(ctx, conf.DBUser, dbUserMigrations.MigrationAssets)
	if err != nil {
		log.WithError(err).Fatal("cannot create db")
	}
	defer dbUser.Close(ctx)

	err = dbUser.Ping(ctx)
	if err != nil {
		log.WithField("conn_string", *configPath).WithError(err).Error("cannot ping db")
		return
	}

	serverParams := api.ServerParams{
		Cfg:  conf.Server,
		Srvc: service.NewService(dbUser),
		Log:  log,
	}
	srv := api.NewServer(&serverParams)

	rungroup, ctx := errgroup.WithContext(ctx)

	log.WithField("address", srv.Addr).Info("starting server")
	rungroup.Go(func() error {
		if er := srv.ListenAndServe(); er != nil && !errors.Is(er, http.ErrServerClosed) {
			return fmt.Errorf("listen and server %w", er)
		}

		return nil
	})

	rungroup.Go(func() error {
		<-ctx.Done()

		if er := srv.Shutdown(context.Background()); er != nil {
			return fmt.Errorf("shutdown http server %w", er)
		}

		return nil
	})

	err = rungroup.Wait()
	if err != nil {
		log.WithError(err).Error("run group exited because of error")
		return
	}

	log.Info("server exited properly")
}
