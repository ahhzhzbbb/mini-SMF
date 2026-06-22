package server

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
)

func getLocalAddress() (string, error) {
	var address string
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return address, err
		}
		defer conn.Close()
		address = conn.LocalAddr().(*net.UDPAddr).IP.String()
	} else {
		address = "localhost"
	}
	return address, nil
}

func Register(gatewayService string) error {
	ips, err := net.LookupIP(gatewayService)
	if err != nil {
		return err
	}

	type reqBody struct {
		ServiceName string `json:"service_name"`
		Ip          string `json:"ip"`
		Port        string `json:"port"`
	}

	ip := ips[rand.IntN(len(ips))]
	gatewayAddress := net.JoinHostPort(ip.String(), os.Getenv("GW_PORT"))
	reqURL := "http://" + gatewayAddress + "/register"

	localIP, err := getLocalAddress()
	if err != nil {
		return err
	}

	bodyData := reqBody{
		ServiceName: "pdu-session",
		Ip:          localIP,
		Port:        "8081",
	}
	jsonBytes, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
