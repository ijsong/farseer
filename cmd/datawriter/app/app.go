package app

import (
	"github.com/ijsong/farseer/pkg/cassandra"
	"github.com/ijsong/farseer/pkg/kafka"
	"github.com/spf13/cobra"
)

func NewDataWriterCommand() *cobra.Command {
	cassandraConfig := &cassandra.CassandraStorageConfig{}
	conf := &DataWriterConfig{
		kafkaConfig:   &kafka.KafkaConfig{},
		storageConfig: cassandraConfig,
	}
	var cmd = &cobra.Command{
		Use:  "datawriter",
		Args: cobra.NoArgs,
	}
	var startCmd = &cobra.Command{
		Use:  "start",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(conf)
		},
	}
	startCmd.Flags().IntVar(&conf.kafkaConfig.FlushTimeoutMs, "kafka_flush_timeout", 1000, "flush timeout in milliseconds when closing kafka producer")
	startCmd.Flags().StringVar(&conf.kafkaConfig.BootstrapServers, "kafka_bootstrap_servers", "localhost", "kafka bootstrap servers")
	startCmd.Flags().StringVar(&conf.kafkaConfig.GroupID, "kafka_group_id", "datawriter", "kafka bootstrap servers")
	startCmd.Flags().StringSliceVar(&cassandraConfig.Hosts, "cassandra_hosts", []string{"localhost"}, "cassandra hosts")
	cmd.AddCommand(startCmd)

	return cmd
}

func run(conf *DataWriterConfig) error {
	app, err := NewDataWriter(conf)
	if err != nil {
		return err
	}
	return app.Start()
}
