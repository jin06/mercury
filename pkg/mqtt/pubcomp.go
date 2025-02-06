package mqtt

import "io"

type Pubcomp struct {
	PacketID       PacketID
	Version        ProtocolVersion
	ReasonCode     ReasonCode
	ReasonString   string
	UserProperties UserProperties
}

func (p *Pubcomp) Encode() ([]byte, error) {
	// result := toHeader(PUBCOMP)
	// if p.Version == MQTT5 {
	// 	result[1] = 0b00000100
	// 	result = append(result, p.PacketID.ToBytes()...)
	// 	result = append(result, byte(p.ReasonCode), 0)
	// 	var length byte = 0
	// 	if reasonString, err := encodeUTF8Str(p.ReasonString); err != nil {
	// 		return []byte{}, nil
	// 	} else {
	// 		length += byte(len(reasonString))
	// 		result = append(result, reasonString...)
	// 	}

	// 	if userProperties, err := p.UserProperties.toBytes(); err != nil {
	// 		return []byte{}, nil
	// 	} else {
	// 		length += byte(len(userProperties))
	// 		result = append(result, userProperties...)
	// 	}
	// 	result[5] = length

	// } else {
	// 	result[1] = 0b00000010
	// 	result = append(result, p.PacketID.ToBytes()...)
	// }
	// return result, nil
	return nil, nil
}

func (p *Pubcomp) Decode(reader io.Reader) error {
	// if packetID, err := readUint16(reader); err != nil {
	// 	return err
	// } else {
	// 	p.PacketID = PacketID(packetID)
	// }
	// if p.Version == MQTT5 {
	// 	if reasonCode, err := readByte(reader); err != nil {
	// 		return err
	// 	} else {
	// 		p.ReasonCode = ReasonCode(reasonCode)
	// 	}
	// 	var userProperties UserProperties
	// 	if err := userProperties.Read(reader); err != nil {
	// 		return err
	// 	}
	// 	p.UserProperties = userProperties
	// }
	return nil
}
