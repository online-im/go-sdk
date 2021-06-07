package pkg

type ConnRsp struct {
	Address string `json:"address"`
	Ok      bool   `json:"ok"`
}

type WrappedGloryHttpRsp struct {
	Result ConnRsp `json:"result"`
	Ok     bool    `json:"ok"`
}
