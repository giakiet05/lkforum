# TÀI LIỆU THIẾT KẾ CƠ SỞ DỮ LIỆU - LKFORUM

## 1. BẢNG USERS (Người dùng)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                    |
| --- | -------------- | ------------ | ------------------------------------------ |
| 1   | id             | ObjectID     | ID duy nhất của người dùng (Primary Key)   |
| 2   | username       | String       | Tên người dùng (unique)                    |
| 3   | email          | String       | Email người dùng (unique)                  |
| 4   | reputation     | Integer      | Điểm uy tín của người dùng                 |
| 5   | password       | String       | Mật khẩu đã mã hóa (chỉ local auth)        |
| 6   | provider       | String       | Phương thức đăng ký: "local" hoặc "google" |
| 7   | provider_id    | String       | ID từ provider (Google OAuth)              |
| 8   | role           | String       | Vai trò: "user" hoặc "admin"               |
| 9   | role_content   | Object       | Thông tin chi tiết theo vai trò            |
| 10  | settings       | Object       | Cài đặt tùy chọn của người dùng            |
| 11  | is_verified    | Boolean      | Email đã được xác thực chưa                |
| 12  | is_banned      | Boolean      | Tài khoản có bị cấm không                  |
| 13  | ban_until      | DateTime     | Thời gian hết hạn cấm (null = vĩnh viễn)   |
| 14  | ban_reason     | String       | Lý do bị cấm                               |
| 15  | created_at     | DateTime     | Thời gian tạo tài khoản                    |
| 16  | deleted_at     | DateTime     | Thời gian xóa tài khoản (soft delete)      |

### 1.1. Role Content - User (role_content.as_user)

| STT | Tên thuộc tính | Kiểu dữ liệu  | Ý nghĩa                                          |
| --- | -------------- | ------------- | ------------------------------------------------ |
| 1   | avatar         | Object        | Ảnh đại diện {url, public_id, uploaded_at}       |
| 2   | cover          | Object        | Ảnh bìa {url, public_id, uploaded_at}            |
| 3   | bio            | String        | Tiểu sử người dùng                               |
| 4   | gender         | String        | Giới tính: "male", "female", "prefer_not_to_say" |
| 5   | date_of_birth  | DateTime      | Ngày sinh                                        |
| 6   | location       | String        | Tỉnh/thành phố (VN provinces)                    |
| 7   | interests      | Array[String] | Sở thích (lập trình, thiết kế, gaming...)        |
| 8   | social_links   | Object        | Liên kết mạng xã hội                             |
| 9   | stats          | Object        | Thống kê hoạt động                               |

### 1.2. Social Links (role_content.as_user.social_links)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa         |
| --- | -------------- | ------------ | --------------- |
| 1   | website        | String       | Website cá nhân |
| 2   | facebook       | String       | Link Facebook   |
| 3   | youtube        | String       | Link YouTube    |
| 4   | github         | String       | Link GitHub     |

### 1.3. Activity Stats (role_content.as_user.stats)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa               |
| --- | -------------- | ------------ | --------------------- |
| 1   | post_count     | Integer      | Số bài đăng           |
| 2   | comment_count  | Integer      | Số bình luận          |
| 3   | total_upvotes  | Integer      | Tổng upvote nhận được |
| 4   | joined_at      | DateTime     | Thời gian tham gia    |
| 5   | last_active_at | DateTime     | Lần hoạt động cuối    |

### 1.4. User Settings (settings)

| STT | Tên thuộc tính                  | Kiểu dữ liệu | Ý nghĩa                                    |
| --- | ------------------------------- | ------------ | ------------------------------------------ |
| 1   | appearance.theme                | String       | Chủ đề: "light", "dark", "auto"            |
| 2   | appearance.font_size            | String       | Kích thước chữ: "small", "medium", "large" |
| 3   | notifications.in_app_enabled    | Boolean      | Bật thông báo trong app                    |
| 4   | notifications.email_enabled     | Boolean      | Bật thông báo email                        |
| 5   | notifications.notify_on_comment | Boolean      | Thông báo khi có comment                   |
| 6   | notifications.notify_on_mention | Boolean      | Thông báo khi được nhắc đến                |
| 7   | notifications.notify_on_upvote  | Boolean      | Thông báo khi được upvote                  |
| 8   | notifications.notify_on_message | Boolean      | Thông báo khi có tin nhắn                  |
| 9   | privacy.show_profile            | Boolean      | Hiển thị profile công khai                 |
| 10  | privacy.show_email              | Boolean      | Hiển thị email trên profile                |
| 11  | privacy.show_post_history       | Boolean      | Hiển thị lịch sử bài đăng                  |
| 12  | privacy.allow_direct_messages   | Boolean      | Cho phép nhận tin nhắn riêng               |
| 13  | privacy.allow_mentions          | Boolean      | Cho phép được nhắc đến                     |
| 14  | content.allow_nsfw              | Boolean      | Cho phép xem nội dung NSFW                 |

---

## 2. BẢNG POSTS (Bài đăng)

| STT | Tên thuộc tính    | Kiểu dữ liệu  | Ý nghĩa                                                             |
| --- | ----------------- | ------------- | ------------------------------------------------------------------- |
| 1   | \_id              | ObjectID      | ID duy nhất của bài đăng (Primary Key)                              |
| 2   | author_id         | ObjectID      | ID tác giả (Foreign Key → users)                                    |
| 3   | community_id      | ObjectID      | ID cộng đồng (Foreign Key → communities)                            |
| 4   | type              | String        | Loại bài đăng: "text", "poll", "video", "image"                     |
| 5   | title             | String        | Tiêu đề bài đăng                                                    |
| 6   | content           | Object        | Nội dung bài đăng                                                   |
| 7   | votes_count       | Object        | Số lượng upvote/downvote {up, down}                                 |
| 8   | comments_count    | Integer       | Số lượng bình luận                                                  |
| 9   | created_at        | DateTime      | Thời gian tạo                                                       |
| 10  | updated_at        | DateTime      | Thời gian cập nhật                                                  |
| 11  | is_deleted        | Boolean       | Bài đăng đã bị xóa chưa                                             |
| 12  | is_hidden         | Boolean       | Bài đăng bị ẩn bởi người dùng                                       |
| 13  | is_edited         | Boolean       | Bài đăng đã được chỉnh sửa                                          |
| 14  | tags              | Array[String] | Các tag của bài đăng                                                |
| 15  | is_draft          | Boolean       | Bài đăng là nháp                                                    |
| 16  | is_ban            | Boolean       | Bài đăng bị ban bởi admin/mod                                       |
| 17  | ban_reason        | String        | Lý do bị ban                                                        |
| 18  | moderation_status | String        | Trạng thái kiểm duyệt: "pending", "approved", "rejected", "skipped" |
| 19  | moderation_result | Object        | Kết quả kiểm duyệt AI                                               |
| 20  | moderated_at      | DateTime      | Thời gian kiểm duyệt                                                |

### 2.1. Post Content (content)

| STT | Tên thuộc tính | Kiểu dữ liệu  | Ý nghĩa                                       |
| --- | -------------- | ------------- | --------------------------------------------- |
| 1   | text           | String        | Nội dung văn bản                              |
| 2   | images         | Array[Object] | Danh sách ảnh {url, public_id, uploaded_at}   |
| 3   | videos         | Array[Object] | Danh sách video {url, public_id, uploaded_at} |
| 4   | poll           | Object        | Thông tin poll nếu type = "poll"              |

### 2.2. Poll Content (content.poll)

| STT | Tên thuộc tính | Kiểu dữ liệu  | Ý nghĩa                        |
| --- | -------------- | ------------- | ------------------------------ |
| 1   | question       | String        | Câu hỏi poll                   |
| 2   | options        | Array[Object] | Các lựa chọn {id, text, votes} |
| 3   | total_votes    | Integer       | Tổng số lượt vote              |
| 4   | expires_at     | DateTime      | Thời gian hết hạn poll         |
| 5   | allow_multiple | Boolean       | Cho phép chọn nhiều đáp án     |

### 2.3. Moderation Result (moderation_result)

| STT | Tên thuộc tính | Kiểu dữ liệu  | Ý nghĩa                                      |
| --- | -------------- | ------------- | -------------------------------------------- |
| 1   | is_violation   | Boolean       | Có vi phạm không                             |
| 2   | confidence     | Float         | Độ tin cậy (0.0 - 1.0)                       |
| 3   | categories     | Array[String] | Loại vi phạm: ["hate_speech", "violence"...] |
| 4   | reason         | String        | Lý do vi phạm (tiếng Việt)                   |
| 5   | checked_text   | Boolean       | Đã kiểm tra văn bản chưa                     |
| 6   | checked_media  | Boolean       | Đã kiểm tra media chưa                       |

---

## 3. BẢNG COMMENTS (Bình luận)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                  |
| --- | -------------- | ------------ | ---------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất của comment (Primary Key)    |
| 2   | author         | Object       | Thông tin tác giả {id, username, avatar} |
| 3   | post_id        | ObjectID     | ID bài đăng (Foreign Key → posts)        |
| 4   | parent_id      | ObjectID     | ID comment cha (cho nested comments)     |
| 5   | content        | String       | Nội dung bình luận                       |
| 6   | votes_count    | Object       | Số lượng upvote/downvote {up, down}      |
| 7   | created_at     | DateTime     | Thời gian tạo                            |
| 8   | deleted_at     | DateTime     | Thời gian xóa (soft delete)              |
| 9   | is_deleted     | Boolean      | Comment đã bị xóa chưa                   |

---

## 4. BẢNG COMMUNITIES (Cộng đồng)

| STT | Tên thuộc tính   | Kiểu dữ liệu  | Ý nghĩa                                          |
| --- | ---------------- | ------------- | ------------------------------------------------ |
| 1   | \_id             | ObjectID      | ID duy nhất của community (Primary Key)          |
| 2   | name             | String        | Tên community (unique, dùng trong URL /lk/:name) |
| 3   | description      | String        | Mô tả community                                  |
| 4   | avatar           | String        | URL avatar community                             |
| 5   | banner           | String        | URL banner community                             |
| 6   | setting          | Object        | Cài đặt community                                |
| 7   | rules            | Array[Object] | Quy tắc community                                |
| 8   | moderators       | Array[Object] | Danh sách moderators                             |
| 9   | member_count     | Integer       | Số lượng thành viên                              |
| 10  | post_count       | Integer       | Số lượng bài đăng                                |
| 11  | create_at        | DateTime      | Thời gian tạo                                    |
| 12  | create_by_id     | ObjectID      | ID người tạo (Foreign Key → users)               |
| 13  | create_by_name   | String        | Tên người tạo                                    |
| 14  | create_by_avatar | String        | Avatar người tạo                                 |
| 15  | is_18_plus       | Boolean       | Community 18+                                    |
| 16  | is_deleted       | Boolean       | Community đã bị xóa                              |
| 17  | is_banned        | Boolean       | Community bị ban                                 |
| 18  | ban_reason       | String        | Lý do bị ban                                     |

### 4.1. Community Settings (setting)

| STT | Tên thuộc tính        | Kiểu dữ liệu | Ý nghĩa                              |
| --- | --------------------- | ------------ | ------------------------------------ |
| 1   | is_private            | Boolean      | Bài đăng chỉ thành viên mới xem được |
| 2   | post_require_approval | Boolean      | Bài đăng cần mod duyệt               |
| 3   | join_require_approval | Boolean      | Tham gia cần mod duyệt               |
| 4   | max_post_length       | Integer      | Độ dài tối đa bài đăng               |

### 4.2. Community Rules (rules)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa         |
| --- | -------------- | ------------ | --------------- |
| 1   | title          | String       | Tiêu đề quy tắc |
| 2   | description    | String       | Mô tả quy tắc   |

### 4.3. Moderators (moderators)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                            |
| --- | -------------- | ------------ | ---------------------------------- |
| 1   | user_id        | ObjectID     | ID moderator (Foreign Key → users) |
| 2   | username       | String       | Tên moderator                      |
| 3   | avatar         | Object       | Avatar moderator                   |
| 4   | is_active      | Boolean      | Moderator đang hoạt động           |
| 5   | assigned_at    | DateTime     | Thời gian được bổ nhiệm            |

---

## 5. BẢNG MEMBERSHIPS (Thành viên cộng đồng)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                  |
| --- | -------------- | ------------ | ---------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)                |
| 2   | user_id        | ObjectID     | ID người dùng (Foreign Key → users)      |
| 3   | community_id   | ObjectID     | ID cộng đồng (Foreign Key → communities) |

**Index:** Unique compound index trên (user_id, community_id)

---

## 6. BẢNG VOTES (Bình chọn)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                               |
| --- | -------------- | ------------ | ------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)             |
| 2   | user_id        | ObjectID     | ID người vote (Foreign Key → users)   |
| 3   | target_type    | String       | Loại đối tượng: "post" hoặc "comment" |
| 4   | target_id      | ObjectID     | ID đối tượng (post hoặc comment)      |
| 5   | value          | Boolean      | true = upvote, false = downvote       |
| 6   | create_at      | DateTime     | Thời gian vote                        |

**Index:** Unique compound index trên (user_id, target_id, target_type)

---

## 7. BẢNG POLL_VOTES (Vote trong Poll)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                |
| --- | -------------- | ------------ | -------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)              |
| 2   | post_id        | ObjectID     | ID bài đăng poll (Foreign Key → posts) |
| 3   | user_id        | ObjectID     | ID người vote (Foreign Key → users)    |
| 4   | option_id      | String       | ID option được chọn                    |
| 5   | created_at     | DateTime     | Thời gian vote                         |

---

## 8. BẢNG SAVED_POSTS (Bài đăng đã lưu)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                            |
| --- | -------------- | ------------ | ---------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)          |
| 2   | user_id        | ObjectID     | ID người lưu (Foreign Key → users) |
| 3   | post_id        | ObjectID     | ID bài đăng (Foreign Key → posts)  |
| 4   | saved_at       | DateTime     | Thời gian lưu                      |

**Index:** Unique compound index trên (user_id, post_id)

---

## 9. BẢNG LIKED_POSTS (Bài đăng đã thích)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                              |
| --- | -------------- | ------------ | ------------------------------------ |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)            |
| 2   | user_id        | ObjectID     | ID người thích (Foreign Key → users) |
| 3   | post_id        | ObjectID     | ID bài đăng (Foreign Key → posts)    |
| 4   | liked_at       | DateTime     | Thời gian thích                      |

**Index:** Unique compound index trên (user_id, post_id)

---

## 10. BẢNG POST_HISTORY (Lịch sử xem bài đăng)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                            |
| --- | -------------- | ------------ | ---------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)          |
| 2   | user_id        | ObjectID     | ID người xem (Foreign Key → users) |
| 3   | post_id        | ObjectID     | ID bài đăng (Foreign Key → posts)  |
| 4   | viewed_at      | DateTime     | Thời gian xem                      |

---

## 11. BẢNG REPORTS (Báo cáo vi phạm)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                   |
| --- | -------------- | ------------ | ----------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)                 |
| 2   | reporter_id    | ObjectID     | ID người báo cáo (Foreign Key → users)    |
| 3   | target_id      | ObjectID     | ID đối tượng bị báo cáo                   |
| 4   | target_type    | String       | Loại đối tượng: "user", "post", "comment" |
| 5   | reason         | String       | Lý do báo cáo                             |
| 6   | description    | String       | Mô tả chi tiết                            |
| 7   | is_deleted     | Boolean      | Report đã bị xóa (đã xử lý)               |
| 8   | deleted_at     | DateTime     | Thời gian xóa                             |
| 9   | created_at     | DateTime     | Thời gian báo cáo                         |

---

## 12. BẢNG NOTIFICATIONS (Thông báo)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                                               |
| --- | -------------- | ------------ | --------------------------------------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)                                             |
| 2   | recipient_id   | ObjectID     | ID người nhận (Foreign Key → users)                                   |
| 3   | actor_id       | ObjectID     | ID người thực hiện hành động                                          |
| 4   | type           | String       | Loại: "comment", "like", "follow", "mention", "new_message", "system" |
| 5   | message        | String       | Nội dung thông báo                                                    |
| 6   | link           | String       | Link đến nội dung liên quan                                           |
| 7   | is_read        | Boolean      | Đã đọc chưa                                                           |
| 8   | metadata       | Object       | Dữ liệu bổ sung (JSON)                                                |
| 9   | created_at     | DateTime     | Thời gian tạo thông báo                                               |

---

## 13. BẢNG MESSAGES (Tin nhắn)

| STT | Tên thuộc tính  | Kiểu dữ liệu | Ý nghĩa                                                        |
| --- | --------------- | ------------ | -------------------------------------------------------------- |
| 1   | \_id            | ObjectID     | ID duy nhất (Primary Key)                                      |
| 2   | channel_id      | ObjectID     | ID kênh chat                                                   |
| 3   | sender_id       | ObjectID     | ID người gửi (Foreign Key → users), null nếu là system message |
| 4   | sender_username | String       | Tên người gửi                                                  |
| 5   | type            | String       | Loại tin nhắn: "user" hoặc "system"                            |
| 6   | content         | String       | Nội dung tin nhắn                                              |
| 7   | is_read         | Boolean      | Đã đọc chưa                                                    |
| 8   | is_send         | Boolean      | Đã gửi thành công chưa                                         |
| 9   | created_at      | DateTime     | Thời gian gửi                                                  |
| 10  | is_deleted      | Boolean      | Tin nhắn đã bị xóa                                             |
| 11  | deleted_at      | DateTime     | Thời gian xóa                                                  |

---

## 14. BẢNG COMMUNITY_BANS (Cấm trong cộng đồng)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                                            |
| --- | -------------- | ------------ | ------------------------------------------------------------------ |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)                                          |
| 2   | user_id        | ObjectID     | ID người bị cấm (Foreign Key → users)                              |
| 3   | community_id   | ObjectID     | ID cộng đồng (Foreign Key → communities)                           |
| 4   | type           | String       | Loại cấm: "banned" (không được vào) hoặc "muted" (không được post) |
| 5   | reason         | String       | Lý do cấm                                                          |
| 6   | banned_by      | ObjectID     | ID người thực hiện cấm (Foreign Key → users)                       |
| 7   | banned_at      | DateTime     | Thời gian bắt đầu cấm                                              |
| 8   | expires_at     | DateTime     | Thời gian hết hạn cấm                                              |
| 9   | is_deleted     | Boolean      | Đã gỡ cấm chưa                                                     |
| 10  | deleted_at     | DateTime     | Thời gian gỡ cấm                                                   |

---

## 15. BẢNG EMAIL_VERIFICATIONS (Xác thực email)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                 |
| --- | -------------- | ------------ | --------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)               |
| 2   | email          | String       | Email cần xác thực                      |
| 3   | otp            | String       | Mã OTP (6 số, đã mã hóa)                |
| 4   | otp_expires_at | DateTime     | Thời gian hết hạn OTP                   |
| 5   | is_verified    | Boolean      | Đã xác thực chưa                        |
| 6   | nonce          | String       | Token ngẫu nhiên để tránh replay attack |
| 7   | created_at     | DateTime     | Thời gian tạo                           |

**Mục đích:** Lưu tạm OTP trước khi user hoàn tất đăng ký

---

## 16. BẢNG PASSWORD_RESETS (Đặt lại mật khẩu)

| STT | Tên thuộc tính | Kiểu dữ liệu | Ý nghĩa                                 |
| --- | -------------- | ------------ | --------------------------------------- |
| 1   | \_id           | ObjectID     | ID duy nhất (Primary Key)               |
| 2   | email          | String       | Email yêu cầu reset                     |
| 3   | otp            | String       | Mã OTP (6 số, đã mã hóa)                |
| 4   | otp_expires_at | DateTime     | Thời gian hết hạn OTP                   |
| 5   | is_verified    | Boolean      | OTP đã được xác thực chưa               |
| 6   | nonce          | String       | Token ngẫu nhiên để tránh replay attack |
| 7   | created_at     | DateTime     | Thời gian tạo                           |

**Mục đích:** Lưu tạm OTP khi user quên mật khẩu

---

## 17. BẢNG DRAFTS (Bài nháp)

| STT | Tên thuộc tính | Kiểu dữ liệu  | Ý nghĩa                                          |
| --- | -------------- | ------------- | ------------------------------------------------ |
| 1   | \_id           | ObjectID      | ID duy nhất (Primary Key)                        |
| 2   | author_id      | ObjectID      | ID tác giả (Foreign Key → users)                 |
| 3   | community_id   | String        | ID cộng đồng (optional)                          |
| 4   | type           | String        | Loại bài đăng: "text", "poll", "video", "image"  |
| 5   | title          | String        | Tiêu đề bài nháp                                 |
| 6   | content        | Object        | Nội dung bài nháp (cấu trúc giống posts.content) |
| 7   | tags           | Array[String] | Các tag                                          |
| 8   | created_at     | DateTime      | Thời gian tạo nháp                               |
| 9   | updated_at     | DateTime      | Thời gian cập nhật nháp                          |

---

## QUAN HỆ GIỮA CÁC BẢNG

### Quan hệ 1-N (One-to-Many)

- **users** → **posts** (1 user có nhiều posts)
- **users** → **comments** (1 user có nhiều comments)
- **users** → **communities** (1 user tạo nhiều communities)
- **posts** → **comments** (1 post có nhiều comments)
- **communities** → **posts** (1 community có nhiều posts)

### Quan hệ N-N (Many-to-Many)

- **users** ↔ **communities** (qua bảng **memberships**)
- **users** ↔ **posts** (saved: qua bảng **saved_posts**)
- **users** ↔ **posts** (liked: qua bảng **liked_posts**)
- **users** ↔ **posts/comments** (voted: qua bảng **votes**)

### Quan hệ Self-Referencing

- **comments** → **comments** (parent_id cho nested comments)

---

## INDEX QUAN TRỌNG

### Users

- `email` (unique)
- `username` (unique)
- `provider_id` (cho OAuth lookup)

### Posts

- `author_id` (tìm posts của user)
- `community_id` (tìm posts trong community)
- `created_at` (sắp xếp theo thời gian)
- `type` (filter theo loại post)

### Comments

- `post_id` (lấy comments của post)
- `parent_id` (lấy reply comments)

### Communities

- `name` (unique, dùng trong URL)
- `create_by_id` (tìm communities của user)

### Memberships

- `(user_id, community_id)` (compound unique)

### Votes

- `(user_id, target_id, target_type)` (compound unique)

### Reports

- `target_id, target_type` (tìm reports của đối tượng)
- `reporter_id` (tìm reports của user)

### Notifications

- `recipient_id, is_read` (lấy unread notifications)
- `created_at` (sắp xếp)

---

## NOTES

1. **ObjectID**: Tất cả ID đều dùng MongoDB ObjectID (24 hex characters)
2. **Soft Delete**: Nhiều bảng dùng `is_deleted` + `deleted_at` thay vì xóa thật
3. **Timestamps**: Hầu hết bảng có `created_at`, một số có `updated_at`
4. **Indexes**: Cần tạo indexes cho các trường hay query để tối ưu performance
5. **Foreign Keys**: MongoDB không enforce foreign keys, cần validate ở application layer
6. **Embedded vs Referenced**: Một số data được embed (như comments.author), một số reference (như post.author_id)

---

**Tổng số bảng: 17**
**Database: MongoDB**
**Ngày tạo: 20/01/2026**
