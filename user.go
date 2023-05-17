package shell

import "context"
import "encoding/json"
import "fmt"
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

type IpfsFileServer struct {
	Status string             `json:"status"`
	Data   IpfsFileServerData `json:"data"`
}
type IpfsFileServerData struct {
	FileIpfsCid     string `json:"file_ipfsCid"`
	DownloadAddress string `json:"download_address"`
	DownloadURL     string `json:"download_url"`
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
			subdomainList = append(subdomainList, gateway.Gateway)
		}
	}
	return subdomainList, err
}

func (s *Shell) GetIpfsFileServer(ipfsCid string) (serverData IpfsFileServerData, err error) {
	url := "file/get_server"
	resq, err := s.Request(url).Option("ipfsCid", ipfsCid).AclGet(context.Background())
	if err != nil {
		return serverData, err
	}
	if resq.Output == nil {
		return serverData, fmt.Errorf("resq.Output is nil")
	}
	decoder := json.NewDecoder(resq.Output)
	var server IpfsFileServer
	decoder.Decode(&server)
	serverData = server.Data
	return serverData, err
}
