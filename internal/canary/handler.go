package canary

import (
	"encoding/json"
	"log"
	"net/http"

	"tailscale.com/client/tailscale"
)

type Handler struct {
	pinger *Pinger
}

func NewHandler(lc *tailscale.LocalClient) *Handler {
	return &Handler{
		pinger: NewPinger(lc),
	}
}

func (h *Handler) GetPeers(w http.ResponseWriter, r *http.Request) {
	peers, err := h.pinger.GetPeers(r.Context())
	if err != nil {
		log.Printf("Failed to get peers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	var req PingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IP == "" {
		http.Error(w, "IP address required", http.StatusBadRequest)
		return
	}

	result, err := h.pinger.Ping(r.Context(), req.IP)
	if err != nil {
		log.Printf("Failed to ping %s: %v", req.IP, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) PingAll(w http.ResponseWriter, r *http.Request) {
	results, err := h.pinger.PingAll(r.Context())
	if err != nil {
		log.Printf("Failed to ping all: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
