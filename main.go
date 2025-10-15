package main

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	alidnsClient, err := NewAlidnsClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	ips, _ := GetAllIPv6()
	if err := alidnsClient.SetNewIPV6(ips); err != nil {
		fmt.Println("更新失败，请检查配置", err)
	}
	last := make([]string, 0)
	for range time.NewTicker(time.Minute).C {
		ips, _ := GetAllIPv6()
		slices.Sort(ips)
		if slices.Equal(ips, last) {
			continue
		}
		if err := alidnsClient.SetNewIPV6(ips); err != nil {
			fmt.Println("更新失败，请检查配置", err)
		} else {
			last = slices.Clone(ips)
			slices.Sort(last)
		}
	}
}
