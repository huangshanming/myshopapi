package model

const (
	TaskPeriodDaily = "daily"
	TaskPeriodOnce  = "once"

	TaskStatusOngoing   = "ongoing"
	TaskStatusClaimable = "claimable"
	TaskStatusClaimed   = "claimed"

	PointChangeTaskClaim      = "task_claim"
	PointChangeAdminAdjust    = "admin_adjust"
	PointChangeMallExchange   = "points_mall_exchange"
	PointChangeMallRefund     = "points_mall_refund"
	PointChangeLotteryCost    = "lottery_cost"
	PointChangeLotteryRefund  = "lottery_cost_refund"
	PointChangeLotteryReward  = "lottery_reward"

	TaskBizDateOnce = "1970-01-01"
)
