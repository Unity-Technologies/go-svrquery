package sqp

import (
	"bytes"
	"encoding/binary"
)

// Challenge sends a challenge request and validates a response
func (q *queryer) Challenge() error {
	if err := q.sendChallenge(); err != nil {
		return err
	}

	pktType, err := q.reader.ReadByte()
	if err != nil {
		return err
	} else if pktType != ChallengeResponseType {
		return NewErrMalformedPacketf("was expecting %v for response type, got %v", ChallengeResponseType, pktType)
	}

	q.challengeID, err = q.readChallenge()
	return err
}

// sendChallenge writes a challenge request
func (q *queryer) sendChallenge() error {
	pkt := &bytes.Buffer{}
	if err := binary.Write(pkt, binary.BigEndian, ChallengeRequestType); err != nil {
		return err
	}

	// Add 4 bytes of padding to make the request equal in size to the response so
	// these requests aren't attractive amplication vectors
	if _, err := pkt.Write([]byte{0, 0, 0, 0}); err != nil {
		return err
	}

	_, err := q.c.Write(pkt.Bytes())
	return err
}

// readChallenge reads the challenge from the response body and returns it.
func (q *queryer) readChallenge() (uint32, error) {
	return q.reader.ReadUint32()
}

// validateChallenge reads and validates the challenge of a request against our current challengeID.
func (q *queryer) validateChallenge() error {
	if id, err := q.readChallenge(); err != nil {
		return err
	} else if id != q.challengeID {
		return NewErrMalformedPacketf("unexpected challengeID 0x%0x wanted 0x%0x", id, q.challengeID)
	}
	return nil
}
