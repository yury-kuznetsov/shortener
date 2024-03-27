package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/storage/database"
	"github.com/yury-kuznetsov/shortener/internal/storage/file"
	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
)

func main() {
	printBuildData()

	config.Init()

	storage, err := buildStorage()
	if err != nil {
		panic(err)
	}
	coder := uricoder.NewCoder(storage)

	// запустим два сервера: http и grpc
	var wg sync.WaitGroup
	wg.Add(2)
	httpSrv, err := startHTTPServer(coder, &wg)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
	grpcSrv, lis, err := startGrpcServer(coder, &wg)
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	// готовим канал для прослушивания системных сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// ожидаем сигнала остановки из канала `stop`
	<-stop

	// даем серверу 5 секунд на завершение обработки текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// завершаем "мягко" работу серверов
	if err := httpSrv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown: %v", err)
	}
	grpcSrv.GracefulStop()
	if err := lis.Close(); err != nil {
		fmt.Printf("Server Shutdown: %v", err)
	}

	wg.Wait()
	log.Println("Both servers are stopped successfully")
}

func buildStorage() (uricoder.Storage, error) {
	if len(config.Options.Database) > 0 {
		return database.NewStorage(config.Options.Database)
	}
	if len(config.Options.FilePath) > 0 {
		return file.NewStorage(config.Options.FilePath)
	}
	return memory.NewStorage(), nil
}
