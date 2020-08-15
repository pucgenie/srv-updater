package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func openSocket() (net.PacketConn, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, fmt.Errorf("failed open socket: %w", err)
	}
	syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)

	return net.FilePacketConn(os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd)))
}

func defaultIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func lookupIP(host string) ([]net.IP, error) {
	return net.DefaultResolver.LookupIP(context.Background(), "ip4", host)
}

func buildUDPPacket(dst, src net.UDPAddr, payload []byte) ([]byte, error) {
	buffer := gopacket.NewSerializeBuffer()

	ip := &layers.IPv4{
		DstIP:    dst.IP,
		SrcIP:    src.IP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(src.Port),
		DstPort: layers.UDPPort(dst.Port),
	}
	if err := udp.SetNetworkLayerForChecksum(ip); err != nil {
		return nil, fmt.Errorf("failed calculate checksum: %w", err)
	}
	if err := gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true}, ip, udp, gopacket.Payload(payload)); err != nil {
		return nil, fmt.Errorf("failed to serialize packet: %w", err)
	}

	return buffer.Bytes(), nil
}
