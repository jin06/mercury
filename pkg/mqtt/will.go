package mqtt

import "io"

type Will struct {
	Topic      string
	Message    string
	QoS        QoS
	Retain     bool
	Properties *Properties
}

type WillProperties struct {
	// DelayInterval  second
	DelayInterval          uint32
	PayloadFormatIndicator bool
	//PublicationExpiryInterval  second
	PublicationExpiryInterval uint32
	ContentType               string
	ResponseTopic             string
	CorrelationData           string
}

func decodeWillProperties(reader io.Reader) (result *WillProperties, err error) {
	// total, err := readByte(reader)
	// if err != nil {
	// 	return result, err
	// }
	// result = &WillProperties{}
	// for i := 0; i < int(total); {
	// 	var identifier byte
	// 	if identifier, err = readByte(reader); err != nil {
	// 		return
	// 	}
	// 	i++
	// 	switch identifier {
	// 	case 0x18:
	// 		{
	// 			i += 4
	// 			if result.DelayInterval, err = readUint32(reader); err != nil {
	// 				return
	// 			}
	// 		}
	// 	case 0x01:
	// 		// payload format
	// 		{
	// 			i++
	// 			if result.PayloadFormatIndicator, err = readBool(reader); err != nil {
	// 				return
	// 			}
	// 		}
	// 	case 0x02:
	// 		// publication expiry interval
	// 		{
	// 			i += 4
	// 			if result.PublicationExpiryInterval, err = readUint32(reader); err != nil {
	// 				return
	// 			}

	// 		}
	// 	case 0x03:
	// 		// content type
	// 		{
	// 			var n int
	// 			if result.ContentType, n, err = readStrN(reader); err != nil {
	// 				return
	// 			}
	// 			i += n
	// 		}
	// 	case 0x08:
	// 		// reponse topic
	// 		{
	// 			var n int
	// 			if result.ResponseTopic, n, err = readStrN(reader); err != nil {
	// 				return
	// 			}
	// 			i += n
	// 		}
	// 	case 0x09:
	// 		// correlation data
	// 		{
	// 			var n int
	// 			if result.CorrelationData, n, err = readStrN(reader); err != nil {
	// 				return
	// 			}
	// 			i += n
	// 		}
	// 	}
	// }
	return
}

func encodeWillProperties(will *WillProperties) (result []byte, err error) {
	if will == nil {
		return
	}
	return
}
