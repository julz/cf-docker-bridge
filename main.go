package main

import (
	"encoding/json"
	"flag"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

var socket = flag.String("socket", "socker.dock", "socket to listen on")

func main() {
	flag.Parse()

	listener, err := net.Listen("unix", *socket)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1.20/containers/create", handleCreate)
	r.HandleFunc("/v1.20/containers/json", handleList)
	r.HandleFunc("/v1.20/containers/{id}/start", handleStart)
	r.HandleFunc("/v1.20/containers/{id}/attach", handleAttach)
	r.HandleFunc("/v1.20/containers/{id}/wait", handleWait)
	r.HandleFunc("/v1.20/containers/{id}", handleDelete).Methods("DELETE")

	logrus.Info("Listening on " + *socket)
	logrus.Infof("export DOCKER_HOST=unix://%s # <-- run this line in your shell to use with docker", *socket)
	panic(http.Serve(listener, r))
}

func handleList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]types.Container{
		{
			ID:      "not-implemented",
			Names:   []string{"/mycfapp"},
			Image:   "lattice-app",
			Command: "",
		},
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Image string
	}

	json.NewDecoder(r.Body).Decode(&req)

	name := r.FormValue("name")
	if name == "" {
		name = uuid.NewV4().String()
	}

	run(log(exec.Command("cf", "push", "--no-start", "-o", req.Image, name)))
	json.NewEncoder(w).Encode(&types.ContainerCreateResponse{ID: name})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["id"]
	run(exec.Command("cf", "delete", name))
}

func handleAttach(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["id"]
	stdout := stdcopy.NewStdWriter(w, stdcopy.Stdout)
	cmd := exec.Command("cf", "logs", name)
	cmdOut, _ := cmd.StdoutPipe()
	go io.Copy(stdout, cmdOut)
	run(cmd)
}

func handleWait(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(types.ContainerWaitResponse{})
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["id"]
	run(log(exec.Command("cf", "start", name)))
}

func log(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = logrus.StandardLogger().Writer()
	return cmd
}

func run(cmd *exec.Cmd) {
	logrus.Info(" > " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		logrus.Error(" ! ", err)
	}
}
