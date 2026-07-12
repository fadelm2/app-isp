package entity

type RadCheck struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string `gorm:"column:username"`
	Attribute string `gorm:"column:attribute"`
	Op        string `gorm:"column:op"`
	Value     string `gorm:"column:value"`
}

func (r *RadCheck) TableName() string {
	return "radcheck"
}

type RadReply struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string `gorm:"column:username"`
	Attribute string `gorm:"column:attribute"`
	Op        string `gorm:"column:op"`
	Value     string `gorm:"column:value"`
}

func (r *RadReply) TableName() string {
	return "radreply"
}

type RadAcct struct {
	RadAcctId           int64   `gorm:"column:radacctid;primaryKey;autoIncrement"`
	AcctSessionId       string  `gorm:"column:acctsessionid"`
	AcctUniqueId        string  `gorm:"column:acctuniqueid"`
	Username            string  `gorm:"column:username"`
	GroupName           string  `gorm:"column:groupname"`
	Realm               *string `gorm:"column:realm"`
	NasIpAddress        string  `gorm:"column:nasipaddress"`
	NasPortId           *string `gorm:"column:nasportid"`
	NasPortType         *string `gorm:"column:nasporttype"`
	AcctStartTime       *string `gorm:"column:acctstarttime"`
	AcctUpdateTime      *string `gorm:"column:acctupdatetime"`
	AcctStopTime        *string `gorm:"column:acctstoptime"`
	AcctInterval        *int    `gorm:"column:acctinterval"`
	AcctSessionTime     *uint   `gorm:"column:acctsessiontime"`
	AcctAuthentic       *string `gorm:"column:acctauthentic"`
	ConnectInfoStart    *string `gorm:"column:connectinfo_start"`
	ConnectInfoStop     *string `gorm:"column:connectinfo_stop"`
	AcctInputOctets     *int64  `gorm:"column:acctinputoctets"`
	AcctOutputOctets    *int64  `gorm:"column:acctoutputoctets"`
	CalledStationId     string  `gorm:"column:calledstationid"`
	CallingStationId    string  `gorm:"column:callingstationid"`
	AcctTerminateCause  string  `gorm:"column:acctterminatecause"`
	FramedIpAddress     string  `gorm:"column:framedipaddress"`
	AcctStartDelay      *int    `gorm:"column:acctstartdelay"`
	AcctStopDelay       *int    `gorm:"column:acctstopdelay"`
	XAscendSessionSvrKey *string `gorm:"column:xascendsessionsvrkey"`
}

func (r *RadAcct) TableName() string {
	return "radacct"
}
