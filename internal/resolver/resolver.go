package resolver

import (
	"log"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
)

type Resolver struct {
	Upstreams []string
	Cache     *cache.Cache
	// Add a map to associate DNS record types with their respective resolve functions
	resolveFuncs map[uint16]func(dns.Msg) (dns.Msg, error)
}

func New() *Resolver {
	r := &Resolver{
		Upstreams: []string{"8.8.8.8:53"}, // Assuming Google's DNS as the upstream
		Cache:     cache.New(5*time.Minute, 10*time.Minute),
	}
	r.initResolveFuncs()
	return r
}

// initResolveFuncs initializes the map of DNS record types to resolve functions.
func (r *Resolver) initResolveFuncs() {
	r.resolveFuncs = map[uint16]func(dns.Msg) (dns.Msg, error){
		dns.TypeA:     r.resolveA,
		dns.TypeAAAA:  r.resolveAAAA,
		dns.TypeCNAME: r.resolveCNAME,
		dns.TypeMX:    r.resolveMX,
		dns.TypeNS:    r.resolveNS,
		dns.TypeTXT:   r.resolveTXT,
	}
}

func (r *Resolver) Resolve(req dns.Msg) (dns.Msg, error) {
	cacheKey := req.Question[0].Name + "_" + dns.TypeToString[req.Question[0].Qtype]

	if cachedResp, found := r.Cache.Get(cacheKey); found {
		log.Printf("Cache hit for %s", cacheKey)
		cachedMsg := cachedResp.(dns.Msg)

		// Adjust the ID of the cached response to match the incoming request's ID
		cachedMsg.Id = req.Id

		return cachedMsg, nil
	}
	log.Printf("Cache miss for %s", cacheKey)

	resolveFunc, ok := r.resolveFuncs[req.Question[0].Qtype]
	if !ok {
		return dns.Msg{}, nil // Return empty message for unsupported types
	}

	response, err := resolveFunc(req)
	if err != nil {
		return dns.Msg{}, err
	}

	// Cache the response
	r.Cache.Set(cacheKey, response, cache.DefaultExpiration)

	return response, nil
}

func (r *Resolver) resolveA(req dns.Msg) (dns.Msg, error) {

	c := new(dns.Client)

	// Set up a message to query the external DNS server
	m := new(dns.Msg)
	m.SetQuestion(req.Question[0].Name, dns.TypeA)

	// Perform the query
	resp, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		// Handle error
		return dns.Msg{}, err
	}

	// Construct the response to the original query
	var response dns.Msg
	response.SetReply(&req)
	response.Authoritative = false // Set false since it's not authoritative
	response.Answer = resp.Answer

	return response, nil
}

func (r *Resolver) resolveAAAA(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveCNAME(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveMX(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveNS(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveTXT(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) createResponse(req dns.Msg, record dns.RR) dns.Msg {
	resp := dns.Msg{}
	resp.SetReply(&req)
	resp.Answer = append(resp.Answer, record)
	return resp
}
