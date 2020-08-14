package main

import (
	crand "crypto/rand"
	"encoding/json"

	"golang.org/x/crypto/ssh"
)

type Packet struct {
	FQDN      string `json:"fqdn"`
	Weight    uint   `json:"weight"`
	Priority  uint   `json:"priority"`
	Key       string `json:"key"`
	Timestamp int64  `json:"ts"`
}

func (p *Packet) MarshalAndSign(signer ssh.Signer) []byte {
	jsonBytes, _ := json.Marshal(p)

	// Sign JSON data
	signature, err := signer.Sign(crand.Reader, jsonBytes)
	if err != nil {
		panic(err)
	}

	// Build payload
	payload := make([]byte, len(jsonBytes)+1+len(signature.Blob))
	copy(payload, jsonBytes)
	copy(payload[len(jsonBytes)+1:], signature.Blob)

	return payload
}
