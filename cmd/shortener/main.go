package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	handlers "github.com/yury-kuznetsov/shortener/internal/app"
	"github.com/yury-kuznetsov/shortener/internal/auth"
	"github.com/yury-kuznetsov/shortener/internal/gzip"
	"github.com/yury-kuznetsov/shortener/internal/logger"
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

	r := buildRouter(coder)

	// создаем сервер
	server := &http.Server{Addr: config.Options.HostAddr, Handler: r}

	// готовим канал для прослушивания системных сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// запускаем сервера в отдельной горутине
	go func() {
		var err error
		if config.Options.Secure {
			certFile, keyFile := getCertAndKey()
			err = http.ListenAndServeTLS(config.Options.HostAddr, certFile, keyFile, r)
		} else {
			err = http.ListenAndServe(config.Options.HostAddr, r)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server Listen: %v", err)
		}
	}()

	// ожидаем сигнала остановки из канала `stop`
	<-stop

	// даем серверу 5 секунд на завершение обработки текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// завершаем "мягко" работу сервера
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown: %v", err)
	}
}

func getCertAndKey() (certFile, keyFile string) {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Shortener"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM
	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		log.Fatal(err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Fatal(err)
	}

	return certPEM.String(), privateKeyPEM.String()
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

func buildRouter(coder *uricoder.Coder) *chi.Mux {
	sugar := logger.NewLogger()

	r := chi.NewRouter()
	r.Get("/{code}", auth.Handle(gzip.Handle(sugar.Handle(handlers.DecodeHandler(coder))), true))
	r.Get("/ping", auth.Handle(gzip.Handle(sugar.Handle(handlers.PingHandler(coder))), true))
	r.Get("/api/user/urls", auth.Handle(gzip.Handle(sugar.Handle(handlers.UserUrlsHandler(coder))), false))
	r.Delete("/api/user/urls", auth.Handle(gzip.Handle(sugar.Handle(handlers.DeleteUrlsHandler(coder))), true))
	r.Post("/api/shorten/batch", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeBatchHandler(coder))), true))
	r.Post("/api/shorten", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeJSONHandler(coder))), true))
	r.Post("/", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeHandler(coder))), true))
	r.MethodNotAllowed(auth.Handle(gzip.Handle(sugar.Handle(handlers.NotAllowedHandler())), true))

	// обработчики для pprof
	r.Handle("/debug/pprof/*", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	// обработчик для favicon.ico (иначе перехватит DecodeHandler)
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return r
}
