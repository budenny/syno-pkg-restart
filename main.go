package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	SYNOPKG = "/usr/syno/bin/synopkg"
)

type Cfg struct {
	Token    string
	HostRoot string
}

var cfg Cfg

func initCfg() error {
	cfg.Token = os.Getenv("TOKEN")
	if cfg.Token != "" {
		os.Setenv("TOKEN", "")
	}

	cfg.HostRoot = os.Getenv("HOST_ROOT")
	if cfg.HostRoot == "" {
		return errors.New("HOST_ROOT is not set")
	}

	log.Println("initCfg - done")
	return nil
}

func verifyAccessToSynoPkg() error {
	if err := syscall.Chroot(cfg.HostRoot); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	_, err := os.Stat(SYNOPKG)
	if err != nil {
		return err
	}

	log.Println("verifyAccessToSynoPkg - done")
	return nil
}

func restartSynoSvc(svcName string) error {
	cmd := exec.Command(SYNOPKG, "restart", svcName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func verifyAuthorization(r *http.Request) bool {
	if cfg.Token != "" {
		auth := r.Header.Get("Authorization")
		authSplitted := strings.Split(auth, "Bearer ")
		if len(authSplitted) != 2 || authSplitted[1] != cfg.Token {
			log.Printf("unauthorized: %s\n", auth)
			return false
		}
	}
	return true
}

func handler(resp http.ResponseWriter, req *http.Request) {
	if !verifyAuthorization(req) {
		http.Error(resp, "unauthorized", http.StatusUnauthorized)
		return
	}

	if req.Method != "GET" {
		log.Printf("unsupported method: %s\n", req.Method)
		http.Error(resp, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	svc := req.URL.Query().Get("svc")
	if svc == "" {
		log.Println("svc param is not set")
		http.Error(resp, "svc param is not set", http.StatusBadRequest)
		return
	}

	if err := restartSynoSvc(svc); err != nil {
		log.Printf("restartSynoSvc(%s) failed: %s\n", svc, err)
		http.Error(resp, "failed to restart service", http.StatusInternalServerError)
		return
	}
}

func main() {
	if err := initCfg(); err != nil {
		log.Panic(err)
	}

	if err := verifyAccessToSynoPkg(); err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/", handler)

	log.Println("start listening on :80 ...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
