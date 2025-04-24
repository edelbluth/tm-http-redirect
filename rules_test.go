package tm_http_redirect_test

import (
	"fmt"
	"testing"

	"git.codebau.dev/goblins/commons/pkg/ref"
	to "github.com/edelbluth/tm_http_redirect"
)

type RuleExpectation struct {
	InputUrl    string
	ExpectedUrl *string
}

type RuleTestCase struct {
	Rule  to.Redirect
	Tests []RuleExpectation
	Name  string
}

type RuleParsingTestCase struct {
	Rules         []to.Redirect
	ExpectedError bool
}

func TestRules(t *testing.T) {
	t.Parallel()
	log := to.NamedLogger("test")
	log.Collecting = false
	tests := []RuleTestCase{
		{
			Name: "Simple rule",
			Rule: to.Redirect{
				From: "/a",
				To:   "/b",
				Code: nil,
			},
			Tests: []RuleExpectation{
				{
					InputUrl:    "/a",
					ExpectedUrl: ref.AsRef("/b"),
				},
				{
					InputUrl:    "/b",
					ExpectedUrl: nil,
				},
				{
					InputUrl:    "/",
					ExpectedUrl: nil,
				},
				{
					InputUrl:    "",
					ExpectedUrl: nil,
				},
			},
		},
		{
			Name: "Replacing rule (single element)",
			Rule: to.Redirect{
				From: "/a/(.*)",
				To:   "/b/${1}/c",
				Code: nil,
			},
			Tests: []RuleExpectation{
				{
					InputUrl:    "/a/",
					ExpectedUrl: ref.AsRef("/b//c"),
				},
				{
					InputUrl:    "/a",
					ExpectedUrl: nil,
				},
				{
					InputUrl:    "/a/go/to/here",
					ExpectedUrl: ref.AsRef("/b/go/to/here/c"),
				},
			},
		},
		{
			Name: "Replacing rule (multiple elements)",
			Rule: to.Redirect{
				From: "/a/(.*)/b/(.*)",
				To:   "/b/${1}/c/${2}/d",
				Code: nil,
			},
			Tests: []RuleExpectation{
				{
					InputUrl:    "/a/",
					ExpectedUrl: nil,
				},
				{
					InputUrl:    "/a/1/b/2",
					ExpectedUrl: ref.AsRef("/b/1/c/2/d"),
				},
				{
					InputUrl:    "/a/1/b/2/d/e/f",
					ExpectedUrl: ref.AsRef("/b/1/c/2/d/e/f/d"),
				},
			},
		},
		{
			Name: "Replacing rule (default)",
			Rule: to.Redirect{
				From: ".*",
				To:   "/dashboard/",
				Code: nil,
			},
			Tests: []RuleExpectation{
				{
					InputUrl:    "/",
					ExpectedUrl: ref.AsRef("/dashboard/"),
				},
				{
					InputUrl:    "",
					ExpectedUrl: ref.AsRef("/dashboard/"),
				},
				{
					InputUrl:    "/not-dashboard/",
					ExpectedUrl: ref.AsRef("/dashboard/"),
				},
				{
					InputUrl:    "/dashboard",
					ExpectedUrl: ref.AsRef("/dashboard/"),
				},
			},
		},
	}
	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%s-%d", test.Name, i), func(subTest *testing.T) {
			subTest.Parallel()
			rules, err := to.ParseRules(&[]to.Redirect{test.Rule}, log)
			if err != nil {
				subTest.Fatalf("error parsing rules: %s", err.Error())
			}
			if len(rules) != 1 {
				subTest.Fatalf("parsed %d rules, expected only 1", len(rules))
			}
			for j, expectation := range test.Tests {
				expectation := expectation
				subTest.Run(fmt.Sprintf("%s-%s-%d", test.Name, expectation.InputUrl, j), func(subSubTest *testing.T) {
					subSubTest.Parallel()
					rule := rules[0]
					result := rule.Handle(expectation.InputUrl)
					if result == nil && expectation.ExpectedUrl != nil {
						subSubTest.Fatalf("expected %s, got nil", *expectation.ExpectedUrl)
					}
					if result != nil && expectation.ExpectedUrl == nil {
						subSubTest.Fatalf("expected nil, got %s", *result)
					}
					if result != nil && expectation.ExpectedUrl != nil && *result != *expectation.ExpectedUrl {
						subSubTest.Fatalf("expected %s, got %s", *expectation.ExpectedUrl, *result)
					}
				})
			}
		})
	}
}

func TestParsingRules(t *testing.T) {
	t.Parallel()
	log := to.DefaultLogger()
	testCases := []RuleParsingTestCase{
		{
			ExpectedError: false,
			Rules: []to.Redirect{
				{
					From: "/a",
					To:   "/b",
				},
				{
					From: "/c",
					To:   "/e",
				},
				{
					From: "/a/.*",
					To:   "/b",
				},
				{
					From: "/a/(.*)",
					To:   "/b/${1}",
				},
				{
					From: `/a/(\d+)/u`,
					To:   "/b/${1}",
				},
			},
		},
		{
			ExpectedError: true,
			Rules: []to.Redirect{
				{
					From: "",
					To:   "/a",
				},
			},
		},
		{
			ExpectedError: true,
			Rules: []to.Redirect{
				{
					From: "/a",
					To:   "",
				},
			},
		},
		{
			ExpectedError: true,
			Rules: []to.Redirect{
				{
					From: "+[",
					To:   "/b",
				},
			},
		},
	}
	for i, test := range testCases {
		test := test
		t.Run(fmt.Sprintf("Test-%d", i), func(subTest *testing.T) {
			subTest.Parallel()
			_, err := to.ParseRules(&test.Rules, log)
			if !test.ExpectedError && err != nil {
				subTest.Fatalf("expected no error, but got %q", err.Error())
			}
			if test.ExpectedError && err == nil {
				subTest.Fatal("expected error, but did not get one")
			}
		})
	}
}
