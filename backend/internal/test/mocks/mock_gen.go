// Package mocks internal/test/mocks/mock_gen.go
package mocks

//go:generate mockgen -source=../../platform/bus/bus.go -destination=../../platform/bus/mock_bus.go -package=bus

// --------------------
// Repo Mocks
// --------------------

//go:generate mockgen -source=../../repo/user_repo.go -destination=../../repo/mocks/mock_user_repo.go -package=mocks
//go:generate mockgen -source=../../repo/community_repo.go -destination=../../repo/mocks/mock_community_repo.go -package=mocks
//go:generate mockgen -source=../../repo/membership_repo.go -destination=../../repo/mocks/mock_membership_repo.go -package=mocks
//go:generate mockgen -source=../../repo/post_repo.go -destination=../../repo/mocks/mock_post_repo.go -package=mocks
//go:generate mockgen -source=../../repo/poll_vote_repo.go -destination=../../repo/mocks/mock_poll_vote_repo.go -package=mocks
//go:generate mockgen -source=../../repo/post_vote_repo.go -destination=../../repo/mocks/mock_post_vote_repo.go -package=mocks
//go:generate mockgen -source=../../repo/comment_repo.go -destination=../../repo/mocks/mock_comment_repo.go -package=mocks
//go:generate mockgen -source=../../repo/notification_repo.go -destination=../../repo/mocks/mock_notification_repo.go -package=mocks
//go:generate mockgen -source=../../repo/channel_repo.go -destination=../../repo/mocks/mock_channel_repo.go -package=mocks
//go:generate mockgen -source=../../repo/message_repo.go -destination=../../repo/mocks/mock_message_repo.go -package=mocks
//go:generate mockgen -source=../../repo/post_history_repo.go -destination=../../repo/mocks/mock_post_history_repo.go -package=mocks
//go:generate mockgen -source=../../repo/email_verification_repo.go -destination=../../repo/mocks/mock_email_verification_repo.go -package=mocks

// --------------------
// Service Mocks
// --------------------

//go:generate mockgen -source=../../service/auth_service.go -destination=../../service/mocks/mock_auth_service.go -package=mocks
//go:generate mockgen -source=../../service/user_service.go -destination=../../service/mocks/mock_user_service.go -package=mocks
//go:generate mockgen -source=../../service/community_service.go -destination=../../service/mocks/mock_community_service.go -package=mocks
//go:generate mockgen -source=../../service/membership_service.go -destination=../../service/mocks/mock_membership_service.go -package=mocks
//go:generate mockgen -source=../../service/post_service.go -destination=../../service/mocks/mock_post_service.go -package=mocks
//go:generate mockgen -source=../../service/comment_service.go -destination=../../service/mocks/mock_comment_service.go -package=mocks
//go:generate mockgen -source=../../service/reputation_service.go -destination=../../service/mocks/mock_reputation_service.go -package=mocks
//go:generate mockgen -source=../../service/notification_service.go -destination=../../service/mocks/mock_notification_service.go -package=mocks
//go:generate mockgen -source=../../service/channel_service.go -destination=../../service/mocks/mock_channel_service.go -package=mocks
//go:generate mockgen -source=../../service/message_service.go -destination=../../service/mocks/mock_message_service.go -package=mocks
//go:generate mockgen -source=../../service/post_history_service.go -destination=../../service/mocks/mock_post_history_service.go -package=mocks
