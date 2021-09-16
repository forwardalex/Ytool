package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init(fp string) (err error) {
	viper.SetConfigFile(fp)

	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed,err:%v\n", err)
		return err
	}
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarsh failed,err:%v\n", err)

	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("conf was update")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		}
	})
	return
}
