package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

type Record struct {
	Target string
	Values []string
	Type string
}

type DNS struct {
	NameServers []string `yaml:"nameservers"`
	Records []Record `yaml:"records"`
}

func loadConfig() (*DNS, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var dns DNS
	err = yaml.Unmarshal(data, &dns)
	if err != nil {
		return nil, err
	}

	return &dns, nil
}

func Boot(s *DNS) (*dns.Server) {
	log.Printf("Serving %s and forwarding rest to %s\n", s.Records, s.NameServers)

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		for _, q := range req.Question {
			log.Printf("DNS query for %#v", q.Name)

			for _, record := range s.Records {
				for _, value := range record.Values {
					if strings.HasSuffix(q.Name, value) {
						m := new(dns.Msg)
						m.SetReply(req)
						m.Authoritative = true
						defer w.WriteMsg(m)

						log.Printf("Resolve DNS query for %#v to %s (%s)", q.Name, record.Target, record.Type)
						
						switch record.Type {
							case "A":
								m.Answer = append(m.Answer, &dns.A{
									A:   net.ParseIP(record.Target),
									Hdr: dns.RR_Header{Name: q.Name, Class: q.Qclass, Ttl: 0, Rrtype: dns.TypeA},
								})
							case "CNAME":
								m.Answer = append(m.Answer, &dns.CNAME{
									Target: record.Target,
									Hdr: dns.RR_Header{Name: q.Name, Class: q.Qclass, Ttl: 0, Rrtype: dns.TypeCNAME},
								})
							case "SRV":
								m.Answer = append(m.Answer, &dns.SRV{
									Target: record.Target,
									Hdr: dns.RR_Header{Name: q.Name, Class: q.Qclass, Ttl: 0, Rrtype: dns.TypeSRV},
								})
							default:
								log.Printf("Unknown record type %s", record.Type)
							}
						return
					}
				}
			}
		}

		log.Println("Forwarding DNS query")

		client := &dns.Client{Net: "udp", SingleInflight: true}

		for _, ns := range s.NameServers {
			if r, _, err := client.Exchange(req, ns+":53"); err == nil {
				if r.Rcode == dns.RcodeSuccess {
					r.Compress = true
					w.WriteMsg(r)
					for _, a := range r.Answer {
						log.Printf("Answer from %s: %v\n", ns, a)
					}
					return
				}
			}
		}

		log.Println("Failure to forward request")
		m := new(dns.Msg)
		m.SetReply(req)
		m.SetRcode(req, dns.RcodeServerFailure)
		w.WriteMsg(m)
	})

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		for {
			select {
			case s := <-sig:
				log.Fatalf("fatal: signal %s received\n", s)
			}
		}
	}()

	server := &dns.Server{Addr: ":53", Net: "udp", TsigSecret: nil}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to setup server: %v\n", err)
	}
	return server
}

func main() {
	dns, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	Boot(dns)
}