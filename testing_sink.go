package gosteno

type TestingSink struct {
	Records []*Record
}

func EnterTestMode() {
	testSink := NewTestingSink()
	stenoConfig := Config{
		Sinks: []Sink{testSink},
	}
	Init(&stenoConfig)
}

func NewTestingSink() *TestingSink {
	return &TestingSink{
		Records: make([]*Record, 0, 10),
	}
}

func (tSink *TestingSink) AddRecord(record *Record) {
	tSink.Records = append(tSink.Records, record)
}

func (tSink *TestingSink) Flush() {

}

func (tSink *TestingSink) SetCodec(codec Codec) {

}

func (tSink *TestingSink) GetCodec() Codec {
	return nil
}
