package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rajsinghtech/tailtunnel/internal/ssh"
	"github.com/rajsinghtech/tailtunnel/internal/tailscale"
)

type Handler struct {
	ts         *tailscale.TailscaleClient
	sshHandler *ssh.SSHHandler
}

func NewHandler(ts *tailscale.TailscaleClient) *Handler {
	return &Handler{
		ts: ts,
		sshHandler: &ssh.SSHHandler{
			DialFunc: ts.DialSSH,
		},
	}
}

func (h *Handler) GetMachines(w http.ResponseWriter, r *http.Request) {
	machines, err := h.ts.GetSSHMachines(r.Context())
	if err != nil {
		log.Printf("Failed to get machines: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(machines)
}

func (h *Handler) SSHWebSocket(w http.ResponseWriter, r *http.Request) {
	machine := chi.URLParam(r, "machine")
	if machine == "" {
		http.Error(w, "machine parameter required", http.StatusBadRequest)
		return
	}

	user := r.URL.Query().Get("user")
	if user == "" {
		user = "root"
	}

	h.sshHandler.HandleWebSocket(w, r, machine, user)
}
