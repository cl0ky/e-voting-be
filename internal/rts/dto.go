package rts

type RTItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type RTListResponse struct {
	RTs []RTItem `json:"rts"`
}
