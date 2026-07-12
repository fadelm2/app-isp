package model

type DashboardStatsResponse struct {
	TotalCustomers     int64   `json:"total_customers"`
	ActiveCustomers    int64   `json:"active_customers"`
	SuspendedCustomers int64   `json:"suspended_customers"`
	OwedCustomers      int64   `json:"owed_customers"`
	TodayPayments      float64 `json:"today_payments"`
	MonthlyRevenue     float64 `json:"monthly_revenue"`
	RouterStatus       string  `json:"router_status"` // online, offline
	OnlineUsers        int     `json:"online_users"`
	OfflineUsers       int     `json:"offline_users"`
}
