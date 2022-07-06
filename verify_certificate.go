package mgmt

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"sync"
)

func VerifyConnectionIgnoreExpiredCertificate(cs tls.ConnectionState) error {
	opts := x509.VerifyOptions{
		DNSName:       cs.ServerName,
		Intermediates: x509.NewCertPool(),
	}
	for _, cert := range cs.PeerCertificates[1:] {
		opts.Intermediates.AddCert(cert)
	}
	_, err := cs.PeerCertificates[0].Verify(opts)
	return err
}

// errNotParsed is returned when a certificate without ASN.1 contents is
// verified. Platform-specific verification needs the ASN.1 contents.
var errNotParsed = errors.New("x509: missing ASN.1 contents; use ParseCertificate")

const (
	leafCertificate = iota
	intermediateCertificate
	rootCertificate
)

var (
	once           sync.Once
	systemRoots    *x509.CertPool
	systemRootsErr error
)

func systemRootsPool() *x509.CertPool {
	once.Do(initSystemRoots)
	return systemRoots
}

func initSystemRoots() {
	// systemRoots, systemRootsErr = loadSystemRoots()
	systemRoots, systemRootsErr = x509.SystemCertPool()
	if systemRootsErr != nil {
		systemRoots = nil
	}
}

// SystemRootsError results when we fail to load the system root certificates.
type SystemRootsError struct {
	Err error
}

func (se SystemRootsError) Error() string {
	msg := "x509: failed to load system roots and no roots provided"
	if se.Err != nil {
		return msg + "; " + se.Err.Error()
	}
	return msg
}

func (se SystemRootsError) Unwrap() error { return se.Err }

func VerifyCertificateIgnoreExpired(c *x509.Certificate, opts x509.VerifyOptions) (chains [][]*x509.Certificate, err error) {
	// Platform-specific verification needs the ASN.1 contents so
	// this makes the behavior consistent across platforms.
	if len(c.Raw) == 0 {
		return nil, errNotParsed
	}
	// for i := 0; i < opts.Intermediates.len(); i++ {
	// 	c, err := opts.Intermediates.cert(i)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("crypto/x509: error fetching intermediate: %w", err)
	// 	}
	// 	if len(c.Raw) == 0 {
	// 		return nil, errNotParsed
	// 	}
	// }

	// // Use platform verifiers, where available, if Roots is from SystemCertPool.
	// if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
	// 	if opts.Roots == nil {
	// 		return c.systemVerify(&opts)
	// 	}
	// 	if opts.Roots != nil && opts.Roots.systemPool {
	// 		platformChains, err := c.systemVerify(&opts)
	// 		// If the platform verifier succeeded, or there are no additional
	// 		// roots, return the platform verifier result. Otherwise, continue
	// 		// with the Go verifier.
	// 		if err == nil || opts.Roots.len() == 0 {
	// 			return platformChains, err
	// 		}
	// 	}
	// }

	if opts.Roots == nil {
		opts.Roots = systemRootsPool()
		if err != nil {
			return
		}
		if opts.Roots == nil {
			return nil, x509.SystemRootsError{systemRootsErr}
		}
	}

	err = c.isValid(leafCertificate, nil, &opts)
	if err != nil {
		return
	}

	if len(opts.DNSName) > 0 {
		err = c.VerifyHostname(opts.DNSName)
		if err != nil {
			return
		}
	}

	var candidateChains [][]*x509.Certificate
	if opts.Roots.contains(c) {
		candidateChains = append(candidateChains, []*x509.Certificate{c})
	} else {
		if candidateChains, err = c.buildChains(nil, []*x509.Certificate{c}, nil, &opts); err != nil {
			return nil, err
		}
	}

	keyUsages := opts.KeyUsages
	if len(keyUsages) == 0 {
		keyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	// If any key usage is acceptable then we're done.
	for _, usage := range keyUsages {
		if usage == x509.ExtKeyUsageAny {
			return candidateChains, nil
		}
	}

	for _, candidate := range candidateChains {
		if checkChainForKeyUsage(candidate, keyUsages) {
			chains = append(chains, candidate)
		}
	}

	if len(chains) == 0 {
		return nil, x509.CertificateInvalidError{c, x509.IncompatibleUsage, ""}
	}

	return chains, nil
}

func checkChainForKeyUsage(chain []*x509.Certificate, keyUsages []x509.ExtKeyUsage) bool {
	usages := make([]x509.ExtKeyUsage, len(keyUsages))
	copy(usages, keyUsages)

	if len(chain) == 0 {
		return false
	}

	usagesRemaining := len(usages)

	// We walk down the list and cross out any usages that aren't supported
	// by each certificate. If we cross out all the usages, then the chain
	// is unacceptable.

NextCert:
	for i := len(chain) - 1; i >= 0; i-- {
		cert := chain[i]
		if len(cert.ExtKeyUsage) == 0 && len(cert.UnknownExtKeyUsage) == 0 {
			// The certificate doesn't have any extended key usage specified.
			continue
		}

		for _, usage := range cert.ExtKeyUsage {
			if usage == x509.ExtKeyUsageAny {
				// The certificate is explicitly good for any usage.
				continue NextCert
			}
		}

		const invalidUsage x509.ExtKeyUsage = -1

	NextRequestedUsage:
		for i, requestedUsage := range usages {
			if requestedUsage == invalidUsage {
				continue
			}

			for _, usage := range cert.ExtKeyUsage {
				if requestedUsage == usage {
					continue NextRequestedUsage
				}
			}

			usages[i] = invalidUsage
			usagesRemaining--
			if usagesRemaining == 0 {
				return false
			}
		}
	}

	return true
}
