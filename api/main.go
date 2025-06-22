package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/tjroeder/censys_take_home/cache/grpcserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// kvRequestBody defines the request body accepted by handleSet
type kvRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// kvHandler handles requests for key-value cache endpoint requests
type kvHandler struct {
	grpcClient grpcserver.CacheClient
}

// handleSet handles POST /v1/keyvalues requests
func (h *kvHandler) handleSet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var body kvRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// TODO: key validation, probably length cap, and character values

	slog.InfoContext(ctx, fmt.Sprintf("Received: POST %s, Body: %+v", r.URL.Path, body))
	_, err := h.grpcClient.Set(ctx, &grpcserver.SetRequest{
		Key:   body.Key,
		Value: []byte(body.Value),
	})
	if err != nil {
		// TODO: log 500 error to monitoring
		slog.ErrorContext(ctx, fmt.Sprintf("Response: POST %s, Status Code: %d", r.URL.Path, http.StatusInternalServerError))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("Response: POST %s, Status Code: %d", r.URL.Path, http.StatusCreated))
	w.WriteHeader(http.StatusCreated)
}

// handleGet handles GET /v1/keyvalues/{key} requests
func (h *kvHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	key := r.PathValue("key")
	slog.InfoContext(ctx, fmt.Sprintf("Received: GET %s", r.URL.Path))

	resp, err := h.grpcClient.Get(ctx, &grpcserver.GetRequest{Key: key})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			slog.WarnContext(ctx, fmt.Sprintf("Response: GET %s, Status Code: %d", r.URL.Path, http.StatusNotFound))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		// TODO: log 500 error to monitoring
		slog.ErrorContext(ctx, fmt.Sprintf("Response: GET %s, Status Code: %d", r.URL.Path, http.StatusInternalServerError))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.InfoContext(ctx, fmt.Sprintf("Response: GET %s, Status Code: %d", r.URL.Path, http.StatusOK))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"value": string(resp.Value),
	})
}

// handleDelete handles DELETE /v1/keyvalues/{key} requests
func (h *kvHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	key := r.PathValue("key")

	slog.InfoContext(ctx, fmt.Sprintf("Received: DELETE %s", r.URL.Path))
	_, err := h.grpcClient.Delete(ctx, &grpcserver.DeleteRequest{Key: key})
	if err != nil {
		// TODO: log 500 error to monitoring
		slog.ErrorContext(ctx, fmt.Sprintf("Response: DELETE %s, Status Code: %d", r.URL.Path, http.StatusInternalServerError))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("Response: DELETE %s, Status Code: %d", r.URL.Path, http.StatusNoContent))
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	target := os.Getenv("GRPC_TARGET")
	if target == "" {
		target = "localhost:50051"
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcserver.NewCacheClient(conn)
	h := &kvHandler{grpcClient: client}

	r := http.NewServeMux()
	r.HandleFunc("GET /v1/keyvalues/{key}", h.handleGet)
	r.HandleFunc("POST /v1/keyvalues", h.handleSet)
	r.HandleFunc("DELETE /v1/keyvalues/{key}", h.handleDelete)

	log.Println("api-gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
