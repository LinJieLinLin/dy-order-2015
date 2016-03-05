package alipayModule

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRSA(t *testing.T) {

	alipay_config := &Alipay_config_struct{
		Partner:   "Partner",
		Key:       "Key",
		Sign_type: "RSA",
		Private_key_path: []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`),
		Ali_public_key_path: []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`),
		Input_charset: "utf-8",
		Cacert:        "Cacert",
		Transport:     "http",
	}

	Init(alipay_config)

	Convey("rsa Sign by private_key", t, func() {
			sign, err := rsaSign("Test.\n", alipay_config.Private_key)
			t.Logf("%v\n", sign)
			So(err, ShouldEqual, nil)
			So(sign, ShouldNotEqual, "")

			Convey("rsa verifySign by private_key", func() {
					encodeString, err := base64EnCode(sign)
					So(err, ShouldEqual, nil)
					t.Logf("encodeString : %v\n", encodeString)
					err = rsaVerify("Test.\n", encodeString, alipay_config.Public_key)
					So(err, ShouldEqual, nil)
				})
		})

	Convey("rsaDecrypt by private_key", t, func() {

			rsaString, err := RsaEncrypt("Test.\n", alipay_config.Public_key)
			So(err, ShouldEqual, nil)
			t.Logf("rsaString : %v", rsaString)

			Decrypt, err := rsaDecrypt(rsaString, alipay_config.Private_key)
			t.Logf("Decrypt : %v", Decrypt)
			So(err, ShouldEqual, nil)
			So(Decrypt, ShouldEqual, "Test.\n")
		})

}

func RsaEncrypt(origData string, public_key *rsa.PublicKey) (string, error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, public_key, []byte(origData))
	if nil != err {
		return "", err
	}
	return string(b), nil
}

//move to real code alipay_core.go
//func base64EnCode(sign string) (string, error) {
//	base64, err := base64.StdEncoding.DecodeString(sign)
//	if err != nil {
//		fmt.Errorf("base64 DecodeString error \n")
//		return "", err
//	}
//	return string(base64), nil
//}

func base64DeCode(sign string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		fmt.Errorf("base64 StdEncoding error \n")
		return "", err
	}
	return string(data), nil
}
