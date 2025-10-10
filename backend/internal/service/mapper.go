package service

import (
	"errors"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mapCreateRequestToPostModel chuyển đổi DTO request thành model để lưu vào DB.
func mapCreateRequestToPostModel(req *dto.CreatePostRequest, authorID primitive.ObjectID) (*model.Post, error) {
	communityID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		AuthorID:    authorID,
		CommunityID: communityID,
		Type:        model.PostType(req.Type),
		Title:       req.Title,
		Content:     &model.PostContent{},
	}

	switch post.Type {
	case model.PostTypeText:
		post.Content.Text = req.Text
	case model.PostTypeImage:
		images := make([]model.Image, len(req.Images))
		for i, imgReq := range req.Images {
			images[i] = model.Image{URL: imgReq.URL}
		}
		post.Content.Images = images
	case model.PostTypePoll:
		if req.Poll == nil {
			return nil, errors.New("poll data is required for poll type post")
		}
		options := make([]model.PollOption, len(req.Poll.Options))
		for i, optText := range req.Poll.Options {
			options[i] = model.PollOption{ID: primitive.NewObjectID(), Text: optText, Votes: 0}
		}
		totalVotes := 0
		post.Content.Poll = &model.Poll{
			Question:      req.Poll.Question,
			Options:       options,
			TotalVotes:    &totalVotes,
			ExpiresAt:     req.Poll.ExpiresAt,
			AllowMultiple: req.Poll.AllowMultiple,
		}
	case model.PostTypeVideo:
		// Thêm bước kiểm tra để đảm bảo an toàn
		if req.Video == nil {
			return nil, errors.New("video data is required for video type post")
		}
		post.Content.Video = &model.Video{
			URL:       req.Video.URL,
			Thumbnail: req.Video.Thumbnail,
		}

	}
	return post, nil
}

// mapPostModelToResponse chuyển đổi model từ DB thành DTO để trả về cho client.
func mapPostModelToResponse(post *model.Post, userVote *model.Vote, userPollVotes []*model.PollVote) *dto.PostResponse {
	res := &dto.PostResponse{
		ID:             post.ID.Hex(),
		AuthorID:       post.AuthorID.Hex(),
		AuthorUsername: post.AuthorUsername,
		AuthorAvatar:   post.AuthorAvatar,
		CommunityID:    post.CommunityID.Hex(),
		CommunityName:  post.CommunityName,
		Title:          post.Title,
		Type:           string(post.Type),
		VotesCount:     mapVotesToResponse(post.VotesCount),
		CommentsCount:  post.CommentsCount,
		CreatedAt:      post.CreatedAt,
		UpdatedAt:      post.UpdatedAt,
	}

	if post.Content != nil {
		res.Content = &dto.PostContentResponse{}
		if post.Content.Text != "" {
			res.Content.Text = post.Content.Text
		}
		if len(post.Content.Images) > 0 {
			res.Content.Images = make([]dto.ImageResponse, len(post.Content.Images))
			for i, img := range post.Content.Images {
				res.Content.Images[i] = dto.ImageResponse{ID: img.ID.Hex(), URL: img.URL}
			}
		}
		if post.Content.Poll != nil {
			res.Content.Poll = mapPollToResponse(post.Content.Poll, userPollVotes)
		}
		if post.Content.Video != nil {
			res.Content.Video = &dto.VideoResponse{
				URL:       post.Content.Video.URL,
				Thumbnail: post.Content.Video.Thumbnail,
				// Giả sử VideoResponse có ID, nếu không thì bỏ qua
				// ID: post.Content.Video.ID.Hex(),
			}
		}
	}

	if userVote != nil {
		if userVote.Value {
			res.UserVote = "up"
		} else {
			res.UserVote = "down"
		}
	}
	return res
}

// mapVotesToResponse chuyển đổi VotesCount model thành response.
func mapVotesToResponse(votes *model.VotesCount) *dto.VotesCountResponse {
	if votes == nil {
		return &dto.VotesCountResponse{Up: 0, Down: 0, Score: 0}
	}
	return &dto.VotesCountResponse{
		Up:    votes.Up,
		Down:  votes.Down,
		Score: votes.Up - votes.Down,
	}
}

// mapPollToResponse chuyển đổi Poll model thành response.
func mapPollToResponse(poll *model.Poll, userVotes []*model.PollVote) *dto.PollResponse {
	if poll == nil {
		return nil
	}

	totalVotes := 0
	if poll.TotalVotes != nil {
		totalVotes = *poll.TotalVotes
	}

	// Tạo một map để tra cứu vote của user nhanh hơn
	userVoteMap := make(map[primitive.ObjectID]bool)
	userVoteIDs := make([]string, len(userVotes))
	for i, v := range userVotes {
		userVoteMap[v.OptionID] = true
		userVoteIDs[i] = v.OptionID.Hex()
	}

	optionsRes := make([]dto.PollOptionResponse, len(poll.Options))
	for i, opt := range poll.Options {
		percentage := 0.0
		if totalVotes > 0 {
			percentage = (float64(opt.Votes) / float64(totalVotes)) * 100
		}
		optionsRes[i] = dto.PollOptionResponse{
			ID:         opt.ID.Hex(),
			Text:       opt.Text,
			Votes:      opt.Votes,
			Percentage: percentage,
		}
	}

	return &dto.PollResponse{
		Question:      poll.Question,
		Options:       optionsRes,
		TotalVotes:    totalVotes,
		UserVoteIDs:   userVoteIDs,
		ExpiresAt:     poll.ExpiresAt,
		AllowMultiple: poll.AllowMultiple,
	}
}
