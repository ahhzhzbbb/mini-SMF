package handler

import (
	"encoding/json"
	"fmt"
	"mini-SMF/pdu-session/internal/config"
	"net/http"
	"os"
	"time"
)

func HandlerGetMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!!!"))
	})
}

func HandlerGetServerConfigInfo(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Host: %s\nPort: %s\n", config.Host, config.Port)
	})
}

func HandlerPDUSessionEstablishment(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		type rqBody struct {
			Supi         string `json:"supi,omitempty"`         //Subscriber Permanent Identifier
			Gpsi         string `json:"gpsi,omitempty"`         //Generic Public Subscription Identifier
			PduSessionId int    `json:"pduSessionId,omitempty"` //pdu session id
			Dnn          string `json:"dnn,omitempty"`          //Data Network Name
			SNssai       SNssai `json:"sNssai"`                 //single Network Slice
			ServingNfId  string `json:"servingNfId,omitempty"`  //Serving Network Function ID
			AnType       string `json:"anType,omitempty"`       //Access NetWork Type
		}

		instanceID, err := os.Hostname()
		if err != nil {
			instanceID = "unknown-node"
		}

		type rspBody struct {
			SmContextRef string `json:"smContextRef,omitempty"`
			Supi         string `json:"supi,omitempty"`
			PduSessionId int    `json:"pduSessionId,omitempty"`
			HandledBy    string `json:"handledBy,omitempty"`
			Status       string `json:"status,omitempty"`
		}

		var req rqBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		smCtx := NewSmContext(req.Supi, req.PduSessionId, req.Dnn, "ACTIVE", "10.0.0.1", req.SNssai)
		smCtx.SaveInDatabase()

		var rsp rspBody
		rsp.SmContextRef = "http://gw/nsmf-pdusession/v1/sm-contexts/" + smCtx.ContextID
		rsp.Supi = req.Supi
		rsp.PduSessionId = req.PduSessionId
		rsp.HandledBy = instanceID
		rsp.Status = "ACTIVE"

		time.Sleep(2 * time.Second)

		if err := json.NewEncoder(w).Encode(&rsp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Fprintf(w, "Hello Client, I am instance: %s, (Listening on %s:%s)",
		// 	instanceID, config.Host, config.Port)
	})
}

func HandlerGetHeath(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instanceID, err := os.Hostname()
		if err != nil {
			instanceID = "unknown-node"
		}
		fmt.Fprintf(w, "Hello Client, I am instance: %s, (Listening on %s:%s)",
			instanceID, config.Host, config.Port)
	})
}
