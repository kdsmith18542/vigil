package api_server

import (
	"encoding/json"
	"net/http"

	"github.com/Vigil-Labs/vgl/internal/rpcserver"
)

// APIServer holds the RPC server instance to interact with the node.
type APIServer struct {
	rpcServer *rpcserver.Server
}

// NewAPIServer creates a new API server instance.
func NewAPIServer(rpcServer *rpcserver.Server) *APIServer {
	return &APIServer{
		rpcServer: rpcServer,
	}
}

// RegisterAPIRoutes registers the API routes for mining information and ticket pool value.
func (s *APIServer) RegisterAPIRoutes() {
	http.HandleFunc("/api/getmininginfo", s.getMiningInfoHandler)
	http.HandleFunc("/api/getticketpoolvalue", s.getTicketPoolValueHandler)
	http.HandleFunc("/api/getstakinginfo", s.getStakingInfoHandler)
}

// getMiningInfoHandler handles requests for mining information.
func (s *APIServer) getMiningInfoHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.rpcServer.GetMiningInfo(nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// getStakingInfoHandler handles requests for staking information.
func (s *APIServer) getStakingInfoHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.rpcServer.GetStakingInfo(nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// getTicketPoolValueHandler handles requests for the ticket pool value.
func (s *APIServer) getTicketPoolValueHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.rpcServer.GetTicketPoolValue(nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}




