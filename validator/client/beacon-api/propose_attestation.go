package beacon_api

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

func (c *beaconApiValidatorClient) proposeAttestation(ctx context.Context, attestation *ethpb.Attestation) (*ethpb.AttestResponse, error) {
	if err := validateNilAttestation(attestation); err != nil {
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

func (c *beaconApiValidatorClient) proposeAttestationElectra(ctx context.Context, attestation *ethpb.AttestationElectra) (*ethpb.AttestResponse, error) {
	if err := validateNilAttestation(attestation); err != nil {
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

func validateNilAttestation(attestation ethpb.Att) error {
	if attestation == nil || attestation.IsNil() {
		return errors.New("attestation can't be nil")
	}
	if attestation.GetData().Source == nil {
		return errors.New("attestation's source can't be nil")
	}
	if attestation.GetData().Target == nil {
		return errors.New("attestation's target can't be nil")
	}
	if attestation.GetAggregationBits() == nil {
		return errors.New("attestation's bitfield can't be nil")
	}
	return nil
}
