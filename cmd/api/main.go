package api

import (
	"github.com/caarlos0/env/v6"
	"github.com/cruffinoni/neobrain-todolist/internal/config"
	"github.com/cruffinoni/neobrain-todolist/internal/database"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	var configuration config.Global
	if err := env.Parse(&configuration, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatalf("error during parsing config: %v", err)
	}
	db, err := database.NewDB(configuration.Database)
	if err != nil {
		log.Fatalf("can't initialize connection to the database: %v", err)
	}

	//var serverStopped = false
	router := gin.New()
	//onSignalFn := func(sig os.Signal) {
	//	if serverStopped {
	//		return
	//	}
	//	serverStopped = true
	//	err := srv.Shutdown(context.Background())
	//	if err != nil && !errors.Is(err, http.ErrServerClosed) {
	//		log.Printf("error while shutting down the server: %v", err)
	//	}
	//}
	//ctx, cancelFn := contextual.NewContext(
	//	contextual.WithSignalListener(onSignalFn, syscall.SIGINT, syscall.SIGTERM),
	//)
}
