package bus

// --- Post Events ---

type PostCreatedEvent struct {
	AuthorID string
}

func (e PostCreatedEvent) Topic() string {
	return "post.created"
}
func (e PostCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID}
}

type PostUpvotedEvent struct {
	AuthorID string
}

func (e PostUpvotedEvent) Topic() string {
	return "post.upvoted"
}
func (e PostUpvotedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID}
}

type PostDownvotedEvent struct {
	AuthorID string
	VoterID  string
}

func (e PostDownvotedEvent) Topic() string {
	return "post.downvoted"
}
func (e PostDownvotedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID, "voterId": e.VoterID}
}

// --- Comment Events ---

type CommentCreatedEvent struct {
	AuthorID string
}

func (e CommentCreatedEvent) Topic() string {
	return "comment.created"
}
func (e CommentCreatedEvent) Payload() map[string]interface{} {
	return map[string]interface{}{"authorId": e.AuthorID}
}
