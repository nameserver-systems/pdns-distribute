//nolint:prealloc
package powerdnsutils

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func GetZoneFilePerAXFR(address, zoneid string, readtimeout time.Duration) (string, error) {
	// does not work for zones contain LUA records
	// used for dnssec signed zones
	return doZoneTransfer(zoneid, address, readtimeout)
}

func doZoneTransfer(zoneid string, address string, readtimeout time.Duration) (zonefile string, err error) {
	var conn net.Conn

	var channel chan *dns.Envelope

	var d net.Dialer

	axfrquery := prepareAxfrQuery(zoneid)

	if conn, err = d.Dial("tcp", address); err != nil {
		return
	}

	dnscon := &dns.Conn{Conn: conn}
	transfer := &dns.Transfer{Conn: dnscon, ReadTimeout: readtimeout}

	if channel, err = transfer.In(axfrquery, address); err != nil {
		return
	}

	answerrr, err := receiveAllRecords(channel)

	zonefile = buildZoneFileFromRecords(answerrr, zonefile)

	return
}

func receiveAllRecords(channel chan *dns.Envelope) ([]dns.RR, error) {
	var answerrr []dns.RR

	for envmsg := range channel {
		if envmsg.Error != nil {
			return nil, envmsg.Error
		}

		answerrr = append(answerrr, envmsg.RR...)
	}

	return answerrr, nil
}

func buildZoneFileFromRecords(answerrr []dns.RR, zonefile string) string {
	for _, rr := range answerrr {
		// if rr.Header().Rrtype != dns.TypeNSEC && rr.Header().Rrtype != dns.TypeNSEC3 && rr.Header().Rrtype != dns.TypeNSEC3PARAM {
		zonefile += rr.String() + string('\n')
	}

	zonefile = strings.TrimSuffix(zonefile, string('\n'))

	return zonefile
}

func prepareAxfrQuery(zoneid string) *dns.Msg {
	axfrquery := new(dns.Msg)
	axfrquery.SetAxfr(dns.Fqdn(zoneid))

	return axfrquery
}

func GetSOARecord(zoneid, address, net string) (zonesoa string, err error) {
	var answer *dns.Msg

	soarequest := new(dns.Msg)
	soarequest.SetQuestion(dns.Fqdn(zoneid), dns.TypeSOA)

	c := new(dns.Client)
	c.Net = net

	if answer, _, err = c.Exchange(soarequest, address); err != nil {
		return
	}

	if answer.Rcode != dns.RcodeSuccess {
		err = errors.New("soa request was not successful")

		return
	}

	zonesoa = answer.Answer[0].String()

	return
}
