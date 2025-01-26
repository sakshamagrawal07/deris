package commands

const (
	GET        = "get"       //Get command
	GET_FORMAT = "get <key>" //Get command format

	EXIT = "exit" //Exit command

	EXPIRE_KEY        = "expire"                         //Expire a key after given time
	EXPIRE_KEY_FORMAT = "expire <key> <time in seconds>" //Expire key command format

	DELETE_KEY        = "delete"       //Delete key command
	DELETE_KEY_FORMAT = "delete <key>" //Delete key command format

	STRING_SET        = "set"               //Set command
	STRING_SET_FORMAT = "set <key> <value>" //Set command format

	STRING_SET_NOT_EXISTS        = "setnx"               //Set if not exists command
	STRING_SET_NOT_EXISTS_FORMAT = "setnx <key> <value>" // Set if not exists command format

	LIST_LEFT_PUSH        = "lpush"               //Left push command
	LIST_LEFT_PUSH_FORMAT = "lpush <key> <value>" //Left push command format

	LIST_LEFT_POP        = "lpop"       //Left pop command
	LIST_LEFT_POP_FORMAT = "lpop <key>" //Left pop command format

	LIST_RIGHT_PUSH        = "rpush"               //Right push command
	LIST_RIGHT_PUSH_FORMAT = "rpush <key> <value>" //Right push command format

	LIST_RIGHT_POP        = "rpop"       //Right pop command
	LIST_RIGHT_POP_FORMAT = "rpop <key>" //Right pop command format

	LIST_LENGTH        = "llen"       //Length of list command
	LIST_LENGTH_FORMAT = "llen <key>" //Length of list command format

	SUBSCRIBE_KEY        = "subscribe"       //Subscribe to a channel
	SUBSCRIBE_KEY_FORMAT = "subscribe <key>" //Subscribe command format
	
	UNSUBSCRIBE_KEY        = "unsubscribe"       //Unsubscribe to a channel
	UNSUBSCRIBE_KEY_FORMAT = "unsubscribe <key>" //Unsubscribe command format

	PUBLISH_KEY        = "publish"              //Publish to a channel
	PUBLISH_KEY_FORMAT = "publish <key> <value" //Publish command format

)
