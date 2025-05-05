package eventencoder

import (
	"encoding/base64"

	"github.com/aghaghiamh/gocast/QAGame/contract/goproto/matching"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/aghaghiamh/gocast/QAGame/pkg/typemapper"
	"google.golang.org/protobuf/proto"
)

func MatchedPlayerUsersDecoder(encodedProto string) (*matching.MatchedPlayers, error) {
	const op = richerr.Operation("brokermsg.MatchedPlayerUsersDecoder")

	payloadB, dErr := base64.StdEncoding.DecodeString(encodedProto)
	if dErr != nil {
		return nil, richerr.New(op).WithError(dErr)
	}

	pbMp := &matching.MatchedPlayers{}
	if err := proto.Unmarshal(payloadB, pbMp); err != nil {
		return nil, richerr.New(op).WithError(err)
	}

	return pbMp, nil
}

func MatchedPlayerUsersEncoder(mp entity.MatchedPlayers) (string, error) {
	const op = richerr.Operation("brokermsg.MatchedPlayerUsersEncoder")

	pbMp := matching.MatchedPlayers{
		Category: string(mp.Category),
		UserIds: typemapper.ArrayMapper(mp.UserIDs, func(uID uint) uint64 {
			return uint64(uID)
		}),
	}

	payload, err := proto.Marshal(&pbMp)
	if err != nil {
		return "", richerr.New(op).WithError(err)
	}
	payloadStr := base64.StdEncoding.EncodeToString(payload)

	return payloadStr, nil
}
