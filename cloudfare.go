package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cloudflare/cloudflare-go"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/crowdsecurity/crowdsec/pkg/sqlite"
	"github.com/crowdsecurity/crowdsec/pkg/types"
)

type context struct {
	api       *cloudflare.API
	apiKey    string
	emailAddr string
	zoneID    string
	accountID string
	scope     string
}

func newCloudflareContext(bConfig *blockerConfig) (*context, error) {
	var err error

	ctx := &context{
		apiKey:    bConfig.APIKey,
		emailAddr: bConfig.Email,
		zoneID:    bConfig.ZoneID,
		accountID: bConfig.AccountID,
		scope:     bConfig.Scope,
	}

	ctx.api, err = cloudflare.New(ctx.apiKey, ctx.emailAddr)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func newRuleConfiguration(value string) cloudflare.AccessRuleConfiguration {
	ruleConfig := cloudflare.AccessRuleConfiguration{}
	if value != "" {
		ip := net.ParseIP(value)
		_, cidr, cidrErr := net.ParseCIDR(value)
		_, asnErr := strconv.ParseInt(value, 10, 32)
		if ip != nil {
			ruleConfig.Target = "ip"
			ruleConfig.Value = ip.String()
		} else if cidrErr == nil {
			cidr.IP = cidr.IP.Mask(cidr.Mask)
			ruleConfig.Target = "ip_range"
			ruleConfig.Value = cidr.String()
		} else if asnErr == nil {
			ruleConfig.Target = "asn"
			ruleConfig.Value = value
		} else {
			ruleConfig.Target = "country"
			ruleConfig.Value = value
		}
	}
	return ruleConfig
}

func (c *context) newAccessRule(ba types.BanApplication) error {
	var err error

	ruleConfig := newRuleConfiguration(ba.IpText)
	rule := cloudflare.AccessRule{
		Configuration: ruleConfig,
	}

	response, err := c.listAccessRule(rule)
	if err != nil {
		return err
	}
	rule.Mode = "block"
	rule.Notes = ba.Reason

	if len(response.Result) > 0 { // if rule already exist, return
		log.Debugf("ip : '%s' already banned", ba.IpText)
		return nil

	}
	log.Debugf("creating access rule for : %s", ba.IpText)
	switch c.scope {
	case "account":
		_, err = c.api.CreateAccountAccessRule(c.accountID, rule)
	case "zone":
		_, err = c.api.CreateZoneAccessRule(c.zoneID, rule)
	default:
		_, err = c.api.CreateUserAccessRule(rule)
	}
	if err != nil {
		log.Errorf("Error creating firewall access rule: %s", err)
	}

	return nil
}

func (c *context) listAllAccessRule() (*cloudflare.AccessRuleListResponse, error) {
	var err error

	rule := cloudflare.AccessRule{
		Mode: "all",
	}

	var response *cloudflare.AccessRuleListResponse
	switch c.scope {
	case "account":
		response, err = c.api.ListAccountAccessRules(c.accountID, rule, 1)
	case "zone":
		response, err = c.api.ListZoneAccessRules(c.zoneID, rule, 1)
	default:
		response, err = c.api.ListUserAccessRules(rule, 1)
	}
	if err != nil {
		return nil, err
	}
	return response, nil

}

func (c *context) listAccessRule(rule cloudflare.AccessRule) (*cloudflare.AccessRuleListResponse, error) {
	var response *cloudflare.AccessRuleListResponse
	var err error

	switch c.scope {
	case "account":
		response, err = c.api.ListAccountAccessRules(c.accountID, rule, 1)
	case "zone":
		response, err = c.api.ListZoneAccessRules(c.zoneID, rule, 1)
	default:
		response, err = c.api.ListUserAccessRules(rule, 1)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *context) deleteAllRules() error {
	var cpt int
	var wg sync.WaitGroup

	rules, err := c.listAllAccessRule()
	if err != nil {
		return err
	}

	for _, r := range rules.Result {
		wg.Add(1)
		go c.deleteRuleWorker(&wg, r.ID)
		cpt++
	}
	wg.Wait()
	log.Printf("%d bans deleted", cpt)
	return nil
}

func (c *context) deleteRuleWorker(wg *sync.WaitGroup, ruleID string) {
	defer wg.Done()
	err := c.deleteAccessRule(ruleID)
	if err != nil {
		log.Errorf(err.Error())
	}
}

func (c *context) deleteAccessRule(ruleID string) error {
	switch c.scope {
	case "account":
		_, err := c.api.DeleteAccountAccessRule(c.accountID, ruleID)
		if err != nil {
			return fmt.Errorf("error deleting account rule")
		}
	case "zone":
		_, err := c.api.DeleteZoneAccessRule(c.zoneID, ruleID)
		if err != nil {
			return fmt.Errorf("error deleting zone rule")
		}
	default:
		_, err := c.api.DeleteUserAccessRule(ruleID)
		if err != nil {
			return fmt.Errorf("error deleting user rule")
		}
	}
	return nil
}

func (c *context) deleteRule(ba types.BanApplication) error {
	var err error

	ruleConfig := newRuleConfiguration(ba.IpText)
	rule := cloudflare.AccessRule{
		Configuration: ruleConfig,
	}

	response, err := c.listAccessRule(rule)
	if err != nil {
		return err
	}

	if len(response.Result) > 0 {
		for _, r := range response.Result {
			err := c.deleteAccessRule(r.ID)
			if err != nil {
				log.Errorf("error while removing rule for '%s' : %s", ba.IpText, err)
			}
			log.Debugf("access rule for '%s' deleted", ba.IpText)
		}
	}
	return nil
}

func (c *context) Run(dbCTX *sqlite.Context, frequency time.Duration) {
	lastDelTS := time.Now()
	lastAddTS := time.Now()
	/*start by getting valid bans in db ^^ */
	log.Infof("fetching existing bans from DB")
	bansToAdd, err := getNewBan(dbCTX)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("found %d bans in DB", len(bansToAdd))
	for idx, ba := range bansToAdd {
		log.Debugf("ban %d/%d", (idx + 1), len(bansToAdd))
		go c.newAccessRule(ba)

	}
	/*go for loop*/
	for {
		time.Sleep(frequency)
		bas, err := getDeletedBan(dbCTX, lastDelTS)
		if err != nil {
			log.Fatal(err)
		}

		lastDelTS = time.Now()

		if len(bas) > 0 {
			log.Infof("%d bans to flush since %s", len(bas), lastDelTS)
		}

		for idx, ba := range bas {
			log.Debugf("delete ban %d/%d", (idx + 1), len(bas))
			go c.deleteRule(ba)
		}

		bansToAdd, err := getLastBan(dbCTX, lastAddTS)
		if err != nil {
			log.Fatal(err)
		}
		lastAddTS = time.Now()
		if len(bansToAdd) > 0 {
			log.Printf("Adding %d new bans", len(bansToAdd))
		}

		for idx, ba := range bansToAdd {
			log.Debugf("ban %d/%d", (idx + 1), len(bansToAdd))
			go c.newAccessRule(ba)
		}
	}
}
