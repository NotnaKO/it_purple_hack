package price_retrival

type Retriever struct {
}

type searchRequest struct {
	location *LocationNode
	category *CategoryNode
	userID   uint64
}

// Возвращает цену в копейках
func (r *Retriever) Search(info *ConnectionInfo) (uint64, error) {
	location := IDToLocationNodeMap[info.LocationID]
	category := IDToCategoryNodeMap[info.MicrocategoryID]
	return r.search(searchRequest{
		location: location,
		category: category,
		userID:   info.UserID,
	})
}

func (r *Retriever) search(request searchRequest) (uint64, error) {
	panic("unimplemented")
}
