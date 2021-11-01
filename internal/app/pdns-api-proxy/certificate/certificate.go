package certificate

import (
	"os"
	"path/filepath"

	"github.com/mvmaasakkers/certificates/cert"
	"github.com/nameserver-systems/pdns-distribute/internal/app/pdns-api-proxy/config"
)

func EnsureCertificatePaths(fqdn string, serviceconfig *config.ServiceConfiguration) (certpath, keypath string,
	err error) {

	certpath = serviceconfig.CertPath
	keypath = serviceconfig.KeyPath

	if certpath == "" || keypath == "" {
		userdir, direrr := os.UserHomeDir()
		if direrr != nil {
			return "", "", direrr
		}

		certpath = filepath.Join(userdir, "cert.pem")
		keypath = filepath.Join(userdir, "key.pem")
	}

	if !filepath.IsAbs(certpath) || !filepath.IsAbs(keypath) {
		return "", "", errCertPathNotAbsolute
	}

	if !fileExist(certpath) || !fileExist(keypath) {
		certerr := generateCertificates(certpath, keypath, fqdn)
		if certerr != nil {
			return "", "", certerr
		}
	}

	return certpath, keypath, nil
}

func generateCertificates(certificatepath, keypath, fqdn string) error {
	certrequest := createCertRequest(fqdn)

	basepath := filepath.Dir(certificatepath)

	cacertpath := filepath.Join(basepath, "cacert.pem")
	cakeypath := filepath.Join(basepath, "cakey.pem")

	cacertdata, cakeydata, caerr := cert.GenerateCA(certrequest)
	if caerr != nil {
		return caerr
	}

	err := writeCAFiles(cacertpath, cacertdata, cakeypath, cakeydata)
	if err != nil {
		return err
	}

	certdata, keydata, certerr := cert.GenerateCertificate(certrequest, cacertdata, cakeydata)
	if certerr != nil {
		return caerr
	}

	err2 := writeCertFiles(certificatepath, certdata, keypath, keydata)
	if err2 != nil {
		return err2
	}

	return nil
}

func writeCertFiles(certificatepath string, cert []byte, keypath string, key []byte) error {
	wrcerterr := os.WriteFile(certificatepath, cert, 0o600)
	if wrcerterr != nil {
		return wrcerterr
	}

	wrkeyerr := os.WriteFile(keypath, key, 0o600)
	if wrkeyerr != nil {
		return wrkeyerr
	}

	return nil
}

func writeCAFiles(cacertpath string, cacert []byte, cakeypath string, cakey []byte) error {
	wrcacerterr := os.WriteFile(cacertpath, cacert, 0o600)
	if wrcacerterr != nil {
		return wrcacerterr
	}

	wrcakeyerr := os.WriteFile(cakeypath, cakey, 0o600)

	if wrcakeyerr != nil {
		return wrcakeyerr
	}

	return nil
}

func createCertRequest(fqdn string) *cert.Request {
	certrequest := cert.NewRequest()
	certrequest.CommonName = fqdn
	certrequest.SubjectAltNames = []string{fqdn}

	return certrequest
}

func fileExist(filepath string) bool {
	_, err := os.Stat(filepath)

	return !os.IsNotExist(err)
}
