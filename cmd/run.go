package cmd

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/Projectopher/config"
	"github.com/saleh-ghazimoradi/Projectopher/infra/mongodb"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
	"github.com/saleh-ghazimoradi/Projectopher/internal/server"
	"github.com/saleh-ghazimoradi/Projectopher/internal/service"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if len(groups) == 0 && a.Key == slog.TimeKey {
					t := a.Value.Time()
					a.Value = slog.StringValue(t.Format("2006-01-02T15:04:05"))
				}
				return a
			},
		}))

		cfg, err := config.GetInstance()
		if err != nil {
			logger.Error("failed to get config", "error", err.Error())
			os.Exit(1)
		}

		mongo := mongodb.NewMongoDB(
			mongodb.WithHost(cfg.MongoDB.Host),
			mongodb.WithPort(cfg.MongoDB.Port),
			mongodb.WithUser(cfg.MongoDB.User),
			mongodb.WithPass(cfg.MongoDB.Pass),
			mongodb.WithDBName(cfg.MongoDB.DBName),
			mongodb.WithAuthSource(cfg.MongoDB.AuthSource),
			mongodb.WithMaxPoolSize(cfg.MongoDB.MaxPoolSize),
			mongodb.WithMinPoolSize(cfg.MongoDB.MinPoolSize),
			mongodb.WithTimeout(cfg.MongoDB.Timeout),
		)

		client, mongodb, err := mongo.Connect()
		if err != nil {
			logger.Error("failed to connect", "error", err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := client.Disconnect(context.Background()); err != nil {
				logger.Error("failed to disconnect", "error", err.Error())
				os.Exit(1)
			}
		}()

		middleware := middlewares.NewMiddleware(cfg, logger)
		validator := helper.NewValidator()

		movieRepository := repository.NewMovieRepository(mongodb, "movie")
		genreRepository := repository.NewGenresRepository(mongodb, "genre")
		rankRepository := repository.NewRankingsRepository(mongodb, "rank")
		userRepository := repository.NewUsersRepository(mongodb, "user")
		tokenRepository := repository.NewTokenRepository(mongodb, "token")

		movieService := service.NewMovieService(movieRepository)
		genreService := service.NewGenreService(genreRepository)
		rankService := service.NewRankingsService(rankRepository)
		authService := service.NewAuthService(cfg, userRepository, tokenRepository)
		userService := service.NewUserService(userRepository)

		healthHandler := handlers.NewHealthHandler(cfg)
		movieHandler := handlers.NewMovieHandler(movieService, validator)
		genreHandler := handlers.NewGenreHandler(genreService)
		rankHandler := handlers.NewRankingHandler(rankService)
		authHandler := handlers.NewAuthHandler(validator, authService)
		userHandler := handlers.NewUserHandler(validator, userService)

		healthRoute := routes.NewHealthRoute(healthHandler)
		movieRoute := routes.NewMovieRoute(middleware, movieHandler)
		genreRoute := routes.NewGenreRoutes(genreHandler)
		rankRoute := routes.NewRankRoutes(rankHandler)
		authRoute := routes.NewAuthRoute(authHandler)
		userRoute := routes.NewUserRoute(userHandler)

		register := routes.NewRegister(
			routes.WithHealthRoute(healthRoute),
			routes.WithAuthRoute(authRoute),
			routes.WithMovieRoute(movieRoute),
			routes.WithGenreRoute(genreRoute),
			routes.WithRankRoute(rankRoute),
			routes.WithUserRoute(userRoute),
			routes.WithMiddleware(middleware),
		)

		httpServer := server.NewServer(
			server.WithHost(cfg.Server.Host),
			server.WithPort(cfg.Server.Port),
			server.WithHandler(register.RegisterRoutes()),
			server.WithReadTimeout(cfg.Server.ReadTimeout),
			server.WithWriteTimeout(cfg.Server.WriteTimeout),
			server.WithIdleTimeout(cfg.Server.IdleTimeout),
			server.WithErrorLog(slog.NewLogLogger(logger.Handler(), slog.LevelError)),
			server.WithLogger(logger),
		)

		logger.Info("starting server", "addr", cfg.Server.Host+":"+cfg.Server.Port, "env", cfg.Application.Environment)
		if err := httpServer.Connect(); err != nil {
			logger.Error("failed to connect", "error", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
