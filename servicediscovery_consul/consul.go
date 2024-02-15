package servicediscoveryconsul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	port      = 3004
	serviceId = "order-service"
)

func RegisterService() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	addr := "localhost"
	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    "order-server",
		Port:    port,
		Address: addr,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", addr, port, serviceId),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	log.Printf(fmt.Sprintf("%s:%d/%s", addr, port, serviceId))

	regiErr := consul.Agent().ServiceRegister(registration)

	if regiErr != nil {
		log.Printf("Failed to register service: %s:%v ", addr, port)
	} else {
		log.Printf("successfully register service: %s:%v", addr, port)
	}
}
