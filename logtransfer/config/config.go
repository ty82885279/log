package config

type Config struct {
	Kafka KafkaCfg `ini:"kafka"`
	Es    EsCfg    `ini:"es"`
}
type KafkaCfg struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
type EsCfg struct {
	Address string `ini:"address"`
}
type LogData struct {
	Topic string
	Data  interface{}
}
