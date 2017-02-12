package retrieving

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/nokka/armory/character"
)

type retrieveCharacterRequest struct {
	CharacterName string
}

type retrieveCharacterResponse struct {
	Character *character.Character `json:"character"`
}

func makeRetrieveCharacterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(retrieveCharacterRequest)
		character, err := s.RetrieveCharacter(req.CharacterName)

		if err != nil {
			return nil, err
		}

		return retrieveCharacterResponse{Character: character}, nil
	}
}
