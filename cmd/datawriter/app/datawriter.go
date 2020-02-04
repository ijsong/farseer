package app

type DataWriter struct {
	conf *DataWriterConfig
}

func NewDataWriter(conf *DataWriterConfig) (*DataWriter, error) {
	dw := &DataWriter{
		conf: conf,
	}
	return dw, nil
}

func (dw *DataWriter) Start() error {
	return nil
}

func (dw *DataWriter) Stop() error {
	return nil
}
