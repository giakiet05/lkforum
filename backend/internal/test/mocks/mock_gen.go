// Package mocks internal/test/mocks/mock_gen.go
package mocks

//go:generate mockgen -source=../../platform/bus/bus.go -destination=../../platform/bus/mock_bus.go -package=bus
//go:generate mockgen -source=../../auth/token_service.go -destination=../../auth/mock_token_service.go -package=auth

// --------------------
// Repo Mocks
// --------------------

//go:generate mockgen -source=../../repo/channel_repo.go -destination=../../repo/mocks/mock_channel_repo.go -package=mocks
//go:generate mockgen -source=../../repo/comment_repo.go -destination=../../repo/mocks/mock_comment_repo.go -package=mocks
//go:generate mockgen -source=../../repo/community_repo.go -destination=../../repo/mocks/mock_community_repo.go -package=mocks
//go:generate mockgen -source=../../repo/draft_repo.go -destination=../../repo/mocks/mock_draft_repo.go -package=mocks
//go:generate mockgen -source=../../repo/email_verification_repo.go -destination=../../repo/mocks/mock_email_verification_repo.go -package=mocks
//go:generate mockgen -source=../../repo/membership_repo.go -destination=../../repo/mocks/mock_membership_repo.go -package=mocks
//go:generate mockgen -source=../../repo/message_repo.go -destination=../../repo/mocks/mock_message_repo.go -package=mocks
//go:generate mockgen -source=../../repo/notification_repo.go -destination=../../repo/mocks/mock_notification_repo.go -package=mocks
//go:generate mockgen -source=../../repo/poll_vote_repo.go -destination=../../repo/mocks/mock_poll_vote_repo.go -package=mocks
//go:generate mockgen -source=../../repo/post_history_repo.go -destination=../../repo/mocks/mock_post_history_repo.go -package=mocks
//go:generate mockgen -source=../../repo/post_repo.go -destination=../../repo/mocks/mock_post_repo.go -package=mocks
//go:generate mockgen -source=../../repo/report_repo.go -destination=../../repo/mocks/mock_report_repo.go -package=mocks
//go:generate mockgen -source=../../repo/saved_post_repo.go -destination=../../repo/mocks/mock_saved_post_repo.go -package=mocks
//go:generate mockgen -source=../../repo/user_repo.go -destination=../../repo/mocks/mock_user_repo.go -package=mocks
//go:generate mockgen -source=../../repo/vote_repo.go -destination=../../repo/mocks/mock_vote_repo.go -package=mocks
