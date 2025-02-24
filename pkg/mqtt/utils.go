package mqtt

import "regexp"

func checkLength(start int, length int) error {
	if start >= length {
		return ErrBytesShorter
	}
	return nil
}

// ValidateMQTTTopic validates if the MQTT topic name is valid according to MQTT rules.
func ValidateMQTTTopic(topic string) bool {
	// Regular expression to validate MQTT topic format
	// This regex checks if the topic contains only alphanumeric characters, slashes, and special MQTT characters (+, #)
	re := regexp.MustCompile(`^([a-zA-Z0-9_+/#-]+)$`)

	// Ensure the topic length is within the allowed size
	if len(topic) > 65535 {
		return false
	}

	// Match the topic format
	return re.MatchString(topic)
}
