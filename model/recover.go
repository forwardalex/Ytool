package model

type PanicReq struct {
	Service        string `json:"service"`
	ErrorInfo      string `json:"error_info"`
	Stack          string `json:"stack"`
	LogId          string `json:"log_id"`
	FuncName       string `json:"func_name"`
	Host           string `json:"host"`
	PodName        string `json:"pod_name"`
	LastCommitUser string `json:"last_commit_user"`
	CommitTime     string `json:"commit_time"`
	FileName       string `json:"file_name"`
}
