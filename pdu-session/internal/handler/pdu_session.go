package handler

import (
	"fmt"

	"github.com/google/uuid"
)

type SNssai struct {
	Sst int    `json:"sst,omitempty"` //Service type (eMBB, URLLC, mMTC)
	Sd  string `json:"sd,omitempty"`  //Slice Differntiator
}

type SmContext struct {
	ContextID    string
	Supi         string
	PduSessionID int
	Dnn          string
	SNssai       SNssai
	Status       string
	IP           string
}

func NewSmContext(
	supi string,
	pduSessionID int,
	dnn string,
	status string,
	ip string,
	sNssai SNssai,
) *SmContext {
	return &SmContext{
		ContextID:    uuid.NewString(),
		Supi:         supi,
		PduSessionID: pduSessionID,
		Dnn:          dnn,
		SNssai:       sNssai,
		Status:       status,
		IP:           ip,
	}
}

func (ctx *SmContext) SaveInDatabase() {
	fmt.Println("Saving...")
}
