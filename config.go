package tm_http_redirect

type Redirect struct {
	From string `json:"from,omitempty" yaml:"from,omitempty"`
	To   string `json:"to,omitempty" yaml:"to,omitempty"`
	Code *int   `json:"code,omitempty" yaml:"code,omitempty"`
}

type Config struct {
	Redirects *[]Redirect `json:"redirects,omitempty" yaml:"redirects,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Redirects: &[]Redirect{},
	}
}
