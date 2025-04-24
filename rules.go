package tm_http_redirect

import (
	"regexp"
)

type Rule struct {
	FromMatcher *regexp.Regexp
	To          string
	Code        int
}

func ParseRules(configs *[]Redirect, log *Logger) ([]*Rule, error) {
	rules := []*Rule{}
	thereAreErrorsInRules := false
	for index, redirect := range *configs {
		ruleHasError := false
		if redirect.To == "" {
			ruleHasError = true
			log.Error(`rule %d: no "to" defined`, index)
		}
		if redirect.From == "" {
			ruleHasError = true
			log.Error(`rule %d: no "from" defined`, index)
		}
		if ruleHasError {
			thereAreErrorsInRules = true
			continue
		}
		ruleRegExp, err := regexp.Compile(redirect.From)
		if err != nil {
			log.Error(`rule %d: unable to compile "from" rule: %s`, index, err.Error())
			thereAreErrorsInRules = true
			continue
		}
		var status int
		if redirect.Code != nil {
			status = *redirect.Code
		} else {
			status = DefaultRedirectionStatusCode
		}
		rules = append(rules, &Rule{
			FromMatcher: ruleRegExp,
			To:          redirect.To,
			Code:        status,
		})
	}
	if thereAreErrorsInRules {
		return nil, ErrRuleParsingFailed
	}
	numberOfRules := len(rules)
	if numberOfRules == 0 {
		log.Error("no rules parsed (out of %d configured rules)", len(*configs))
		return nil, ErrRuleParsingFailed
	}
	if numberOfRules != 1 {
		log.Info("parsed %d rules", numberOfRules)
	} else {
		log.Info("parsed %d rule", numberOfRules)
	}
	return rules, nil
}

func (r *Rule) Handle(url string) *string {
	if r.FromMatcher.MatchString(url) {
		redirect := r.FromMatcher.ReplaceAllString(url, r.To)
		return &redirect
	}
	return nil
}
