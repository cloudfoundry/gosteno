package steno

import ()

type NullSink struct {
	records []*Record
}

func (nullSink *NullSink) AddRecord(record *Record) {
	nullSink.records = append(nullSink.records, record)
}

func (nullSink *NullSink) Flush() {

}

func (nullSink *NullSink) SetCodec(codec Codec) {

}

func (nullSink *NullSink) GetCodec() Codec {
	return nil
}

func (nullSink *NullSink) MarshalJSON() ([]byte, error) {
	return []byte(""), nil
}
