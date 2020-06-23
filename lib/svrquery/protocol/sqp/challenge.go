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

	return q.readChallenge()
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

// readChallenge reads the challenge response body and stores the challenge id
// for use in subsequent requests.
func (q *queryer) readChallenge() (err error) {
	q.challengeID, err = q.reader.ReadUint32()
	return err
}

// resetChallenge resets the challengeID so a new one will be generated when needed.
func (q *queryer) resetChallenge() {
	q.challengeID = 0
}
