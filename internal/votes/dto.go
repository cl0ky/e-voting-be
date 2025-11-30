package votes

type CommitVoteRequest struct {
	ElectionId string `json:"election_id" binding:"required,uuid"`
	HashVote   string `json:"hash_vote" binding:"required"`
}

type RevealVoteRequest struct {
	ElectionId  string `json:"election_id" binding:"required,uuid"`
	CandidateId string `json:"candidate_id" binding:"required,uuid"`
	Nonce       string `json:"nonce" binding:"required"`
}
