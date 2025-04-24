package tm_http_redirect

import "errors"

var ErrMisconfiguration = errors.New("middelware misconfiguration - see log for details")
var ErrNoConfigurationFound = errors.New("middleware misconfiguration - no configuration found - see log for details")
var ErrRuleParsingFailed = errors.New("middleware misconfiguration - parsing of rules failed - see log for details")
var ErrLogCollectionInactive = errors.New("log collection was inactive, no logs collected")
