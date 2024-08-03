package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func StartHttp(mux *http.ServeMux) {
	logrus.Println("Launching http server...")
	host := viper.GetString("server.address")
	if host == "" {
		host = ":9111"
	}

	if err := http.ListenAndServe(host, mux); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
