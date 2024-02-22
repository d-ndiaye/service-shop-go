package main

import (
	"bitbucket.easy.de/users/n.gauche/service-shop-go/crypt"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/health"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/product"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/store"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/internal/storeProduct"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/api"
	"bitbucket.easy.de/users/n.gauche/service-shop-go/pkg/config"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const file = "config/service.yaml"

var configFile string
var ssl bool

func main() {
	flag.Parse()
	err, conf := config.Load(configFile)
	if err != nil {
		fmt.Println("Read config error")
		os.Exit(1)
	}
	repos, err := storeProduct.NewRepositoryStoreProduct(conf.Mongo)
	serviceProductStore := storeProduct.New(repos)
	storeProductHandler := api.StoreProductHandler{
		Sp: serviceProductStore,
	}

	repo, err := store.NewRepositoryStore(conf.Mongo)
	storeService := store.New(repo)
	storeHandler := api.StoreHandler{
		S: storeService,
	}

	repository, err := product.NewRepository(conf.Mongo)
	if err != nil {
		fmt.Println("could not initialize product repository")
		os.Exit(1)
	}
	health.InitHealthSystem(conf.Healthy, repository)
	productService := product.New(repository)
	productHandler := api.ProductHandler{
		S: productService,
	}
	storeRoutes := storeHandler.Routes()
	storeRoutes.Route("/", func(r chi.Router) {
		r.Mount("/", storeProductHandler.Routes())
	})
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Mount("/product", productHandler.Routes())
		r.Mount("/store", storeRoutes)
		r.Mount("/health", health.Routes())
	})

	gc := crypt.GenerateCertificate{
		Organization: "EASY SOFTWARE",
		Host:         "127.0.0.1",
		ValidFor:     10 * 365 * 24 * time.Hour,
		IsCA:         false,
		EcdsaCurve:   "P384",
		Ed25519Key:   false,
	}

	if conf.Sslport > 0 {
		ssl = true
		fmt.Println("ssl active")
	}
	var sslsrv *http.Server
	var srv *http.Server
	if ssl {
		tlsConfig, err := gc.GenerateTLSConfig()
		if err != nil {
			fmt.Printf("could not create tls config. %s\n", err.Error())
		}
		sslsrv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(conf.Sslport),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
			TLSConfig:    tlsConfig,
		}
		go func() {
			fmt.Printf("starting https server on address: %s\n", sslsrv.Addr)
			if err := sslsrv.ListenAndServeTLS("", ""); err != nil {
				fmt.Printf("error starting server: %s\n", err.Error())
			}
		}()
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(conf.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
		}
		go func() {
			fmt.Printf("starting http server on address: %s\n", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				fmt.Printf("error starting server: %s\n", err.Error())
			}
		}()
	} else {
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(conf.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
		}
		go func() {
			fmt.Printf("starting http server on address: %s\n", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				fmt.Printf("error starting server: %s\n", err.Error())
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("waiting for clients")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_ = srv.Shutdown(ctx)
	if ssl {
		_ = sslsrv.Shutdown(ctx)
	}

	fmt.Println("finished")

	os.Exit(0)
}

func init() {
	flag.StringVarP(&configFile, "config", "c", file, "this is the path and filename to the config file")
}
