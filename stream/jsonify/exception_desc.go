package jsonify

type DubboCode uint8

const (
	OK                                DubboCode = 20
	CHANNEL_INACTIVE                            = 35
	BAD_REQUEST                                 = 40
	CLIENT_TIMEOUT                              = 30
	SERVER_TIMEOUT                              = 31
	BAD_RESPONSE                                = 50
	SERVICE_NOT_FOUND                           = 60
	SERVICE_ERROR                               = 70
	SERVER_ERROR                                = 80
	CLIENT_ERROR                                = 90
	SERVER_THREADPOOL_EXHAUSTED_ERROR           = 100
)

var dubboExceptionDesc = []string{
	CLIENT_TIMEOUT:                    "client side timeout.",
	SERVER_TIMEOUT:                    "server side timeout.",
	CHANNEL_INACTIVE:                  "channel inactive, directly return the unfinished requests.",
	BAD_REQUEST:                       "request format error.",
	BAD_RESPONSE:                      "response format error.",
	SERVICE_NOT_FOUND:                 "service not found.",
	SERVICE_ERROR:                     "service error.",
	SERVER_ERROR:                      "internal server error.",
	CLIENT_ERROR:                      "internal server error.",
	SERVER_THREADPOOL_EXHAUSTED_ERROR: "server side threadpool exhausted and quick return.",
}

type KafakaCode int16

const (
	UNKNOWN_SERVER_ERROR                             = -1
	NONE                                  KafakaCode = 0
	OFFSET_OUT_OF_RANGE                              = 1
	CORRUPT_MESSAGE                                  = 2
	UNKNOWN_TOPIC_OR_PARTITION                       = 3
	INVALID_FETCH_SIZE                               = 4
	LEADER_NOT_AVAILABLE                             = 5
	NOT_LEADER_OR_FOLLOWER                           = 6
	REQUEST_TIMED_OUT                                = 7
	BROKER_NOT_AVAILABLE                             = 8
	REPLICA_NOT_AVAILABLE                            = 9
	MESSAGE_TOO_LARGE                                = 10
	STALE_CONTROLLER_EPOCH                           = 11
	OFFSET_METADATA_TOO_LARGE                        = 12
	NETWORK_EXCEPTION                                = 13
	COORDINATOR_LOAD_IN_PROGRESS                     = 14
	COORDINATOR_NOT_AVAILABLE                        = 15
	NOT_COORDINATOR                                  = 16
	INVALID_TOPIC_EXCEPTION                          = 17
	RECORD_LIST_TOO_LARGE                            = 18
	NOT_ENOUGH_REPLICAS                              = 19
	NOT_ENOUGH_REPLICAS_AFTER_APPEND                 = 20
	INVALID_REQUIRED_ACKS                            = 21
	ILLEGAL_GENERATION                               = 22
	INCONSISTENT_GROUP_PROTOCOL                      = 23
	INVALID_GROUP_ID                                 = 24
	UNKNOWN_MEMBER_ID                                = 25
	INVALID_SESSION_TIMEOUT                          = 26
	REBALANCE_IN_PROGRESS                            = 27
	INVALID_COMMIT_OFFSET_SIZE                       = 28
	TOPIC_AUTHORIZATION_FAILED                       = 29
	GROUP_AUTHORIZATION_FAILED                       = 30
	CLUSTER_AUTHORIZATION_FAILED                     = 31
	INVALID_TIMESTAMP                                = 32
	UNSUPPORTED_SASL_MECHANISM                       = 33
	ILLEGAL_SASL_STATE                               = 34
	UNSUPPORTED_VERSION                              = 35
	TOPIC_ALREADY_EXISTS                             = 36
	INVALID_PARTITIONS                               = 37
	INVALID_REPLICATION_FACTOR                       = 38
	INVALID_REPLICA_ASSIGNMENT                       = 39
	INVALID_CONFIG                                   = 40
	NOT_CONTROLLER                                   = 41
	INVALID_REQUEST                                  = 42
	UNSUPPORTED_FOR_MESSAGE_FORMAT                   = 43
	POLICY_VIOLATION                                 = 44
	OUT_OF_ORDER_SEQUENCE_NUMBER                     = 45
	DUPLICATE_SEQUENCE_NUMBER                        = 46
	INVALID_PRODUCER_EPOCH                           = 47
	INVALID_TXN_STATE                                = 48
	INVALID_PRODUCER_ID_MAPPING                      = 49
	INVALID_TRANSACTION_TIMEOUT                      = 50
	CONCURRENT_TRANSACTIONS                          = 51
	TRANSACTION_COORDINATOR_FENCED                   = 52
	TRANSACTIONAL_ID_AUTHORIZATION_FAILED            = 53
	SECURITY_DISABLED                                = 54
	OPERATION_NOT_ATTEMPTED                          = 55
	KAFKA_STORAGE_ERROR                              = 56
	LOG_DIR_NOT_FOUND                                = 57
	SASL_AUTHENTICATION_FAILED                       = 58
	UNKNOWN_PRODUCER_ID                              = 59
	REASSIGNMENT_IN_PROGRESS                         = 60
	DELEGATION_TOKEN_AUTH_DISABLED                   = 61
	DELEGATION_TOKEN_NOT_FOUND                       = 62
	DELEGATION_TOKEN_OWNER_MISMATCH                  = 63
	DELEGATION_TOKEN_REQUEST_NOT_ALLOWED             = 64
	DELEGATION_TOKEN_AUTHORIZATION_FAILED            = 65
	DELEGATION_TOKEN_EXPIRED                         = 66
	INVALID_PRINCIPAL_TYPE                           = 67
	NON_EMPTY_GROUP                                  = 68
	GROUP_ID_NOT_FOUND                               = 69
	FETCH_SESSION_ID_NOT_FOUND                       = 70
	INVALID_FETCH_SESSION_EPOCH                      = 71
	LISTENER_NOT_FOUND                               = 72
	TOPIC_DELETION_DISABLED                          = 73
	FENCED_LEADER_EPOCH                              = 74
	UNKNOWN_LEADER_EPOCH                             = 75
	UNSUPPORTED_COMPRESSION_TYPE                     = 76
	STALE_BROKER_EPOCH                               = 77
	OFFSET_NOT_AVAILABLE                             = 78
	MEMBER_ID_REQUIRED                               = 79
	PREFERRED_LEADER_NOT_AVAILABLE                   = 80
	GROUP_MAX_SIZE_REACHED                           = 81
	FENCED_INSTANCE_ID                               = 82
	ELIGIBLE_LEADERS_NOT_AVAILABLE                   = 83
	ELECTION_NOT_NEEDED                              = 84
	NO_REASSIGNMENT_IN_PROGRESS                      = 85
	GROUP_SUBSCRIBED_TO_TOPIC                        = 86
	INVALID_RECORD                                   = 87
	UNSTABLE_OFFSET_COMMIT                           = 88
	THROTTLING_QUOTA_EXCEEDED                        = 89
	PRODUCER_FENCED                                  = 90
	RESOURCE_NOT_FOUND                               = 91
	DUPLICATE_RESOURCE                               = 92
	UNACCEPTABLE_CREDENTIAL                          = 93
	INCONSISTENT_VOTER_SET                           = 94
	INVALID_UPDATE_VERSION                           = 95
	FEATURE_UPDATE_FAILED                            = 96
	PRINCIPAL_DESERIALIZATION_FAILURE                = 97
	SNAPSHOT_NOT_FOUND                               = 98
	POSITION_OUT_OF_RANGE                            = 99
	UNKNOWN_TOPIC_ID                                 = 100
	DUPLICATE_BROKER_REGISTRATION                    = 101
	BROKER_ID_NOT_REGISTERED                         = 102
	INCONSISTENT_TOPIC_ID                            = 103
	INCONSISTENT_CLUSTER_ID                          = 104
	TRANSACTIONAL_ID_NOT_FOUND                       = 105
)

const (
	UNKNOWN_SERVER_ERROR_DESC = "The server experienced an unexpected error when processing the request."
)

var kafkaExceptionDesc = []string{
	OFFSET_OUT_OF_RANGE:                   "The requested offset is not within the range of offsets maintained by the server.",
	CORRUPT_MESSAGE:                       "This message has failed its CRC checksum, exceeds the valid size, has a null key for a compacted topic, or is otherwise corrupt.",
	UNKNOWN_TOPIC_OR_PARTITION:            "This server does not host this topic-partition.",
	INVALID_FETCH_SIZE:                    "The requested fetch size is invalid.",
	LEADER_NOT_AVAILABLE:                  "There is no leader for this topic-partition as we are in the middle of a leadership election.",
	NOT_LEADER_OR_FOLLOWER:                "For requests intended only for the leader, this error indicates that the broker is not the current leader. For requests intended for any replica, this error indicates that the broker is not a replica of the topic partition.",
	REQUEST_TIMED_OUT:                     "The request timed out.",
	BROKER_NOT_AVAILABLE:                  "The broker is not available.",
	REPLICA_NOT_AVAILABLE:                 "The replica is not available for the requested topic-partition. Produce/Fetch requests and other requests intended only for the leader or follower return NOT_LEADER_OR_FOLLOWER if the broker is not a replica of the topic-partition.",
	MESSAGE_TOO_LARGE:                     "The request included a message larger than the max message size the server will accept.",
	STALE_CONTROLLER_EPOCH:                "The controller moved to another broker.",
	OFFSET_METADATA_TOO_LARGE:             "The metadata field of the offset request was too large.",
	NETWORK_EXCEPTION:                     "The server disconnected before a response was received.",
	COORDINATOR_LOAD_IN_PROGRESS:          "The coordinator is loading and hence can't process requests.",
	COORDINATOR_NOT_AVAILABLE:             "The coordinator is not available.",
	NOT_COORDINATOR:                       "This is not the correct coordinator.",
	INVALID_TOPIC_EXCEPTION:               "The request attempted to perform an operation on an invalid topic.",
	RECORD_LIST_TOO_LARGE:                 "The request included message batch larger than the configured segment size on the server.",
	NOT_ENOUGH_REPLICAS:                   "Messages are rejected since there are fewer in-sync replicas than required.",
	NOT_ENOUGH_REPLICAS_AFTER_APPEND:      "Messages are written to the log, but to fewer in-sync replicas than required.",
	INVALID_REQUIRED_ACKS:                 "Produce request specified an invalid value for required acks.",
	ILLEGAL_GENERATION:                    "Specified group generation id is not valid.",
	INCONSISTENT_GROUP_PROTOCOL:           "The group member's supported protocols are incompatible with those of existing members or first group member tried to join with empty protocol type or empty protocol list.",
	INVALID_GROUP_ID:                      "The configured groupId is invalid.",
	UNKNOWN_MEMBER_ID:                     "The coordinator is not aware of this member.",
	INVALID_SESSION_TIMEOUT:               "The session timeout is not within the range allowed by the broker (as configured by group.min.session.timeout.ms and group.max.session.timeout.ms).",
	REBALANCE_IN_PROGRESS:                 "The group is rebalancing, so a rejoin is needed.",
	INVALID_COMMIT_OFFSET_SIZE:            "The committing offset data size is not valid.",
	TOPIC_AUTHORIZATION_FAILED:            "Topic authorization failed.",
	GROUP_AUTHORIZATION_FAILED:            "Group authorization failed.",
	CLUSTER_AUTHORIZATION_FAILED:          "Cluster authorization failed.",
	INVALID_TIMESTAMP:                     "The timestamp of the message is out of acceptable range.",
	UNSUPPORTED_SASL_MECHANISM:            "The broker does not support the requested SASL mechanism.",
	ILLEGAL_SASL_STATE:                    "Request is not valid given the current SASL state.",
	UNSUPPORTED_VERSION:                   "The version of API is not supported.",
	TOPIC_ALREADY_EXISTS:                  "Topic with this name already exists.",
	INVALID_PARTITIONS:                    "Number of partitions is below 1.",
	INVALID_REPLICATION_FACTOR:            "Replication factor is below 1 or larger than the number of available brokers.",
	INVALID_REPLICA_ASSIGNMENT:            "Replica assignment is invalid.",
	INVALID_CONFIG:                        "Configuration is invalid.",
	NOT_CONTROLLER:                        "This is not the correct controller for this cluster.",
	INVALID_REQUEST:                       "This most likely occurs because of a request being malformed by the client library or the message was sent to an incompatible broker. See the broker logs for more details.",
	UNSUPPORTED_FOR_MESSAGE_FORMAT:        "The message format version on the broker does not support the request.",
	POLICY_VIOLATION:                      "Request parameters do not satisfy the configured policy.",
	OUT_OF_ORDER_SEQUENCE_NUMBER:          "The broker received an out of order sequence number.",
	DUPLICATE_SEQUENCE_NUMBER:             "The broker received a duplicate sequence number.",
	INVALID_PRODUCER_EPOCH:                "Producer attempted to produce with an old epoch.",
	INVALID_TXN_STATE:                     "The producer attempted a transactional operation in an invalid state.",
	INVALID_PRODUCER_ID_MAPPING:           "The producer attempted to use a producer id which is not currently assigned to its transactional id.",
	INVALID_TRANSACTION_TIMEOUT:           "The transaction timeout is larger than the maximum value allowed by the broker (as configured by transaction.max.timeout.ms).",
	CONCURRENT_TRANSACTIONS:               "The producer attempted to update a transaction while another concurrent operation on the same transaction was ongoing.",
	TRANSACTION_COORDINATOR_FENCED:        "Indicates that the transaction coordinator sending a WriteTxnMarker is no longer the current coordinator for a given producer.",
	TRANSACTIONAL_ID_AUTHORIZATION_FAILED: "Transactional Id authorization failed.",
	SECURITY_DISABLED:                     "Security features are disabled.",
	OPERATION_NOT_ATTEMPTED:               "The broker did not attempt to execute this operation. This may happen for batched RPCs where some operations in the batch failed, causing the broker to respond without trying the rest.",
	KAFKA_STORAGE_ERROR:                   "Disk error when trying to access log file on the disk.",
	LOG_DIR_NOT_FOUND:                     "The user-specified log directory is not found in the broker config.",
	SASL_AUTHENTICATION_FAILED:            "SASL Authentication failed.",
	UNKNOWN_PRODUCER_ID:                   "This exception is raised by the broker if it could not locate the producer metadata associated with the producerId in question. This could happen if, for instance, the producer's records were deleted because their retention time had elapsed. Once the last records of the producerId are removed, the producer's metadata is removed from the broker, and future appends by the producer will return this exception.",
	REASSIGNMENT_IN_PROGRESS:              "A partition reassignment is in progress.",
	DELEGATION_TOKEN_AUTH_DISABLED:        "Delegation Token feature is not enabled.",
	DELEGATION_TOKEN_NOT_FOUND:            "Delegation Token is not found on server.",
	DELEGATION_TOKEN_OWNER_MISMATCH:       "Specified Principal is not valid Owner/Renewer.",
	DELEGATION_TOKEN_REQUEST_NOT_ALLOWED:  "Delegation Token requests are not allowed on PLAINTEXT/1-way SSL channels and on delegation token authenticated channels.",
	DELEGATION_TOKEN_AUTHORIZATION_FAILED: "Delegation Token authorization failed.",
	DELEGATION_TOKEN_EXPIRED:              "Delegation Token is expired.",
	INVALID_PRINCIPAL_TYPE:                "Supplied principalType is not supported.",
	NON_EMPTY_GROUP:                       "The group is not empty.",
	GROUP_ID_NOT_FOUND:                    "The group id does not exist.",
	FETCH_SESSION_ID_NOT_FOUND:            "The fetch session ID was not found.",
	INVALID_FETCH_SESSION_EPOCH:           "The fetch session epoch is invalid.",
	LISTENER_NOT_FOUND:                    "There is no listener on the leader broker that matches the listener on which metadata request was processed.",
	TOPIC_DELETION_DISABLED:               "Topic deletion is disabled.",
	FENCED_LEADER_EPOCH:                   "The leader epoch in the request is older than the epoch on the broker.",
	UNKNOWN_LEADER_EPOCH:                  "The leader epoch in the request is newer than the epoch on the broker.",
	UNSUPPORTED_COMPRESSION_TYPE:          "The requesting client does not support the compression type of given partition.",
	STALE_BROKER_EPOCH:                    "Broker epoch has changed.",
	OFFSET_NOT_AVAILABLE:                  "The leader high watermark has not caught up from a recent leader election so the offsets cannot be guaranteed to be monotonically increasing.",
	MEMBER_ID_REQUIRED:                    "The group member needs to have a valid member id before actually entering a consumer group.",
	PREFERRED_LEADER_NOT_AVAILABLE:        "The preferred leader was not available.",
	GROUP_MAX_SIZE_REACHED:                "The consumer group has reached its max size.",
	FENCED_INSTANCE_ID:                    "The broker rejected this static consumer since another consumer with the same group.instance.id has registered with a different member.id.",
	ELIGIBLE_LEADERS_NOT_AVAILABLE:        "Eligible topic partition leaders are not available.",
	ELECTION_NOT_NEEDED:                   "Leader election not needed for topic partition.",
	NO_REASSIGNMENT_IN_PROGRESS:           "No partition reassignment is in progress.",
	GROUP_SUBSCRIBED_TO_TOPIC:             "Deleting offsets of a topic is forbidden while the consumer group is actively subscribed to it.",
	INVALID_RECORD:                        "This record has failed the validation on broker and hence will be rejected.",
	UNSTABLE_OFFSET_COMMIT:                "There are unstable offsets that need to be cleared.",
	THROTTLING_QUOTA_EXCEEDED:             "The throttling quota has been exceeded.",
	PRODUCER_FENCED:                       "There is a newer producer with the same transactionalId which fences the current one.",
	RESOURCE_NOT_FOUND:                    "A request illegally referred to a resource that does not exist.",
	DUPLICATE_RESOURCE:                    "A request illegally referred to the same resource twice.",
	UNACCEPTABLE_CREDENTIAL:               "Requested credential would not meet criteria for acceptability.",
	INCONSISTENT_VOTER_SET:                "Indicates that the either the sender or recipient of a voter-only request is not one of the expected voters.",
	INVALID_UPDATE_VERSION:                "The given update version was invalid.",
	FEATURE_UPDATE_FAILED:                 "Unable to update finalized features due to an unexpected server error.",
	PRINCIPAL_DESERIALIZATION_FAILURE:     "Request principal deserialization failed during forwarding. This indicates an internal error on the broker cluster security setup.",
	SNAPSHOT_NOT_FOUND:                    "Requested snapshot was not found.",
	POSITION_OUT_OF_RANGE:                 "Requested position is not greater than or equal to zero, and less than the size of the snapshot.",
	UNKNOWN_TOPIC_ID:                      "This server does not host this topic ID.",
	DUPLICATE_BROKER_REGISTRATION:         "This broker ID is already in use.",
	BROKER_ID_NOT_REGISTERED:              "The given broker ID was not registered.",
	INCONSISTENT_TOPIC_ID:                 "The log's topic ID did not match the topic ID in the request.",
	INCONSISTENT_CLUSTER_ID:               "The clusterId in the request does not match that found on the server.",
	TRANSACTIONAL_ID_NOT_FOUND:            "The transactionalId could not be found.",
}
