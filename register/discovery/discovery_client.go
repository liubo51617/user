package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//service instance struct
type InstanceInfo struct {
	ID                string            `json:"id"`                  //服务ID
	Service           string            `json:"service,omitempty"`   //服务发现时返回的服务名
	Name              string            `json:"name"`                //服务名
	Tags              string            `json:"tags,omitempty"`      //标签，可用于服务过滤
	Address           string            `json:"address"`             //服务实例HOST
	Port              int               `json:"port"`                //服务实例端口
	Meta              map[string]string `json:"meta,omitempty"`      //元数据
	EnableTagOverride bool              `json:"enable_tag_override"` //是否允许标签覆盖
	Check             `json:"check,omitempty"`                       //健康检查相关配置
	Weights           `json:"weights,omitempty"`                     //权重
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"` // 多久之后注销服务
	Args                           []string `json:"Args,omitempty"`                 // 请求参数
	HTTP                           string   `json:"HTTP"`                           // 健康检查地址
	Interval                       string   `json:"Interval,omitempty"`             // Consul 主动检查间隔
	TTL                            string   `json:"TTL,omitempty"`                  // 服务实例主动维持心跳间隔，与Interval只存其一
}

type Weights struct {
	Passing int `json:"Passing"`
	Warning int `json:"Warning"`
}

type DiscoveryClient struct {
	host string //consul的host
	port int    //consul 的port
}

func NewDiscoveryClient(host string, port int) *DiscoveryClient {
	return &DiscoveryClient{
		host: host,
		port: port,
	}
}

//服务注册
func (discoveryclient *DiscoveryClient) Register(ctx context.Context, serviceName, instanceId, healthCheckUrl, instanceHost string, instancePort int, meta map[string]string, weights *Weights) error {
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}
	if weights != nil {
		instanceInfo.Weights = *weights
	} else {
		instanceInfo.Weights = Weights{
			Passing: 10,
			Warning: 1,
		}
	}

	byteData, err := json.Marshal(instanceInfo)

	if err != nil {
		log.Print("json marshal err: %s", err)
		return err
	}

	req, err := http.NewRequest("PUT", "http://"+discoveryclient.host+":"+strconv.Itoa(discoveryclient.port)+"/v1/agent/service/register", bytes.NewReader(byteData))
	if err != nil {
		log.Println("http request err： %s", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("register service err : %s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("register service http request errCode : %v", resp.StatusCode)
		return fmt.Errorf("register service http request errCode : %v", resp.StatusCode)
	}

	log.Println("register service success")
	return nil
}

//服务注销
func (discoveryClient *DiscoveryClient) Deregister(ctx context.Context, instanceId string) error {
	req, err := http.NewRequest("PUT", "http://"+discoveryClient.host+":"+strconv.Itoa(discoveryClient.port)+"/v1/agent/service/deregister/"+instanceId, nil)
	if err != nil {
		log.Println("http request err： %s", err)
		return err
	}
	client := http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
	if err != nil {
		log.Println("http response err： %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("deresigister service http request errCode : %v", resp.StatusCode)
		return fmt.Errorf("deresigister service http request errCode : %v", resp.StatusCode)
	}

	log.Println("deregister service success")
	return nil
}

//服务发现，通过serviceName
func (discoveryClient *DiscoveryClient) DiscoveryServices(ctx context.Context, serverName string) ([]*InstanceInfo, error) {
	req, err := http.NewRequest("GET", "http://"+discoveryClient.host+":"+strconv.Itoa(discoveryClient.port)+"/v1/health/service/"+serverName, nil)
	if err != nil {
		log.Println("http request err： %s", err)
		return nil, err
	}
	client := http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
	if err != nil {
		log.Println("http response err： %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("discover service http request errCode : %v", resp.StatusCode)
		return nil, fmt.Errorf("discover service http request errCode : %v", resp.StatusCode)
	}
	var serviceList []struct {
		Service InstanceInfo `json:"Service"`
	}
	err = json.NewDecoder(resp.Body).Decode(&serviceList)
	if err != nil {
		log.Printf("format service info err : %s", err)
		return nil, err
	}

	instances := make([]*InstanceInfo, len(serviceList))
	for i := 0; i < len(instances); i++ {
		instances[i] = &serviceList[i].Service
	}
	return instances, nil
}
