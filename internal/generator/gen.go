package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"promogen/internal/utils"
)

const PromoURL string = "https://discord.com/billing/partner-promotions/"
const TokenEndPoint string = "https://api.discord.gx.games/v1/direct-fulfillment"

func (in *Instance) GetOperaToken() string {
	resp, err := in.Request(http.MethodPost, TokenEndPoint, Payload{
		"partnerUserId": utils.RandomID(),
	}, &Header{
		Cookie: in.Cookie,
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	var data OperaResponseToken
	if err = json.Unmarshal(in.Body, &data); err != nil {
		log.Println(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return data.Token
	default:
		fmt.Println(resp.Status, string(in.Body))
		return ""
	}
}
