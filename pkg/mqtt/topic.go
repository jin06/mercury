package mqtt

import "strings"

const (
	SysTopic   TopicType = "sys"    // system topic type
	ShareType  TopicType = "share"  // share topic type
	CommonType TopicType = "common" // common topic type
)

type TopicType string

type Topic string

func (t *Topic) Valid() error {
	s := string(*t)
	if len(s) == 0 {
		return ErrNotValidTopic
	}
	parts := strings.Split(s, "/")
	for i, part := range parts {
		if strings.Contains(part, "+") && len(part) != 1 {
			return ErrNotValidTopic
		}
		if strings.Contains(part, "#") {
			if len(part) != 1 {
				return ErrNotValidTopic
			}
			if len(parts)-1 > i {
				return ErrNotValidTopic
			}
		}
	}
	return nil
}

func (t *Topic) IsWild() bool {
	s := string(*t)
	return strings.Contains(s, "+") || strings.Contains(s, "#")
}

func (t *Topic) Type() TopicType {
	s := string(*t)
	if strings.HasPrefix(s, "$SYS/") {
		return SysTopic
	}
	if strings.HasPrefix(s, "$share/") {
		return ShareType
	}
	return CommonType
}

func (t *Topic) IsCommon() bool {
	return t.Type() == CommonType
}

func (t *Topic) IsShare() bool {
	return t.Type() == ShareType
}

func ValidTopic(topic UTF8String) error {
	if len(topic) == 0 {
		return ErrNotValidTopic
	}
	s := string(topic)
	parts := strings.Split(s, "/")
	for i, part := range parts {
		if strings.Contains(part, "+") && len(part) != 1 {
			return ErrNotValidTopic
		}
		if strings.Contains(part, "#") {
			if len(part) != 1 {
				return ErrNotValidTopic
			}
			if len(parts)-1 > i {
				return ErrNotValidTopic
			}
		}
	}
	return nil
}

func (u *Topic) Encode() ([]byte, error) {
	return encodeUTF8Str(*u)
}

func (t *Topic) Decode(data []byte) (int, error) {
	str, n, err := decodeUTF8Str(data)
	if err != nil {
		return n, err
	}
	*t = Topic(str)
	return n, nil
}

func (t *Topic) String() string {
	return string(*t)
}

func (t *Topic) TopicFilter() string {
	s := string(*t)
	if t.IsShare() {
		parts := strings.SplitN(s, "/", 3)
		if len(parts) == 3 {
			return parts[2]
		}
	}
	return s
}

func (t *Topic) Split() []string {
	tf := t.TopicFilter()
	return strings.Split(tf, "/")
}
