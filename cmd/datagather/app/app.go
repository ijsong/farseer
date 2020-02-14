package app

import (
	"time"

	"github.com/ijsong/farseer/pkg/kafka"

	"github.com/ijsong/farseer/pkg/queue"
	"github.com/ijsong/farseer/pkg/server"
	"github.com/spf13/cobra"
)

func NewDataGatherCommand() *cobra.Command {
	conf := &DataGatherConfig{
		serverConfig: &server.ServerConfig{},
		queueConfig:  &queue.EmbeddedQueueConfig{},
		kafkaConfig:  &kafka.KafkaConfig{},
	}

	var cmd = &cobra.Command{
		Use:  "datagather",
		Args: cobra.NoArgs,
	}
	var startCmd = &cobra.Command{
		Use:  "start",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(conf)
		},
	}

	startCmd.Flags().StringVar(&conf.serverConfig.ListenAddress, "listen_address", ":9091", "listen address")
	startCmd.Flags().StringVar(&conf.queueConfig.Address, "queue_address", "localhost:4150", "embedded queue address")
	startCmd.Flags().StringVar(&conf.queueConfig.DataPath, "queue_data_path", "data", "embedded queue data path")
	startCmd.Flags().Int64Var(&conf.queueConfig.MemQueueSize, "queue_mem_size", 10000, "embedded queue data path")
	startCmd.Flags().Int64Var(&conf.queueConfig.SyncEvery, "queue_sync_every", 2500, "number of messages per queue persistence (fsync)")
	startCmd.Flags().DurationVar(&conf.queueConfig.SyncTimeout, "queue_sync_timeout", 2*time.Second, "queue persistent interval (fsync)")
	startCmd.Flags().IntVar(&conf.queueConfig.NumberOfProducers, "queue_num_producers", 1, "the number of producers")
	startCmd.Flags().IntVar(&conf.queueConfig.NumberOfConsumers, "queue_num_consumers", 1, "the number of consumers")
	startCmd.Flags().IntVar(&conf.kafkaConfig.FlushTimeoutMs, "kafka_flush_timeout", 1000, "flush timeout in milliseconds when closing kafka producer")
	startCmd.Flags().StringVar(&conf.kafkaConfig.BootstrapServers, "kafka_bootstrap_servers", "localhost", "kafka bootstrap servers")
	cmd.AddCommand(startCmd)
	return cmd
}

func run(conf *DataGatherConfig) error {
	dataGather, err := NewDataGather(conf)
	if err != nil {
		return err
	}
	return dataGather.Start()
}
