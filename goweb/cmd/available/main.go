package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func exists(domain string) (bool, error) {
	const whoisServer string = "com.whois-servers.net"
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return false, err
	}
	defer conn.Close()
	conn.Write([]byte(domain + "rn"))
	s := bufio.NewScanner(conn)
	for s.Scan() {
		if strings.Contains(strings.ToLower(s.Text()), "no match") {
			return false, nil
		}
	}
	return true, nil
}

var marks = map[bool]string{true: "x", false: "o"}

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		domain := s.Text()
		fmt.Print(domain, " ")
		exist, err := exists(domain)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(marks[!exist])
		time.Sleep(1 * time.Second)
	}
}
