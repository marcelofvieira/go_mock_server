package mockconfigurationrepo

import (
	"encoding/json"
	"errors"
	"mock_server_mux/internal/core/domain"
	"mock_server_mux/pkg/apperrors"
)

type MemoryKvs struct {
	kvs map[int][]byte
}

func NewMemKVS() *MemoryKvs {
	return &MemoryKvs{kvs: map[int][]byte{}}
}

func (repo *MemoryKvs) GetAll() ([]domain.MockConfiguration, error) {
	var mockConfigurations []domain.MockConfiguration

	if len(repo.kvs) == 0 {
		return mockConfigurations, apperrors.New(apperrors.NotFound, nil, "not found")

	} else {
		for _, value := range repo.kvs {
			mockConfiguration := domain.MockConfiguration{}

			err := json.Unmarshal(value, &mockConfiguration)
			if err != nil {
				return []domain.MockConfiguration{}, errors.New("fail to get value from kvs")
			}

			mockConfigurations = append(mockConfigurations, mockConfiguration)
		}
		return mockConfigurations, nil
	}

}

func (repo *MemoryKvs) GetById(Id int) (domain.MockConfiguration, error) {
	if value, ok := repo.kvs[Id]; ok {
		mockConfiguration := domain.MockConfiguration{}
		err := json.Unmarshal(value, &mockConfiguration)
		if err != nil {
			return domain.MockConfiguration{}, errors.New("fail to get value from kvs")
		}
		return mockConfiguration, nil
	} else {
		return domain.MockConfiguration{}, apperrors.New(apperrors.NotFound, nil, "mock configuration not found")
	}
}

func (repo *MemoryKvs) Save(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {
	newId := len(repo.kvs) + 1

	mockConfig.Id = newId

	bytes, err := json.Marshal(mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, apperrors.New(apperrors.InvalidInput, err, "fails at marshal into json string")
	}
	repo.kvs[newId] = bytes

	mockConfig.Id = newId

	return mockConfig, nil
}

func (repo *MemoryKvs) Update(mockConfig domain.MockConfiguration) (domain.MockConfiguration, error) {
	_, err := repo.GetById(mockConfig.Id)
	if err != nil {
		return domain.MockConfiguration{}, err
	}

	bytes, err := json.Marshal(mockConfig)
	if err != nil {
		return domain.MockConfiguration{}, apperrors.New(apperrors.InvalidInput, err, "fails at marshal into json string")
	}

	repo.kvs[mockConfig.Id] = bytes

	return mockConfig, nil
}
