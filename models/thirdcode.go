package models

type IThirdCode interface {
	Get(buf []byte) string
	ReportErr()
}

type ThirdCode struct {
}
