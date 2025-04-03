package mqtt

func cloneBytePtr(b *byte) *byte {
	if b == nil {
		return nil
	}
	clone := *b
	return &clone
}

func cloneUint32Ptr(u *uint32) *uint32 {
	if u == nil {
		return nil
	}
	clone := *u
	return &clone
}

func cloneUint16Ptr(u *uint16) *uint16 {
	if u == nil {
		return nil
	}
	clone := *u
	return &clone
}

func cloneStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	clone := *s
	return &clone
}

func cloneBoolPtr(b *bool) *bool {
	if b == nil {
		return nil
	}
	clone := *b
	return &clone
}

func cloneQoSPtr(q *QoS) *QoS {
	if q == nil {
		return nil
	}
	clone := *q
	return &clone
}

func cloneUserProperties(up UserProperties) UserProperties {
	if up == nil {
		return nil
	}
	clone := make(UserProperties, len(up))
	for i, prop := range up {
		clone[i] = UserProperty{
			Key: prop.Key,
			Val: prop.Val,
		}
	}
	return clone
}
