package primaryzoneprovider

import (
	"strings"
	"sync"

	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/config"
	"gitlab.com/nameserver-systems/pdns-distribute/internal/app/pdns-secondary-syncer/powerdns"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice"
	"gitlab.com/nameserver-systems/pdns-distribute/pkg/microservice/logger"
)

var (
	localdnsaxfrqueries       = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_local_dns_proxy_queries_total", ConstLabels: map[string]string{"qtype": "axfr", "state": "permitted"}, Help: "The total count of local dns queries"})
	localdnssoaqueries        = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_local_dns_proxy_queries_total", ConstLabels: map[string]string{"qtype": "soa", "state": "permitted"}, Help: "The total count of local dns queries"})
	localdnsprohibitedqueries = promauto.NewCounter(prometheus.CounterOpts{Name: "secondarysyncer_local_dns_proxy_queries_total", ConstLabels: map[string]string{"qtype": "*", "state": "prohibited"}, Help: "The total count of local dns queries"})
)

var tcpserver, udpserver *dns.Server

func StartAXFRProvider(ms *microservice.Microservice, conf *config.ServiceConfiguration) {
	startDNSServer(conf.AXFRPrimaryAddress)
	setDNSHandler(conf, ms)
}

func StopDNSServer() (err error) {
	if tcpserver != nil {
		if err = tcpserver.Shutdown(); err != nil {
			return
		}
	}

	if udpserver != nil {
		if err = udpserver.Shutdown(); err != nil {
			return
		}
	}

	return
}

func startDNSServer(address string) {
	tcpserver = &dns.Server{Net: "tcp", Addr: address}
	udpserver = &dns.Server{Net: "udp", Addr: address}

	go func() {
		if err := tcpserver.ListenAndServe(); err != nil {
			logger.ErrorErrLog(err)
		}
	}()

	go func() {
		if err := udpserver.ListenAndServe(); err != nil {
			logger.ErrorErrLog(err)
		}
	}()
}

func setDNSHandler(conf *config.ServiceConfiguration, ms *microservice.Microservice) {
	dns.HandleFunc(".", func(writer dns.ResponseWriter, msg *dns.Msg) {
		ansmsg := new(dns.Msg)

		if qtypeNotAllowed(msg) {
			refuseAnswer(writer, msg, ansmsg)

			localdnsprohibitedqueries.Inc()

			return
		}

		if msg.Question[0].Qtype == dns.TypeSOA {
			handleSOA(writer, msg, ansmsg, conf, ms)

			localdnssoaqueries.Inc()

			return
		}

		handleAXFR(writer, msg, conf, ms)
		localdnsaxfrqueries.Inc()
	})
}

func handleAXFR(writer dns.ResponseWriter, msg *dns.Msg, conf *config.ServiceConfiguration, ms *microservice.Microservice) {
	zoneid := getZoneID(msg)

	zonefile, err := retrieveZoneFile(zoneid, conf, ms, false)
	if err != nil {
		logger.ErrorErrLog(err)

		return
	}

	zp := dns.NewZoneParser(strings.NewReader(zonefile), "", "")

	var ansrrset []dns.RR

	for rr, ok := zp.Next(); ok; rr, ok = zp.Next() {
		ansrrset = append(ansrrset, rr)
	}

	if err := zp.Err(); err != nil {
		logger.ErrorErrLog(err)

		return
	}

	ch := make(chan *dns.Envelope)
	tr := new(dns.Transfer)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		if err = tr.Out(writer, msg, ch); err != nil {
			logger.ErrorErrLog(err)
		}

		wg.Done()
	}()
	ch <- &dns.Envelope{RR: ansrrset}
	close(ch)
	wg.Wait()
	writer.Close()
}

func getZoneID(msg *dns.Msg) string {
	zoneid := msg.Question[0].Name
	if !strings.HasSuffix(zoneid, ".") {
		zoneid += "."
	}

	return zoneid
}

func retrieveZoneFile(zoneid string, conf *config.ServiceConfiguration, ms *microservice.Microservice, onlysoa bool) (string, error) {
	payload, marshalerr := powerdns.PrepareZoneDataRequest(zoneid, onlysoa)
	if marshalerr != nil {
		return "", marshalerr
	}

	zonedata, getdataerr := powerdns.GetZoneData(conf, ms, payload)
	if getdataerr != nil {
		return "", getdataerr
	}

	return zonedata.ZoneData, nil
}

func handleSOA(writer dns.ResponseWriter, msg *dns.Msg, ansmsg *dns.Msg, conf *config.ServiceConfiguration, ms *microservice.Microservice) {
	ansmsg.SetReply(msg)

	zoneid := getZoneID(msg)

	zonefile, err := retrieveZoneFile(zoneid, conf, ms, true)
	if err != nil {
		logger.ErrorErrLog(err)

		return
	}

	soarr, err := dns.NewRR(zonefile)
	if err != nil {
		logger.ErrorErrLog(err)

		return
	}

	ansmsg.Insert([]dns.RR{soarr})

	if err = writer.WriteMsg(ansmsg); err != nil {
		logger.ErrorErrLog(err)
	}

	if err = writer.Close(); err != nil {
		logger.ErrorErrLog(err)
	}
}

func refuseAnswer(writer dns.ResponseWriter, msg *dns.Msg, ansmsg *dns.Msg) {
	ansmsg.SetReply(msg)
	ansmsg.Rcode = dns.RcodeRefused

	if err := writer.WriteMsg(ansmsg); err != nil {
		logger.ErrorErrLog(err)
	}

	if err := writer.Close(); err != nil {
		logger.ErrorErrLog(err)
	}
}

func qtypeNotAllowed(msg *dns.Msg) bool {
	return msg.Question[0].Qtype != dns.TypeAXFR && msg.Question[0].Qtype != dns.TypeSOA
}
