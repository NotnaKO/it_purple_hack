package main

type Retriever struct {
}

func (r *Retriever) Search(info *ConnectionInfo) (float64, error) {
	return 100, nil
}
