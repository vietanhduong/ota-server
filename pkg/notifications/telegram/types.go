package telegram

type RequestPayload struct {
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type ResponseBody struct {
	Ok     bool           `json:"ok"`
	Result *ResponseResult `json:"result"`
}

type ResponseResult struct {
	MessageId int               `json:"message_id"`
	From      *ResponseForm     `json:"from"`
	Chat      *ResponseChat     `json:"chat"`
	Date      int               `json:"date"`
	Text      string            `json:"text"`
	Entities  []*ResponseEntity `json:"entities"`
}

type ResponseForm struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type ResponseChat struct {
	Id                         int    `json:"id"`
	Title                      string `json:"title"`
	Type                       string `json:"type"`
	AllMemberAreAdministrators bool   `json:"all_member_are_administrators"`
}

type ResponseEntity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}
