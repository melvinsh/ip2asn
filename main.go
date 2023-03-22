package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	concurrency := 10
	semaphore := make(chan struct{}, concurrency)

	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())

		if net.ParseIP(ip) == nil {
			fmt.Fprintf(os.Stderr, "%s is not a valid IP address\n", ip)
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(ip string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			asn, err := lookupASN(ip)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error looking up ASN for %s: %s\n", ip, err)
				return
			}

			if asn == "NA" {
				return
			}

			fmt.Println(asn)
		}(ip)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
		os.Exit(1)
	}

	wg.Wait()
}

func lookupASN(ip string) (string, error) {
	conn, err := net.Dial("tcp", "whois.cymru.com:43")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "-r %s\n", ip)

	var buf strings.Builder
	if _, err := io.Copy(&buf, conn); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(buf.String()))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "AS") {
			continue
		}
		return extractNumber(line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("ASN not found for %s", ip)
}

func extractNumber(s string) (string, error) {
	fields := strings.Fields(s)

	if len(fields) == 0 {
		return "", fmt.Errorf("invalid string: %s", s)
	}

	return fields[0], nil
}
