package models

type PostTextBody struct {
	Items []PostTextBodyItem `json:"items"`
}

type PostTextBodyItem struct {
	Source  string `json:"source"`
	Message string `json:"message"`
}

func MakePostTextBodyItem(source, message string) PostTextBodyItem {
	return PostTextBodyItem{
		Source:  source,
		Message: message,
	}
}

func MakePostTextBody(source string, messages ...string) PostTextBody {
	var items []PostTextBodyItem
	for _, message := range messages {
		items = append(items, MakePostTextBodyItem(source, message))
	}
	return PostTextBody{
		Items: items,
	}
}
