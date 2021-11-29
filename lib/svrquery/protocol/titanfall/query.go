package titanfall

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/multiplay/go-svrquery/lib/svrquery/common"
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
	"github.com/netdata/go-orchestrator/module"
)

var (

	// minLength is the smallest packet we can expect.
	minLength  = 26
	nonceSize  = 12
	tagSize    = 16
	packetSize = 1200
)

type queryer struct {
	c       protocol.Client
	version byte
}

func newQueryer(version byte) func(c protocol.Client) protocol.Queryer {
	return func(c protocol.Client) protocol.Queryer {
		return &queryer{
			c:       c,
			version: version,
		}
	}
}

func encrypt(b []byte) ([]byte, error) {
	key := os.Getenv("AES_KEY")
	if key == "" {
		return nil, fmt.Errorf("no aes key found. (Is AES_KEY env var set?)")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	hexKey := []byte(hex.EncodeToString(keyBytes))
	fmt.Println("key length", len(hexKey))
	fmt.Println("key", hexKey)

	c, err := aes.NewCipher(hexKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	fmt.Println("Nonce size", gcm.NonceSize())
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	fmt.Println("Nonce: ", nonce)

	// This will output in the form CipherTest | Tag and will need rearranging
	ciperTextAndTag := gcm.Seal(nil, nonce, b, nil)
	fmt.Println(ciperTextAndTag)

	// Rearange output to nonce | tag | ciphertext
	newCipherText := nonce
	newCipherText = append(newCipherText, ciperTextAndTag[len(ciperTextAndTag)-tagSize:]...)
	newCipherText = append(newCipherText, ciperTextAndTag[:len(ciperTextAndTag)-tagSize]...)

	fmt.Println(newCipherText)

	return newCipherText, nil
}

func decrypt(b []byte) ([]byte, error) {
	key := os.Getenv("AES_KEY")
	if key == "" {
		return nil, fmt.Errorf("no aes key found. (Is AES_KEY env var set?)")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher([]byte(hex.EncodeToString(keyBytes)))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	if len(b) < gcm.NonceSize() {
		return nil, fmt.Errorf("incoming bytes smaller than %d", gcm.NonceSize())
	}

	nonce, b := b[:gcm.NonceSize()], b[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, b, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Query implements protocol.Queryer.
func (q *queryer) Query() (protocol.Responser, error) {
	b := make([]byte, 1200)
	copy(b, q.serverInfoPkt())

	if key := q.c.Key(); key != "" {
		if q.version < 5 {
			// If keyed data asked for bump version sent to supported version level.
			b[5] = ServerInfoVersionKeyed
		}
		copy(b[6:], key)
	}

	b, err := encrypt(b)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Length: %d\n", len(b))

	if _, err := q.c.Write(b); err != nil {
		return nil, fmt.Errorf("query write: %w", err)
	}

	testRead := make([]byte, 16)
	n, err := q.c.Read(testRead)
	// Output everything
	fmt.Println(n, err, testRead)

	return nil, fmt.Errorf("debug")
	//
	//n, err := q.c.Read(b)
	//if err != nil {
	//	return nil, fmt.Errorf("query read: %w", err)
	//} else if n < minLength {
	//	return nil, fmt.Errorf("packet too short (len: %d)", n)
	//}
	//
	//b, err = decrypt(b)
	//if err != nil {
	//	return nil, err
	//}
	//
	//r := common.NewBinaryReader(b[:n], binary.LittleEndian)
	//i := &Info{}
	//
	//// Header.
	//if err = r.Read(&i.Header); err != nil {
	//	return nil, err
	//} else if i.Command != ServerInfoResponse {
	//	return nil, fmt.Errorf("unexpected cmd %x", i.Command)
	//}
	//
	//if i.Version > 1 {
	//	// InstanceInfo.
	//	if err = q.instanceInfo(r, i); err != nil {
	//		return nil, err
	//	}
	//}
	//
	//// BasicInfo.
	//if err = q.basicInfo(r, i); err != nil {
	//	return nil, err
	//}
	//
	//if i.Version > 4 {
	//	// PerformanceInfo.
	//	if err = r.Read(&i.PerformanceInfo); err != nil {
	//		return nil, err
	//	}
	//}
	//
	//if i.Version > 2 {
	//	if i.Version > 5 {
	//		// MatchState and Teams.
	//		if err = r.Read(&i.MatchState); err != nil {
	//			return nil, err
	//		}
	//	} else if err = r.Read(&i.MatchState.MatchStateV2); err != nil {
	//		return nil, err
	//	}
	//
	//	if err = q.teams(r, i); err != nil {
	//		return nil, err
	//	}
	//}
	//
	//// Clients
	//if err = q.clients(r, i); err != nil {
	//	return nil, err
	//}
	//
	//return i, nil
}

// instanceInfo decodes the instance information from a response.
func (q *queryer) instanceInfo(r *common.BinaryReader, i *Info) (err error) {
	if i.Version > 7 {
		if err = r.Read(&i.InstanceInfoV8); err != nil {
			return err
		}
	} else {
		if err = r.Read(&i.InstanceInfo); err != nil {
			return err
		}
	}
	if i.BuildName, err = r.ReadString(); err != nil {
		return err
	}
	if i.Datacenter, err = r.ReadString(); err != nil {
		return err
	}
	i.GameMode, err = r.ReadString()
	return err
}

// basicInfo decodes basic info from a response.
func (q *queryer) basicInfo(r *common.BinaryReader, i *Info) (err error) {
	if err = r.Read(&i.BasicInfo.Port); err != nil {
		return err
	} else if i.BasicInfo.Platform, err = r.ReadString(); err != nil {
		return err
	} else if i.BasicInfo.PlaylistVersion, err = r.ReadString(); err != nil {
		return err
	} else if err = r.Read(&i.BasicInfo.PlaylistNum); err != nil {
		return err
	} else if i.BasicInfo.PlaylistName, err = r.ReadString(); err != nil {
		return err
	}

	if i.Version > 6 {
		var platformNum byte

		if err = r.Read(&platformNum); err != nil {
			return err
		}
		i.BasicInfo.PlatformPlayers = make(map[string]byte, platformNum)

		for j := 0; j < int(platformNum); j++ {
			platformName, err := r.ReadString()
			if err != nil {
				return err
			}
			var platformPlayers byte
			if err = r.Read(&platformPlayers); err != nil {
				return err
			}
			i.BasicInfo.PlatformPlayers[platformName] = platformPlayers
		}
	}

	if err = r.Read(&i.BasicInfo.NumClients); err != nil {
		return err
	} else if err = r.Read(&i.BasicInfo.MaxClients); err != nil {
		return err
	}
	i.BasicInfo.Map, err = r.ReadString()
	return err
}

// teams decodes teams from a response.
func (q *queryer) teams(r *common.BinaryReader, i *Info) (err error) {
	var id byte
	if err = r.Read(&id); err != nil {
		return err
	} else if id != 255 {
		i.Teams = make([]Team, 0, 2)
		for id != 255 {
			t := Team{ID: id}
			if err = r.Read(&t.Score); err != nil {
				return err
			}
			i.Teams = append(i.Teams, t)

			if err = r.Read(&id); err != nil {
				return err
			}
		}
	}
	return nil
}

// clients decodes clients from a response.
func (q *queryer) clients(r *common.BinaryReader, i *Info) (err error) {
	var id uint64
	if err = r.Read(&id); err != nil {
		return err
	}

	for id > 0 {
		c := Client{ID: id}
		c.Name, err = r.ReadString()
		if err != nil {
			return err
		} else if err = r.Read(&c.TeamID); err != nil {
			return err
		}

		if i.Version > 3 {
			c.Address, err = r.ReadString()
			if err != nil {
				return err
			} else if err = r.Read(&c.Ping); err != nil {
				return err
			} else if err = r.Read(&c.PacketsReceived); err != nil {
				return err
			} else if err = r.Read(&c.PacketsDropped); err != nil {
				return err
			}
		}

		if i.Version > 2 {
			if err = r.Read(&c.Score); err != nil {
				return err
			} else if err = r.Read(&c.Kills); err != nil {
				return err
			} else if err = r.Read(&c.Deaths); err != nil {
				return err
			}
		}

		if err = r.Read(&id); err != nil {
			return err
		}
	}

	return nil
}

// Charts implements protocol.Charter.
func (q *queryer) Charts(serverID int64) module.Charts {
	if q.c.Key() == "" {
		return nil
	}

	cs := *charts.Copy()
	for _, c := range cs {
		c.ID = fmt.Sprintf(c.ID, serverID)
		c.Fam = fmt.Sprintf(c.Fam, serverID)
		for _, d := range c.Dims {
			d.ID = fmt.Sprintf(d.ID, serverID)
		}
		c.MarkNotCreated()
	}
	return cs
}

// serverInfoPkt returns a byte array of info request packet data.
func (q *queryer) serverInfoPkt() []byte {
	return []byte{0xFF, 0xFF, 0xFF, 0xFF, ServerInfoRequest, q.version}
}
