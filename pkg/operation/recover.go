/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package operation

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/trustbloc/sidetree-core-go/pkg/api/batch"
	"github.com/trustbloc/sidetree-core-go/pkg/api/protocol"
	"github.com/trustbloc/sidetree-core-go/pkg/docutil"
	"github.com/trustbloc/sidetree-core-go/pkg/restapi/model"
)

// ParseRecoverOperation will parse recover operation
func ParseRecoverOperation(request []byte, protocol protocol.Protocol) (*batch.Operation, error) {
	schema, err := parseRecoverRequest(request)
	if err != nil {
		return nil, err
	}

	code := protocol.HashAlgorithmInMultiHashCode

	delta, err := parseDelta(schema.Delta, code)
	if err != nil {
		return nil, err
	}

	signedData, err := parseSignedDataForRecovery(schema.SignedData.Payload, code)
	if err != nil {
		return nil, err
	}

	// TODO: Handle recovery key

	return &batch.Operation{
		OperationBuffer:              request,
		Type:                         batch.OperationTypeRecover,
		UniqueSuffix:                 schema.DidSuffix,
		Delta:                        delta,
		EncodedDelta:                 schema.Delta,
		RecoveryRevealValue:          schema.RecoveryRevealValue,
		UpdateCommitment:             delta.UpdateCommitment,
		RecoveryCommitment:           signedData.RecoveryCommitment,
		HashAlgorithmInMultiHashCode: code,
		SignedData:                   schema.SignedData,
	}, nil
}

func parseRecoverRequest(payload []byte) (*model.RecoverRequest, error) {
	schema := &model.RecoverRequest{}
	err := json.Unmarshal(payload, schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func parseDelta(encoded string, code uint) (*model.DeltaModel, error) {
	bytes, err := docutil.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	schema := &model.DeltaModel{}
	err = json.Unmarshal(bytes, schema)
	if err != nil {
		return nil, err
	}

	if err := validateDelta(schema, code); err != nil {
		return nil, err
	}

	return schema, nil
}

func parseSignedDataForRecovery(encoded string, code uint) (*model.RecoverSignedDataModel, error) {
	bytes, err := docutil.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	schema := &model.RecoverSignedDataModel{}
	err = json.Unmarshal(bytes, schema)
	if err != nil {
		return nil, err
	}

	if err := validateSignedDataForRecovery(schema, code); err != nil {
		return nil, err
	}

	return schema, nil
}

func validateSignedDataForRecovery(signedData *model.RecoverSignedDataModel, code uint) error {
	if signedData.RecoveryKey == nil {
		return errors.New("missing recovery key")
	}

	if !docutil.IsComputedUsingHashAlgorithm(signedData.RecoveryCommitment, uint64(code)) {
		return errors.New("next recovery commitment hash is not computed with the latest supported hash algorithm")
	}

	if !docutil.IsComputedUsingHashAlgorithm(signedData.DeltaHash, uint64(code)) {
		return errors.New("patch data hash is not computed with the latest supported hash algorithm")
	}

	return nil
}
