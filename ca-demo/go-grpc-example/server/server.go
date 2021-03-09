package server


import "context"

type ProdService struct {}

func (p *ProdService) GetProdStock(context.Context, *ProdRequest) (*ProdResponse, error) {
	response := &ProdResponse{ProdStock: 30}

	return response, nil
}

