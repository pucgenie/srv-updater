package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	packet := Packet{}
	flag.StringVar(&packet.FQDN, "fqdn", "_service._tcp.example.com", "the hostname of the record you want to update")
	flag.UintVar(&packet.Priority, "priority", 10, "the priority of the record you want to update")
	flag.UintVar(&packet.Weight, "weight", 1, "the weight of the record you want to update")

	interval := flag.Duration("interval", time.Second*28, "sending interval")
	keyPath := flag.String("key", "~/.ssh/id_ed25519", "private key")
	dstPort := flag.Int("dst-port", 55, "the destination port")
	dstHost := flag.String("dst-host", "dynv6.com", "the destination host")
	srcPort := flag.Int("src-port", 10000, "the local port")
	ip := flag.String("src-ip", defaultIP(), "the local ip address")

	flag.Parse()

	// No arguments given?
	if len(os.Args) < 2 {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Expand home directory
	if str := *keyPath; strings.HasPrefix(str, "~/") {
		home, _ := os.UserHomeDir()
		home += str[1:]
		keyPath = &home
	}

	privKey := loadPrivateKey(*keyPath)

	// Marshal JSON data
	packet.Timestamp = time.Now().Unix()
	packet.Key = base64.RawStdEncoding.EncodeToString(privKey.PublicKey().Marshal())
	payload := packet.MarshalAndSign(privKey)

	// Parse local IP address
	srcIP := net.ParseIP(*ip)

	dstIPs, err := lookupIP(*dstHost)
	if err != nil {
		log.Panicln("unable to resolve destination host:", err)
	}

	// open raw socket
	conn, err := openSocket()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	srcAddr := net.UDPAddr{
		IP:   srcIP,
		Port: *srcPort,
	}

	// send packets
	for i := 0; ; i++ {
		dstAddr := net.UDPAddr{
			IP:   dstIPs[i%len(dstIPs)],
			Port: *dstPort,
		}

		log.Printf("sending from %s to %s", srcAddr.String(), dstAddr.String())
		b, err := buildUDPPacket(dstAddr, srcAddr, payload)
		if err != nil {
			panic(err)
		}

		_, err = conn.WriteTo(b, &net.IPAddr{IP: dstAddr.IP})
		if err != nil {
			log.Println(err)
		}

		time.Sleep(*interval)
	}
}

func loadPrivateKey(path string) ssh.Signer {
	// Load private key
	pemBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Parse private key
	privKey, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		panic(err)
	}

	return privKey
}
