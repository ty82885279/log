package config

type App struct {
	Kafka KafkaConfig `ini:"kafka"`
	Etcd  EtcdConfig  `ini:"etcd"`
}
type KafkaConfig struct {
	Address string `ini:"address"`
	MaxSize int    `ini:"chan_max_size"`
}
type EtcdConfig struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeout"`
}
type LogEntry struct {
	Path  string `ini:"path"`
	Topic string `ini:"topic"`
}
