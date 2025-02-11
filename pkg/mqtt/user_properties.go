package mqtt

type UserProperty struct {
	Key string
	Val string
}

type UserProperties []UserProperty

func (u *UserProperties) Encode() (result []byte, err error) {
	result = []byte{}
	for _, pro := range *u {
		data, err := encodeStringPair(pro.Key, pro.Val)
		if err != nil {
			return nil, err
		}
		result = append(result, ID_UserProperties)
		result = append(result, data...)
	}
	return
}

func (u *UserProperties) Read(r *Reader) error {
	key, _, err := r.ReadUTF8Str()
	if err != nil {
		return err
	}
	val, _, err := r.ReadUTF8Str()
	if err != nil {
		return err
	}
	*u = append(*u, UserProperty{key, val})

	return nil
}
