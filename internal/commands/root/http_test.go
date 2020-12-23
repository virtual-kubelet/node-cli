package root

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/node-cli/provider/mock"
	"gotest.tools/assert"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

var (
	testCert = []byte(`
-----BEGIN CERTIFICATE-----
MIIEWDCCAsCgAwIBAgIRANyVR7WKUbcq9aleJdW9ZhgwDQYJKoZIhvcNAQELBQAw
gYExHjAcBgNVBAoTFW1rY2VydCBkZXZlbG9wbWVudCBDQTErMCkGA1UECwwiY3B1
Z3V5ODNAQnJpYW5zLU1TRlQtbWFjYm9vay5sb2NhbDEyMDAGA1UEAwwpbWtjZXJ0
IGNwdWd1eTgzQEJyaWFucy1NU0ZULW1hY2Jvb2subG9jYWwwHhcNMjAxMTMwMTgz
OTQ1WhcNMzAxMTMwMTgzOTQ1WjBkMScwJQYDVQQKEx5ta2NlcnQgZGV2ZWxvcG1l
bnQgY2VydGlmaWNhdGUxOTA3BgNVBAsMMGNwdWd1eTgzQEpQTUFZQU5BR0kwNC5m
YXJlYXN0LmNvcnAubWljcm9zb2Z0LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAN5BoyOVbMilICvc5PChVy/Kk90uDIQl3ryIli0Xom49hFfSjvmW
d2oEtwYRf3PQeQmSCFdR5zQcmTPHCbxWyipYXqeWZVcSnhB5dH5euTi8XzbpYpNP
NEE0fTmK8KbdatqGmHRlOHp8aSABrQLdXwy8bhs44NRlbEHieWvy5mkiy9htcCBi
Gz3qs9ozNRI4U8KXuw3wKChBJuInUW8Sv+/aDFP8doIJF/uSYLXWAI5904aRK4dB
30+9UAemrbUj6cT0p5js6VUJVz2sYRIrKOn3Pwoi8/GIo+tpBXDVnUf+kGWSj29S
se15t3JMEdc2ZBDu8OwA36dte3nkJ7l0o68CAwEAAaNnMGUwDgYDVR0PAQH/BAQD
AgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgw
FoAUStWHau9E6srUyZ7WfQp9zX6uuhIwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG
9w0BAQsFAAOCAYEAqDoiSxfYgTO0UcBvnfjjh+f5lDWkoh+AZNb6josY+koer0z8
zCJKNTSSvYVOm/kFIihiclOrW/nQ0XT6Vep9P3zs4RfU18DgUbohTzFhx1U0SrZP
J4Og5Bbw/Ue8j4AwZtvq8uOFEIlHjqAuZr87yT1HvbDurhEIuKAZveY9ScVR5ZCJ
0KWJr4071F/UYVaE+U9wf9fTq1gaCCHA5IxLSDf7bUyqXGWexwFzykI/GOwh2WDS
MpEOPQBuj6rZQlyAwanolrLjEPrK2sjDcWfYvBXCE6WALXoNKqwBNNtbXU5jXRDJ
4bFfSFUJGqmBgP2XnrbqQsf6luQSeF77MrivX13UWCAndxlCL5wx/t+1cg8Fh3u7
aMysJIscqHYz42vgdR1uQ/qHebgKowF2L/GbXw59Lhcj4vY/1iWxmBA6M9+VNAAt
cuiDQAfG4WGFRdGPaQ5q9w/RJTMAQcjHVitVpoC5ikXk93R4skqWkHUmk+dhCjDW
XNxkOUOi5eYwvkxk
-----END CERTIFICATE-----
`)

	testKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDeQaMjlWzIpSAr
3OTwoVcvypPdLgyEJd68iJYtF6JuPYRX0o75lndqBLcGEX9z0HkJkghXUec0HJkz
xwm8VsoqWF6nlmVXEp4QeXR+Xrk4vF826WKTTzRBNH05ivCm3Wrahph0ZTh6fGkg
Aa0C3V8MvG4bOODUZWxB4nlr8uZpIsvYbXAgYhs96rPaMzUSOFPCl7sN8CgoQSbi
J1FvEr/v2gxT/HaCCRf7kmC11gCOfdOGkSuHQd9PvVAHpq21I+nE9KeY7OlVCVc9
rGESKyjp9z8KIvPxiKPraQVw1Z1H/pBlko9vUrHtebdyTBHXNmQQ7vDsAN+nbXt5
5Ce5dKOvAgMBAAECggEBAK1xHWVca1sc+UEhjYt27Ln/5Wn6UIwjnXEVSdSAmCJd
YVTDnQ2K7T9P1KAosYRokLv2OQojgUC6fJfaYG+YbwWilqNDi2vqvGzwywb+1p4+
6jLI6EM60PV9h6eLFIezTHqiBID4qJ11TvhKNoCAznb66RXXSiSVzWiQ2t5x3Hr3
1z5OulVuyLxM6v0TG4I5PJaNbJSawQTjFkhinH6GWRDzP11NX1BJsBdy15EgICBx
Nb91SDsa+5Ojb64Vl5FG8BE2zDHeEAyKDojCE4GLMY9L3VbXouUK3/9Ha8k6pKB3
Tl2AfChTDAo9Qa+o2pJaTynv3wqFBb4rcWot2SLac1ECgYEA7YRVCrMMsAW7otv6
5+HMV8zXFshgl0E0m7kC5ouxu5ZBO5EYhkpdDNXDX5Mr/BGJ0t4lv37Datjs67Eq
+5B7HhbX+WCPpr6M8jQWhyAwxVKhslEPBiR2ycFSnFE1R+xPfK9q3B2Ezq/rBY5Q
74jC/W2EY9+4/A6c3flU6zWLDskCgYEA741LS26t+Li5Q4eDJZ0J/G56EoyiHRIE
LGXZsegVwqAVdnFn7qW9VLV6mGtC8QdSC9tlbpBmMcv6kTruYBKOgnfFr8iZPL+S
NWZ8WL0AJ6l3nUxBvwkc2/h8N/m/OoH4y39CfKtu2Zu5yxGqBcYN9d3XUTHDIUO+
iD2LKQezgrcCgYAYQL7+TLIq9yrlwliofOIExSHhbayPRVU94XJuYC1R3lHi5zn9
3HIL8Xf1tm1zW8cbBRwNpcAGlQf8OScOcP5hYCvFhxqkCCkUQkVanurb+0gPkT9b
fTWz/E2XMKOkKHklXjQnLcx13ni9JH8XNnvSrPAr0phtBID4GZGWQu1kIQKBgQCj
BrycTGmXYFeszneBTJt0MNdg8laNhCpU8MeznKfaeUnB/rHlpuPv10XknvLCx+Gd
ciVYlmsGLrSKy9lYhqh3v/1IgTNQNWvSbbnoRk/prhpacYA4+4GpbjVTfuMWdUeV
bjkYUS8yZxmNSqs0HLJ5hg04E66hX9I2M/QV60jOhwKBgQCEFqYOpc27mssizx4x
AjG7w9c0IAhDOnyBLl2SKhGBkt0JY37z0SeHOIpoSonSPfe3L2MfetFuNHxg6mA3
KLI7vlVFot19Ihh74SYqHaP69ej+kKEsJyWY/0357PZdoR8ZuaI43hnvOLT/DE26
zHNS9w8NRPV8bCIu1f+dFocMDw==
-----END PRIVATE KEY-----
`)

	testClientCert = []byte(`
-----BEGIN CERTIFICATE-----
MIIEVzCCAr+gAwIBAgIQDolyQDKLLb7NJfwhvzwnOjANBgkqhkiG9w0BAQsFADCB
gTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSswKQYDVQQLDCJjcHVn
dXk4M0BCcmlhbnMtTVNGVC1tYWNib29rLmxvY2FsMTIwMAYDVQQDDClta2NlcnQg
Y3B1Z3V5ODNAQnJpYW5zLU1TRlQtbWFjYm9vay5sb2NhbDAeFw0yMDExMzAxNzU4
MTRaFw0zMDExMzAxNzU4MTRaMGQxJzAlBgNVBAoTHm1rY2VydCBkZXZlbG9wbWVu
dCBjZXJ0aWZpY2F0ZTE5MDcGA1UECwwwY3B1Z3V5ODNASlBNQVlBTkFHSTA0LmZh
cmVhc3QuY29ycC5taWNyb3NvZnQuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEA1slkMwg6e7fSVC/ILs42+TWqFErLVwKvNrXYIyHFvdFPKdNqP+5S
BQVtiuPdnPoklOqZwC2a86dWi1TbiyElPG1hzK4OnrlKZG2IQRsCTJljke+CAZqB
qLdmBdywuVwxZOa7ZerEIBuQnl5A3pLiDNcfBdCXirbaSFr5ktgmarbLwW8WPwmv
f2J6cKWDV3DIx8VUg2u09OAeJ4nvAAsTAUAVMD58zYy3PlfGNeEqCbgYnOzGBgfX
VH/1US6GthvFlPTAbQ5pzr0Hq5LgUbXGhKgAokC9c8QxFwtoxHGCIFRMiYfn1TUp
BkKH5/owCZbT23pROUdWtPFK0aDy7UiMjQIDAQABo2cwZTAOBgNVHQ8BAf8EBAMC
BaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAW
gBRK1Ydq70TqytTJntZ9Cn3Nfq66EjAPBgNVHREECDAGhwR/AAABMA0GCSqGSIb3
DQEBCwUAA4IBgQA46TeXdFatO2SPMhnjcV0GjIpZ8gWpyHMEIcLSg9V8BtF2GbcL
arlMqZ2Na5aa7JZ/QAIYx4UaPbGGTcs4KYOrb0tGBkvGrL1UfHV0GTYs9oSPMEN0
HPHtXFd7jhEHqDkjChjYUfU0sXGZDW3zmjGoSJGsWxAdq8lUYRZwSwabYizhAbFK
pIPwB9eHuM/HZFPRvdjQGIq/ibyXr4bq+jBCzpdCMeO6iAWWMv18739AhitMzz8T
J5K6IUT5y2IWCjnh0Yu5ABIgFfWttGMNq1JSYfE0uc6fPjBpZwLGnyLinOladobI
bH/VqXL2UT7o/93/5qRIgEJShzHklfPXkdgt63CgclOKsek2IxDiiZT4DrCiJ84u
Ib7LIJ13CnQlCVneFHMNiSDRrJvDJFB8MxQSjIFLD1eduwDrMgc66tgPKEb5mJKG
RULZKOXEeRmUxX7QEIkJR7m9R4iTnzn7OkjBgegJZIPHPJOpX6Y2sMZpcVIoUXdA
NFHRrht3mW+OiFQ=
-----END CERTIFICATE-----
`)

	testClientKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDWyWQzCDp7t9JU
L8guzjb5NaoUSstXAq82tdgjIcW90U8p02o/7lIFBW2K492c+iSU6pnALZrzp1aL
VNuLISU8bWHMrg6euUpkbYhBGwJMmWOR74IBmoGot2YF3LC5XDFk5rtl6sQgG5Ce
XkDekuIM1x8F0JeKttpIWvmS2CZqtsvBbxY/Ca9/YnpwpYNXcMjHxVSDa7T04B4n
ie8ACxMBQBUwPnzNjLc+V8Y14SoJuBic7MYGB9dUf/VRLoa2G8WU9MBtDmnOvQer
kuBRtcaEqACiQL1zxDEXC2jEcYIgVEyJh+fVNSkGQofn+jAJltPbelE5R1a08UrR
oPLtSIyNAgMBAAECggEBALUV+mqkJ1qjcqsT1fzQU7zsp8aQALwNQVgpHF8SXDtb
OxkSa+QWtAQTvXV6BCATLcB3wsUqLhf7H5Y9JxQ4D8LQncIJhb4Ajl35kwUBFoEq
Wa5ydfOQJnzukw+iL0U4G1Tsy1Z0BoLjepxq7to4kGku/bLTWNDUtViHix9pKYqR
o3LxO3MsQ/lG3w1SEH7wkj+tTR5paP85iAObywoTVzQpwfPoBnhXngJ/fKhui8eO
ufpDATvG/BbbIsHSv2M1qTGSfwp8a0rAO6UbDl0nHMxIuMCg7hTgKwg3RUfHXRlS
rmDIC7C1Ddq17HP+bBuIbhvcHcwTpb1MsG+S7CJtLMECgYEA9rVhYoeFnv5W7h7x
dK94tz/kMOOPceP5WRPM6tLwVbsp+ACcA2gWw0Wf71cOcFmlIpAUiP1ebn/PWx4z
s1Gpaw0kWxIG1WGWGr/d4pZ4VxJ1iSZ2WOTsx5KIb0v6Xxaj7emIo1dPu3lz5cNY
u18FAiHjJnbdB7/Gl+iu3U2k4WUCgYEA3uA9Q5x/JeSPrLtVfttpGj/Mi7bPgUlI
z4+va25f3ohmnBus12kVXmdZzLJGRA05zr9KXU4CfK4STtdWoXAeDAeuXgMarbyJ
4z/xk58/uW5V1zUZd9UwpaUUzPDu+6vjIb0XYrJIzm4lYQnY1wTboIBEpJngtsqa
zcMwPoAQIAkCgYEA2jVHs4xGpYA0h10bF6f0T7DVNmCwCX4ol58pyjFUnZ9z2YVA
eMriB0lX0qvfe4PuyYlCgIAJvBaT4vXtqJd8D9GJ7HsfTDCKQZKewMFyIyGSkAJS
/wFMZKC4yCgdhWlTCSVb041wWlNsLTcBDolWtrIeZXEQwr/e+ZG2yMraIPkCgYEA
vZszQ3O9z7TUbfSpVVS/427nSuzpN2nrIXlxmQm7UYvlD2WT82YYocl24efAU2CV
D0g5sYsOHpfQR3Z24ryJM17NfnlRlwBQph3eHOJbyhsNuBoaYpHh4um/+mH2TfD7
N9awMGzP955I+nbwHGyrk63Lt+SZAaj3bZliT6mPDlECgYA75EYIOQeS3ZVnXKMc
o954gXZMjSY50cyAKJeoLQ9X6MyRhm/Q5hq+wf1B1gOYeotaKmRDIEo3VPLvxya1
YNg0Nn7C5Wrih5O5JthcxgS5CSYAVn267pAgrkxxhipTeu4ooA5biC6tCfl9OW5b
F3k6QVn+kFMDawTbPpNuofmx5Q==
-----END PRIVATE KEY-----
`)

	testCACert = []byte(`
-----BEGIN CERTIFICATE-----
MIIE0zCCAzugAwIBAgIQe2yY55CxyB5WeRygBI2soTANBgkqhkiG9w0BAQsFADCB
gTEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSswKQYDVQQLDCJjcHVn
dXk4M0BCcmlhbnMtTVNGVC1tYWNib29rLmxvY2FsMTIwMAYDVQQDDClta2NlcnQg
Y3B1Z3V5ODNAQnJpYW5zLU1TRlQtbWFjYm9vay5sb2NhbDAeFw0yMDExMjMyMjQ0
NDhaFw0zMDExMjMyMjQ0NDhaMIGBMR4wHAYDVQQKExVta2NlcnQgZGV2ZWxvcG1l
bnQgQ0ExKzApBgNVBAsMImNwdWd1eTgzQEJyaWFucy1NU0ZULW1hY2Jvb2subG9j
YWwxMjAwBgNVBAMMKW1rY2VydCBjcHVndXk4M0BCcmlhbnMtTVNGVC1tYWNib29r
LmxvY2FsMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAsFmuAz92hzOp
IqtyhwIjAf9JimHOOzV5LvAdcF/mj44HMtS5DeYqAZ8qQFBWGjXeeTd9L5SXNQSl
lZKrUmc4yDpVGp4FcRp1wIn1T14KG2taw4CFmRSnwHdOgU3N4e/OA3SmiLVEtecZ
p0+IqPM2MGpHP9URCzeXDly0lln0wRp04auwgGt3QbliOiTBBAsZI/Nh6eQJaaj6
zkGQjlsw+r+SvGRGfYmD0twUVYqqe9Y0dx12InHEzXk3cCti0h51q326TMwtRRJL
3LhL06H4soEG++3NpWUEj6/ZdngUkzzYKXg+T4zVv7Yd/J+jAuUH1JgAUBkcr1sc
jqENdAksAvGhLa2a6xOSYpqEObQN0ejDKiUJoJcaRlHoN19UwEZsAW33vqBxehmO
fiNKREUv/kIyDODj5B9aiC6MCyx98r0Ks9FKz4Em0Vu2K961C4kJFNyvs5hrPBLN
ykF36y5KC/JzEbUskxvTn71NjFYBecKnGa9Ocx4Dm1Dqd7n3Ru+FAgMBAAGjRTBD
MA4GA1UdDwEB/wQEAwICBDASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBRK
1Ydq70TqytTJntZ9Cn3Nfq66EjANBgkqhkiG9w0BAQsFAAOCAYEArFyHwcpQG98M
rJuPj68uL4awrL7WKLa8V4RrkQ7BiQOMjeN4AQx6YpLRa2nsZ1kfMiduy3hlGlSW
pPnJoHgH8sDgGSeUQ7DsaYXP0aqGsWwIqeqIp5LTdfZ9B8QJSbgsp9zlj2tlB/dv
fgy0XgfMGPjiu2xL4spYlUudQhzMu6Gy9FHEY9ug9Dyd8FRdYi0ExNgwKd/bVMeM
puyJ2L+FGvi9pvQeqKJMT+Blpi5R8Id2dqETgy+QuGbPuZJzkPf/Ft6RRZ/WTm1T
uBgd1dA7EBIlcBoQc9os/O4tRxAW3sz1ZHuscNTRLBN88cG7S16tCjgOJ3RhF1jy
BXlHRpas7mXHQPmS8Y1ElUeBSFCRLdAMLoNLJ8SYv+CrYl/xtWNJGJS4q9n4wral
2v63uZ53OPQXnYK0+6SQ6KqT0Q8hst2AK7Zg5I88/L5b0xPEz7jDsjJjpDrEaCfB
+CDS+4gGSOw/b+fykIT/vD5v+mnom4KX7Xdi/SV3DbbDnXkwH5W/
-----END CERTIFICATE-----
`)
)

var (
	calledAuthenticate = false
	calledAuthorize    = false
	calledAttributes   = false
)

func TestHTTPServer(t *testing.T) {
	pCfg := mock.Config{
		CPU:    "1",
		Memory: "100M",
		Pods:   "100",
	}

	p, err := mock.NewProviderConfig(pCfg, t.Name(), runtime.GOOS, "", 0)
	assert.NilError(t, err)

	dir, err := ioutil.TempDir("", strings.Replace(t.Name(), string(os.PathSeparator), "_", -1))
	assert.NilError(t, err)
	defer os.RemoveAll(dir)
	writeTestCerts(t, dir)

	cert := filepath.Join(dir, "cert.pem")
	key := filepath.Join(dir, "key.pem")
	clientCA := filepath.Join(dir, "client-ca.pem")
	clientKey := filepath.Join(dir, "client-key.pem")
	clientCert := filepath.Join(dir, "client-cert.pem")
	badPath := "/some/nonexistent/path"

	clientCAPool := x509.NewCertPool()
	assert.Assert(t, clientCAPool.AppendCertsFromPEM(testCACert), "could not add ca cert")

	unauthenticatedClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: clientCAPool,
			},
		},
	}

	clientAuthCert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	assert.NilError(t, err)

	authClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientAuthCert},
				RootCAs:      clientCAPool,
			},
		},
	}

	authorizedFakeFilter := &fakeAuth{
		authenticateFunc: func(req *http.Request) (*authenticator.Response, bool, error) {
			calledAuthenticate = true
			return &authenticator.Response{User: &user.DefaultInfo{Name: "test"}}, true, nil
		},
		attributesFunc: func(u user.Info, req *http.Request) authorizer.Attributes {
			calledAttributes = true
			return &authorizer.AttributesRecord{User: u}
		},
		authorizeFunc: func(a authorizer.Attributes) (decision authorizer.Decision, reason string, err error) {
			calledAuthorize = true
			return authorizer.DecisionAllow, "", nil
		},
	}

	unauthorizedFakeFilter := &fakeAuth{
		authenticateFunc: func(req *http.Request) (*authenticator.Response, bool, error) {
			calledAuthenticate = true
			return &authenticator.Response{User: &user.DefaultInfo{Name: "test"}}, true, nil
		},
		attributesFunc: func(u user.Info, req *http.Request) authorizer.Attributes {
			calledAttributes = true
			return &authorizer.AttributesRecord{User: u}
		},
		authorizeFunc: func(a authorizer.Attributes) (decision authorizer.Decision, reason string, err error) {
			calledAuthorize = true
			return authorizer.DecisionDeny, "", nil
		},
	}

	t.Run("tls", func(t *testing.T) {
		for name, cfg := range map[string]*apiServerConfig{
			"no config":                 {},
			"bad key path":              {KeyPath: badPath},
			"bad cert path":             {CertPath: badPath},
			"bad ca cert path":          {CACertPath: badPath},
			"bad ca cert and key path":  {KeyPath: badPath, CACertPath: badPath},
			"bad ca cert and cert path": {CertPath: badPath, CACertPath: badPath},
			"bad cert and key path":     {CertPath: badPath, KeyPath: badPath},
			"ok cert/key paths":         {CertPath: cert, KeyPath: key},
		} {
			cfg := cfg
			t.Run(name, func(t *testing.T) {
				closer := getTestHTTPServer(t, cfg, p)
				assert.NilError(t, err)
				defer closer()

				c, err := net.Dial("tcp", cfg.Addr)
				if c != nil {
					c.Close()
				}
				// no listener shold be running
				assert.ErrorContains(t, err, "connection refused")
			})

		}

		t.Run("all bad paths", func(t *testing.T) {
			// this is a special case, because the user provided all the certs,
			// but they are missing This probably should/could change, in that
			// perhaps invalid config like a key but not cert (or vice-cersa)
			// should error instead of just not listening.
			cfg := &apiServerConfig{
				KeyPath:    badPath,
				CertPath:   badPath,
				CACertPath: badPath,
			}

			_, err := setupHTTPServer(context.Background(), p, cfg)
			assert.Assert(t, os.IsNotExist(errors.Cause(err)), err)
		})

		t.Run("client verification", func(t *testing.T) {
			t.Run("auth not required", func(t *testing.T) {
				t.Run("no client ca but auth required", func(t *testing.T) {
					cfg := &apiServerConfig{
						KeyPath:  key,
						CertPath: cert,
					}
					defer getTestHTTPServer(t, cfg, p)()

					c, err := net.Dial("tcp", cfg.Addr)
					if c != nil {
						c.Close()
					}
					// no listener shold be running
					assert.ErrorContains(t, err, "connection refused")
				})

				t.Run("no client ca but auth not required", func(t *testing.T) {
					cfg := &apiServerConfig{
						KeyPath:                     key,
						CertPath:                    cert,
						AllowUnauthenticatedClients: true,
					}
					defer getTestHTTPServer(t, cfg, p)()

					resp, err := unauthenticatedClient.Get(fmt.Sprintf("https://%s/runningpods", cfg.Addr))
					assert.NilError(t, err)
					resp.Body.Close()
					assert.Equal(t, resp.StatusCode, 200, resp.Status)
				})
			})

			t.Run("authenticated required", func(t *testing.T) {
				cfg := &apiServerConfig{
					KeyPath:    key,
					CertPath:   cert,
					CACertPath: clientCA,
				}
				defer getTestHTTPServer(t, cfg, p)()

				t.Run("unauthenticated client", func(t *testing.T) {
					_, err := unauthenticatedClient.Get(fmt.Sprintf("https://%s/runningpods", cfg.Addr))
					assert.ErrorContains(t, err, "bad certificate")
				})
				t.Run("authenticated client", func(t *testing.T) {
					resp, err := authClient.Get(fmt.Sprintf("https://%s/runningpods", cfg.Addr))
					assert.NilError(t, err)
					resp.Body.Close()
					assert.Equal(t, resp.StatusCode, 200, resp.Status)
				})
			})
		})
	})

	t.Run("webhook auth middleware", func(t *testing.T) {
		cfg := &apiServerConfig{
			KeyPath:    key,
			CertPath:   cert,
			CACertPath: clientCA,
		}
		cfg.AuthWebhookEnabled = true

		type testCase struct {
			name           string
			authClient     *http.Client
			authWebhook    *fakeAuth
			expectedStatus int
		}

		for _, c := range []testCase{
			{
				name:           "unauthenticated client and token auth forbidden",
				authClient:     unauthenticatedClient,
				authWebhook:    unauthorizedFakeFilter,
				expectedStatus: http.StatusForbidden,
			},
			{
				name:           "unauthenticated client but token auth pass",
				authClient:     unauthenticatedClient,
				authWebhook:    authorizedFakeFilter,
				expectedStatus: http.StatusOK,
			},
			{
				name:           "authenticated client but token auth forbidden",
				authClient:     unauthenticatedClient,
				authWebhook:    unauthorizedFakeFilter,
				expectedStatus: http.StatusForbidden,
			},
			{
				name:           "authenticated client and token auth pass",
				authClient:     unauthenticatedClient,
				authWebhook:    authorizedFakeFilter,
				expectedStatus: http.StatusOK,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				cfg.Auth = c.authWebhook

				closer := getTestHTTPServer(t, cfg, p)
				assert.NilError(t, err)
				defer closer()

				resp, err := c.authClient.Get(fmt.Sprintf("https://%s/stats/summary", cfg.Addr))
				assert.NilError(t, err)
				assert.Equal(t, resp.StatusCode, c.expectedStatus, resp.Status)
				assert.Assert(t, calledAuthenticate)
				assert.Assert(t, calledAttributes)
				assert.Assert(t, calledAuthorize)
			})
		}
	})
}

func getTestHTTPServer(t *testing.T, cfg *apiServerConfig, p provider.Provider) func() {
	var (
		closer func()
		err    error
		ctx    = context.Background()
		port   = 11250
	)

	for i := 0; i < 100; i++ {
		cfg.Addr = fmt.Sprintf("127.0.0.1:%d", port+i)
		closer, err = setupHTTPServer(ctx, p, cfg)
		if err == nil {
			t.Log(cfg.Addr)
			return closer
		}
		if closer != nil {
			closer()
		}
	}

	t.Fatalf("%+v", err)
	return nil
}

func writeTestCerts(t *testing.T, dir string) {
	t.Helper()

	err := ioutil.WriteFile(filepath.Join(dir, "cert.pem"), testCert, 0600)
	assert.NilError(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "key.pem"), testKey, 0600)
	assert.NilError(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "client-cert.pem"), testClientCert, 0600)
	assert.NilError(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "client-key.pem"), testClientKey, 0600)
	assert.NilError(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "client-ca.pem"), testCACert, 0600)
	assert.NilError(t, err)
}

type fakeAuth struct {
	authenticateFunc func(*http.Request) (*authenticator.Response, bool, error)
	attributesFunc   func(user.Info, *http.Request) authorizer.Attributes
	authorizeFunc    func(authorizer.Attributes) (authorized authorizer.Decision, reason string, err error)
}

func (f *fakeAuth) AuthenticateRequest(req *http.Request) (*authenticator.Response, bool, error) {
	return f.authenticateFunc(req)
}
func (f *fakeAuth) GetRequestAttributes(u user.Info, req *http.Request) authorizer.Attributes {
	return f.attributesFunc(u, req)
}
func (f *fakeAuth) Authorize(ctx context.Context, a authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	return f.authorizeFunc(a)
}
