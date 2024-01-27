package resolver

import (
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
)

type Resolver struct {
	Upstreams []string
	Cache     *cache.Cache
}

func New(upstreams []string, cache *cache.Cache) *Resolver {
	return &Resolver{
		Upstreams: upstreams,
		Cache:     cache,
	}
}

func (r *Resolver) Resolve(req dns.Msg) (dns.Msg, error) {
	cacheKey := r.generateCacheKey(req)
	if response, found := r.getFromCache(cacheKey, req); found {
		return response, nil
	}

	response, err := r.queryUpstream(req)
	if err != nil {
		return dns.Msg{}, err
	}

	r.Cache.Set(cacheKey, response, cache.DefaultExpiration)
	return response, nil
}

func (r *Resolver) generateCacheKey(req dns.Msg) string {
	return fmt.Sprintf("%s_%s", req.Question[0].Name, dns.TypeToString[req.Question[0].Qtype])
}

func (r *Resolver) getFromCache(key string, req dns.Msg) (dns.Msg, bool) {
	if cachedResp, found := r.Cache.Get(key); found {
		log.Printf("Cache hit for %s", key)
		cachedMsg := cachedResp.(dns.Msg)
		cachedMsg.Id = req.Id
		return cachedMsg, true
	}
	log.Printf("Cache miss for %s", key)
	return dns.Msg{}, false
}

func (r *Resolver) queryUpstream(req dns.Msg) (dns.Msg, error) {
	client := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(req.Question[0].Name, req.Question[0].Qtype)

	upstream := r.chooseUpstream()
	if upstream == "" {
		return dns.Msg{}, fmt.Errorf("no upstream DNS servers configured")
	}

	// Ensure the upstream address includes a port number
	if !strings.Contains(upstream, ":") {
		upstream = fmt.Sprintf("%s:53", upstream) // Default DNS port is 53
	}

	resp, _, err := client.Exchange(msg, upstream)
	if err != nil {
		return dns.Msg{}, err
	}

	response := r.createResponse(req, resp.Answer)
	return response, nil
}

func (r *Resolver) chooseUpstream() string {
	// Implement load balancing or failover logic if needed
	if len(r.Upstreams) > 0 {
		return r.Upstreams[0] // Simple selection for now
	}
	return ""
}

func (r *Resolver) createResponse(req dns.Msg, answers []dns.RR) dns.Msg {
	response := dns.Msg{}
	response.SetReply(&req)
	response.Authoritative = false
	response.Answer = answers
	return response
}
