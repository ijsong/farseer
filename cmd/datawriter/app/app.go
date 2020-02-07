package app

import "github.com/spf13/cobra"

func NewDataWriterCommand() *cobra.Command {
	conf := &DataWriterConfig{}
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
	startCmd.Flags().StringSliceVar(&conf.cassandraConfig.Hosts, "cassandra_hosts", []string{"localhost"}, "cassandra hosts")
	cmd.AddCommand(startCmd)
	return cmd
}

func run(conf *DataWriterConfig) error {
	return nil
}
