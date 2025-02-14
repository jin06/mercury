package mqtt

type Will struct {
	Topic      string
	Message    string
	QoS        QoS
	Retain     bool
	Properties *Properties
}

type WillProperties struct {
	// DelayInterval in seconds
	DelayInterval          uint32
	PayloadFormatIndicator bool
	// PublicationExpiryInterval in seconds
	PublicationExpiryInterval uint32
	ContentType               string
	ResponseTopic             string
	CorrelationData           string
}

func (w *Will) Encode() ([]byte, error) {
	var data []byte

	// Encode Topic
	topicData, err := encodeUTF8Str(w.Topic)
	if err != nil {
		return nil, err
	}
	data = append(data, topicData...)

	// Encode Message
	messageData, err := encodeUTF8Str(w.Message)
	if err != nil {
		return nil, err
	}
	data = append(data, messageData...)

	// Encode QoS and Retain
	data = append(data, byte(w.QoS))
	if w.Retain {
		data[len(data)-1] |= 0b00000001
	}

	// Encode Properties (MQTT 5.0 only)
	if w.Properties != nil {
		propertiesData, err := w.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	return data, nil
}

func (w *Will) Decode(data []byte) (int, error) {
	var start int

	// Decode Topic
	topic, n, err := decodeUTF8Str(data[start:])
	if err != nil {
		return start, err
	}
	w.Topic = topic
	start += n

	// Decode Message
	message, n, err := decodeUTF8Str(data[start:])
	if err != nil {
		return start, err
	}
	w.Message = message
	start += n

	// Decode QoS and Retain
	w.QoS = QoS(data[start] & 0b00000011)
	w.Retain = data[start]&0b00000001 == 0b00000001
	start++

	// Decode Properties (MQTT 5.0 only)
	if len(data) > start {
		w.Properties = new(Properties)
		n, err := w.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start += n
	}

	return len(data), nil
}

func (w *Will) Read(r *Reader) error {
	return nil
}

func (w *Will) Write(wr *Writer) error {
	data, err := w.Encode()
	if err != nil {
		return err
	}
	_, err = wr.Write(data)
	return err
}
