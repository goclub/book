package func_user


func UserName(userID string) string {
	return "nimoc"
}

func MyMessageList(MessageListByUserID func(userID string) []string ) []string {
	return MessageListByUserID("a")
}