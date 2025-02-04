package mqtt

func checkLength(start int, length int) error {
	if start >= length {
		return ErrBytesShorter
	}
	return nil
}
