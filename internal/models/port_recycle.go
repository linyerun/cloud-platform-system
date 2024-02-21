package models

const (
	SinglePortType      = "single"
	TenPortsType        = "ten"
	PortRecycleDocument = "port_recycle"
)

type PortRecycle struct {
	Id        string `bson:"_id"`
	Type      string `bson:"type"`
	PortStart uint   `bson:"port_start"`
}
