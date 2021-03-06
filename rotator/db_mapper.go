package rotator

import (
	"bytes"
	"encoding/base64"
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"github.com/pkg/errors"
)

type DbMapper struct{}

func (DbMapper) MapBase64ToCipherValue(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func (DbMapper) Map(value crypto.EncryptedValue) ([]byte, error) {
	var dbValue []byte
	dbValue = append(dbValue, value.Nonce...)
	dbValue = append(dbValue, value.Salt...)
	dbValue = append(dbValue, value.CipherValue...)

	dbValueWriter := bytes.NewBuffer([]byte{})
	base64Encoder := base64.NewEncoder(base64.StdEncoding, dbValueWriter)

	_, err := base64Encoder.Write(dbValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable to encode encrypted value into base64")
	}

	err = base64Encoder.Close()
	if err != nil {
		return nil, errors.Wrap(err, "unable to close base64 writer")
	}

	return dbValueWriter.Bytes(), nil
}
