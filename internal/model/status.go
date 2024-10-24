package model

type transferProductStatus struct {
	OTW       string
	Failed    string
	Canceled  string
	Completed string
}

var TransferProductStatus transferProductStatus = transferProductStatus{
	OTW:       "OTW",
	Failed:    "Failed",
	Canceled:  "Canceled",
	Completed: "Completed",
}
