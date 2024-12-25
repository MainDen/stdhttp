package models

type ProcessesBody struct {
	Items []ProcessesBodyItem `json:"items"`
}

func (body ProcessesBody) ProcessModels() ProcessModels {
	var processes ProcessModels
	for _, item := range body.Items {
		processes = append(processes, item.ProcessModel())
	}
	return processes
}

type ProcessesBodyItem struct {
	Pid         int      `json:"pid"`
	ClientName  string   `json:"client_name"`
	CommandName string   `json:"command_name"`
	CommandArgs []string `json:"command_args"`
	Expired     bool     `json:"expired"`
	Persistent  bool     `json:"persistent"`
}

func (item ProcessesBodyItem) ProcessModel() ProcessModel {
	return ProcessModel{
		Pid:         item.Pid,
		ClientName:  item.ClientName,
		CommandName: item.CommandName,
		CommandArgs: item.CommandArgs,
		Persistent:  item.Persistent,
		Expired:     item.Expired,
	}
}

func MakeProcessesBodyItem(process ProcessModel) ProcessesBodyItem {
	return ProcessesBodyItem{
		Pid:         process.Pid,
		ClientName:  process.ClientName,
		CommandName: process.CommandName,
		CommandArgs: process.CommandArgs,
		Persistent:  process.Persistent,
		Expired:     process.Expired,
	}
}

func MakeProcessesBody(processes ...ProcessModel) ProcessesBody {
	var items []ProcessesBodyItem
	for _, process := range processes {
		items = append(items, MakeProcessesBodyItem(process))
	}
	return ProcessesBody{Items: items}
}
