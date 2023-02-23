package shell

import "context"
import "encoding/json"

type GateWayList struct {
	Status string `json:"status"`
	Data   []Data `json:"data"`
}
type Data struct {
	UID       string      `json:"uid"`
	UserUID   string      `json:"user_uid"`
	Subdomain string      `json:"subdomain"`
	Gateway   string      `json:"gateway"`
	IsActive  bool        `json:"is_active"`
	ID        int         `json:"id"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	DeletedAt interface{} `json:"deleted_at"`
}

func (s *Shell) GetUserSubdomain(wallet string, source string) (subdomainList []string, err error) {
	url := "gateway/get_user_gateway"
	resq, err := s.Request(url).Option("wallet", wallet).Option("source", source).AclGet(context.Background())
	if err != nil {
		return subdomainList, err
	}
	var gatewayList GateWayList
	decoder := json.NewDecoder(resq.Output)
	decoder.Decode(&gatewayList)
	for _, gateway := range gatewayList.Data {
		if gateway.IsActive {
			subdomainList = append(subdomainList, gateway.Subdomain)
		}
	}
	return subdomainList, err
}
