package discovery

import (
	"context"
	"fmt"
	"github.com/dsxriiiii/l3x_pay/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("registered to consul")
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}

func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	addr, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addr) == 0 {
		return "", fmt.Errorf("got empty %s addrs from consul", serviceName)
	}
	i := rand.Intn(len(addr))
	logrus.Infof("Discovered %d instance of %s, addrs=%v", len(addr), serviceName, addr)
	return addr[i], nil

}
