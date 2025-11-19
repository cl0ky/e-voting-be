package rts

type UseCase interface {
	GetAllRT() ([]RTItem, error)
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (u *useCase) GetAllRT() ([]RTItem, error) {
	rts, err := u.repo.GetAllRT()
	if err != nil {
		return nil, err
	}
	var result []RTItem
	for _, rt := range rts {
		result = append(result, RTItem{
			ID:     rt.Id.String(),
			Name:   rt.Name,
			Region: rt.Region,
		})
	}
	return result, nil
}
