package mqtt

type ReasonCode byte

// var (
// 	BASE_SUCCESS = 0x00
// )

var (
	V3_CONNACK_ACCEPT                ReasonCode = 0x00
	V3_CONNACK_UNACCEPTABLE_VERSION  ReasonCode = 0x01
	V3_CONNACK_INDENTIFIER_REJECT    ReasonCode = 0x02
	V3_CONNACK_SERVER_UNAVAILABLE    ReasonCode = 0x03
	V3_CONNACK_BAD_USERNAME_PASSWORD ReasonCode = 0x04
	V3_CONNACK_NOT_AUTHORIZED        ReasonCode = 0x05
)

var (
	V5_SUCCESS                                ReasonCode = 0x00
	V5_Normal_Disconnection                   ReasonCode = 0x00
	V5_Granted_QoS0                           ReasonCode = 0x00
	V5_Granted_QoS1                           ReasonCode = 0x01
	V5_Granted_QoS2                           ReasonCode = 0x02
	V5_Disconnect_With_Will_Message           ReasonCode = 0x04
	V5_No_Matching_Subscribers                ReasonCode = 0x10
	V5_No_Subscription_Existed                ReasonCode = 0x11
	V5_Continue_Authentication                ReasonCode = 0x18
	V5_ReAuthenticate                         ReasonCode = 0x19
	V5_Unspecified_Error                      ReasonCode = 0x80
	V5_Malformed_Packet                       ReasonCode = 0x81
	V5_Protocol_Error                         ReasonCode = 0x82
	V5_Implementation_Specific_Error          ReasonCode = 0x83
	V5_Unsupported_Protocol_Version           ReasonCode = 0x84
	V5_Client_Identifier_Not_Valid            ReasonCode = 0x85
	V5_Bad_User_Name_OR_Password              ReasonCode = 0x86
	V5_Not_Authorized                         ReasonCode = 0x87
	V5_Server_Unavailable                     ReasonCode = 0x88
	V5_Server_Busy                            ReasonCode = 0x89
	V5_Banned                                 ReasonCode = 0x8A
	V5_Server_Shutting_Down                   ReasonCode = 0x8B
	V5_Bad_Authentication_Method              ReasonCode = 0x8C
	V5_Keep_Alive_Timeout                     ReasonCode = 0x8D
	V5_Session_Taken_Over                     ReasonCode = 0x8E
	V5_Topic_Filter_Invalid                   ReasonCode = 0x8F
	V5_Topic_Name_Invalid                     ReasonCode = 0x90
	V5_Packet_Identifier_In_Use               ReasonCode = 0x91
	V5_Packet_Identifier_Not_Found            ReasonCode = 0x92
	V5_Receive_Maximum_Exceeded               ReasonCode = 0x93
	V5_Topic_Alias_Invalid                    ReasonCode = 0x94
	V5_Packet_Too_Large                       ReasonCode = 0x95
	V5_Message_Rate_Too_High                  ReasonCode = 0x96
	V5_Quota_Exceeded                         ReasonCode = 0x97
	V5_Administrative_Action                  ReasonCode = 0x98
	V5_Payload_Format_Invalid                 ReasonCode = 0x99
	V5_Retain_Not_Supported                   ReasonCode = 0x9A
	V5_QoS_Not_Supported                      ReasonCode = 0x9B
	V5_Use_Another_Server                     ReasonCode = 0x9C
	V5_Server_Moved                           ReasonCode = 0x9D
	V5_Shared_Subscriptions_Not_Supported     ReasonCode = 0x9E
	V5_Connection_Rate_Exceeded               ReasonCode = 0x9F
	V5_Maximum_Connect_Time                   ReasonCode = 0xA0
	V5_Subscription_Identifiers_Not_Supported ReasonCode = 0xA1
	V5_Wildcard_Subscriptions_Not_Supported   ReasonCode = 0xA2
)
