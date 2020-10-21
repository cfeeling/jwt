package jwt

import (
	"crypto/rand"
	"crypto/rsa"
)

type algRSAPSS struct {
	name string
	opts *rsa.PSSOptions
}

func (a *algRSAPSS) Name() string {
	return a.name
}

func (a *algRSAPSS) Sign(key PrivateKey, headerAndPayload []byte) ([]byte, error) {
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKey
	}

	h := a.opts.Hash.New()
	// header.payload
	_, err := h.Write(headerAndPayload)
	if err != nil {
		return nil, err
	}

	hashed := h.Sum(nil)
	return rsa.SignPSS(rand.Reader, privateKey, a.opts.Hash, hashed, a.opts)
}

func (a *algRSAPSS) Verify(key PublicKey, headerAndPayload []byte, signature []byte) error {
	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		if privateKey, ok := key.(*rsa.PrivateKey); ok {
			publicKey = &privateKey.PublicKey
		} else {
			return ErrInvalidKey
		}
	}

	h := a.opts.Hash.New()
	// header.payload
	_, err := h.Write(headerAndPayload)
	if err != nil {
		return err
	}

	hashed := h.Sum(nil)
	return rsa.VerifyPSS(publicKey, a.opts.Hash, hashed, signature, a.opts)
}
