package beacon_api

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

func (c *beaconApiValidatorClient) proposeAttestation(ctx context.Context, attestation *ethpb.Attestation) (*ethpb.AttestResponse, error) {
	if err := checkNilAttestation(attestation); err != nil {
		return nil, err
	}

	marshalledAttestation, err := json.Marshal(jsonifyAttestations([]*ethpb.Attestation{attestation}))
	if err != nil {
		return nil, err
	}

	if err = c.jsonRestHandler.Post(
		ctx,
		"/eth/v1/beacon/pool/attestations",
		nil,
		bytes.NewBuffer(marshalledAttestation),
		nil,
	); err != nil {
		return nil, err
	}

	attestationDataRoot, err := attestation.Data.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to compute attestation data root")
	}

	return &ethpb.AttestResponse{AttestationDataRoot: attestationDataRoot[:]}, nil
}

// checkNilAttestation returns error if attestation or any field of attestation is nil.
func checkNilAttestation(attestation ethpb.Att) error {
	if attestation == nil {
		return errors.New("attestation is nil")
	}

	if attestation.GetData() == nil {
		return errors.New("attestation data is nil")
	}

	if attestation.GetData().Source == nil || attestation.GetData().Target == nil {
		return errors.New("source/target in attestation data is nil")
	}

	if len(attestation.GetAggregationBits()) == 0 {
		return errors.New("attestation aggregation bits is empty")
	}

	if len(attestation.GetSignature()) == 0 {
		return errors.New("attestation signature is empty")
	}

	return nil
}

func (c *beaconApiValidatorClient) proposeAttestationElectra(ctx context.Context, attestation *ethpb.AttestationElectra) (*ethpb.AttestResponse, error) {
	if err := checkNilAttestation(attestation); err != nil {
		return nil, err
	}
	if len(attestation.CommitteeBits) == 0 {
		return nil, errors.New("attestation committee bits is empty")
	}

	marshalledAttestation, err := json.Marshal(jsonifyAttestationsElectra([]*ethpb.AttestationElectra{attestation}))
	if err != nil {
		return nil, err
	}

	if err = c.jsonRestHandler.Post(
		ctx,
		"/eth/v2/beacon/pool/attestations",
		nil,
		bytes.NewBuffer(marshalledAttestation),
		nil,
	); err != nil {
		return nil, err
	}

	attestationDataRoot, err := attestation.Data.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to compute attestation data root")
	}

	return &ethpb.AttestResponse{AttestationDataRoot: attestationDataRoot[:]}, nil
}
