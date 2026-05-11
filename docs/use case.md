# Use Case Descriptions

## User Management

### UC-01 — Sign Up
Users create a new account by providing email, username, and password, verified via OTP.

### UC-02 — Sign In
Authenticated access using credentials or linked Google OAuth account.

### UC-03 — Forgot Password
Users reset their password through email verification (OTP).

### UC-04 — Manage Profile
Users update personal information, including bio, avatar, and cover image.

### UC-05 — Change Settings
Users configure account preferences (privacy, notifications, theme).

### UC-06 — View Activity History
Users review their past posts, comments, and interactions.

### UC-07 — Follow / Unfollow User
Users subscribe to another user’s activity feed.

---

## Content Browsing & Discovery

### UC-08 — Browse Content
Users explore communities and view public posts without necessarily logging in.

### UC-09 — Search Content
Users search for specific keywords, posts, users, or communities.

---

## Post Management

### UC-10 — Create Post
Members publish new content (Text, Image, Video) within a selected community.

### UC-11 — Edit Post
Members modify existing posts (stores edit history).

### UC-12 — Delete Post
Members or Moderators remove posts from public view (soft delete).

### UC-19 — Manage Drafts
Users save unfinished posts as drafts and resume editing later.

---

## Engagement

### UC-13 — Post Comment
Users add a response directly to a post.

### UC-14 — Reply to Comment
Users respond to an existing comment (nested threads supported).

### UC-15 — Vote on Content
Users express approval (Upvote) or disapproval (Downvote) on posts/comments.

### UC-16 — Participate in Poll
Users cast votes in community-created polls.

### UC-17 — Save Post
Users bookmark posts for quick access later.

### UC-18 — Hide Post
Users hide specific posts from their personal feed.

---

## Communication

### UC-20 — Private Messaging
Users send and receive instant 1-1 messages via WebSocket.

### UC-21 — Manage Notifications
Users view and clear real-time alerts for engagement and system updates.

---

## Moderation & Reporting

### UC-22 — Report Violation
Users flag inappropriate content or behavior for moderator review.

### UC-24 — Approve / Reject Post
Moderators review pending posts before they appear in the community.

### UC-25 — Ban / Mute User
Moderators restrict a user’s ability to participate within a community.

---

## Community Management

### UC-23 — Create Community
Users establish a new niche forum and become the “Owner.”

### UC-26 — Configure Community
Owners adjust rules, privacy settings, and appearance of their community.

### UC-27 — Assign Roles
Owners appoint other members as Moderators or Contributors.

---

## System Administration

### UC-28 — Handle Global Reports
System Admins review and act on platform-wide violation reports.

### UC-29 — Manage Platform
System Admins ban/unban users or communities globally.

### UC-30 — View Analytics
System Admins access platform usage and performance metrics.