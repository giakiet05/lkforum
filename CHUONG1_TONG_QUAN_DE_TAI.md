# CHƯƠNG 1: TỔNG QUAN ĐỀ TÀI

## 1.1. HIỆN TRẠNG

### 1.1.1. Tình hình các diễn đàn hiện nay

Hiện nay, các nền tảng diễn đàn phổ biến như Reddit, Facebook Groups đang được sử dụng rộng rãi. Tuy nhiên, các nền tảng này còn một số hạn chế:

**Về quản lý cộng đồng:**

- Quyền tự chủ hạn chế trong việc cấu hình cộng đồng
- Hệ thống phân quyền (owner, moderator) chưa linh hoạt
- Khó kiểm soát việc duyệt bài viết và quản lý thành viên

**Về trải nghiệm người dùng:**

- Giao diện phức tạp, khó sử dụng với người mới
- Thiếu tính năng tìm kiếm và lọc nâng cao
- Hệ thống thông báo chưa tối ưu

**Về bảo mật:**

- Lo ngại về quyền riêng tư dữ liệu cá nhân
- Xác thực người dùng chưa đủ mạnh
- Quản lý phiên đăng nhập còn lỏng lẻo

### 1.1.2. Nhu cầu thực tế

Cần có một nền tảng diễn đàn cho phép:

- Người dùng tự tạo và quản lý cộng đồng riêng
- Hệ thống phân quyền rõ ràng (owner, moderator, member)
- Kiểm duyệt nội dung linh hoạt (duyệt bài, duyệt thành viên)
- Xác thực an toàn với email OTP và Google OAuth
- Giao diện đơn giản, responsive trên mọi thiết bị
- Tin nhắn real-time giữa người dùng
- Không chèn quảng cáo
- Phù hợp với cộng đồng nhỏ

## 1.2. NỘI DUNG VÀ PHẠM VI ĐỀ TÀI

### 1.2.1. Nội dung đề tài

Xây dựng hệ thống diễn đàn cộng đồng LKForum với các nội dung chính:

**1. Công nghệ**

- Frontend: Svelte 5, TypeScript, Vite
- Backend: Go (Gin Framework), RESTful API
- Database: MongoDB Atlas
- Cache: Redis
- Authentication: JWT, Google OAuth
- Real-time: WebSocket
- Deployment: Docker

**2. Chức năng chính**

**Module xác thực và người dùng:**

- Đăng ký với email + xác thực OTP (10 phút)
- Đăng nhập: email/password hoặc Google OAuth
- Quên mật khẩu: gửi OTP reset qua email
- Quản lý profile: username, avatar, banner, bio
- Cài đặt tài khoản

**Module cộng đồng:**

- Tạo cộng đồng: tên (3-21 ký tự), mô tả, avatar, banner
- Cấu hình cộng đồng:
  - Private/Public
  - Yêu cầu duyệt bài (post_require_approval)
  - Yêu cầu duyệt thành viên (join_require_approval)
  - Giới hạn độ dài bài viết
  - Nội dung 18+
- Quy tắc cộng đồng (rules)
- Quản lý moderator (tối đa 25 người)
- Ban/Unban user (ban vĩnh viễn hoặc tạm thời)
- Mute user

**Module bài viết:**

- Tạo bài viết: tiêu đề (3-300 ký tự), nội dung text, hình ảnh, video
- Poll: tối đa 10 options, có thể vote
- Chỉnh sửa bài (lưu lịch sử chỉnh sửa)
- Xóa bài
- Lưu nháp (Draft)
- Vote: upvote/downvote
- Lưu bài viết (Save)
- Ẩn bài viết (Hide)
- Báo cáo bài vi phạm
- URL slug: `/post/title-slug-{id}`

**Module bình luận:**

- Bình luận trên bài viết
- Vote comment
- Xóa comment
- Lọc comment (filter)

**Module tin nhắn:**

- Nhắn tin 1-1 real-time (WebSocket)
- Tạo channel giữa 2 người
- Lịch sử tin nhắn

**Module kiểm duyệt (Moderator):**

- Xem bài viết pending (chờ duyệt)
- Xem bài viết edited (đã chỉnh sửa)
- Duyệt/từ chối bài viết
- Ban/Unban bài viết
- Xử lý báo cáo

**Module quản trị (Admin):**

- Quản lý người dùng: ban/unban, xóa, khôi phục
- Quản lý cộng đồng: ban/unban cộng đồng
- Xem báo cáo toàn hệ thống
- Thống kê: tổng quan, người dùng, nội dung

**Module thông báo:**

- Thông báo real-time
- Đánh dấu đã đọc
- Xóa thông báo

**Module tìm kiếm:**

- Tìm kiếm cộng đồng
- Tìm kiếm người dùng
- Tìm kiếm trong topbar

### 1.2.2. Phạm vi đề tài

**Phạm vi người dùng:**

- Người dùng thông thường (Member)
- Người tạo cộng đồng (Owner)
- Người kiểm duyệt (Moderator)
- Quản trị viên (Admin)

**Phạm vi kỹ thuật:**

- Web application (desktop + mobile responsive)
- RESTful API
- WebSocket cho real-time messaging
- MongoDB collections: users, posts, comments, communities, memberships, votes, messages, notifications, reports, drafts, user_post_history
- Deployment: Docker containers (backend, redis)

**Phạm vi triển khai:**

- Frontend: 2 web apps (user-web port 5173, admin-web port 5174)
- Backend: API server port 8080
- Database: MongoDB Atlas cloud
- Cache: Redis container

## 1.3. CÁC GIỚI HẠN CỦA HỆ THỐNG

### 1.3.1. Giới hạn chức năng

**Không hỗ trợ:**

- Video streaming trực tiếp
- Voice/Video call
- Tin nhắn nhóm (chỉ có 1-1)
- Mobile app native (chỉ web responsive)
- Payment/Marketplace
- Tích hợp mạng xã hội khác (ngoài Google OAuth)

**Giới hạn upload:**

- File: Tối đa 10MB/file, tối đa 10 files/bài
- Video: Tối đa 100MB, định dạng MP4/WebM
- Avatar/Banner: Lưu trên server local

**Giới hạn nội dung:**

- Nested comments: Tối đa 5 levels
- Poll options: Tối đa 10 options
- Moderators: Tối đa 25/cộng đồng
- Title: 3-300 ký tự
- Post content: Tối đa 40,000 ký tự (có thể cấu hình theo cộng đồng)
- Community name: 3-21 ký tự
- Username: 3-20 ký tự

### 1.3.2. Giới hạn kỹ thuật

**Thời gian token:**

- Access token: 1 giờ
- Refresh token: 7 ngày
- OTP code: 10 phút

**Rate limiting:**

- 100 requests/phút/user

**Database:**

- MongoDB Atlas Free Tier: 512MB storage

**Browser hỗ trợ:**

- Chrome 90+, Firefox 88+, Safari 14+
- Yêu cầu JavaScript enabled

**Ngôn ngữ:**

- Tiếng Việt và tiếng Anh

### 1.3.3. Giới hạn bảo mật

**Xác thực:**

- Chỉ hỗ trợ email/password và Google OAuth
- Không có 2FA/MFA
- Không có biometric authentication

**Quản lý dữ liệu:**

- Soft delete user: Có thể khôi phục trong 30 ngày
- Hard delete: Không hỗ trợ
- Backup: Chưa tự động

## 1.4. BIỂU MẪU VÀ QUY ĐỊNH LIÊN QUAN

### 1.4.1. Biểu mẫu hệ thống

**1. Form đăng ký (Register)**

```
Fields:
- Email: required, format email
- Username: required, 3-20 chars, alphanumeric + underscore
- Password: required, min 8 chars
- Confirm Password: required, must match
→ Gửi OTP qua email
```

**2. Form xác thực OTP**

```
Fields:
- OTP Code: required, 6 digits
- Expires: 10 minutes
```

**3. Form đăng nhập (Login)**

```
Fields:
- Email: required
- Password: required
hoặc
- Login with Google button
```

**4. Form quên mật khẩu**

```
Step 1 - Nhập email:
- Email: required

Step 2 - Xác thực OTP:
- OTP Code: required, 6 digits

Step 3 - Reset password:
- New Password: required, min 8 chars
- Confirm Password: required, must match
```

**5. Form tạo cộng đồng**

```
Fields:
- Name*: 3-21 chars, unique, alphanumeric + underscore
- Description: optional, max 500 chars
- Avatar: optional, image file
- Banner: optional, image file
- Is 18+: checkbox
- Settings:
  + is_private: boolean
  + post_require_approval: boolean
  + join_require_approval: boolean
  + max_post_length: number
- Rules: array of {title, description}
```

**6. Form tạo bài viết**

```
Fields:
- Community*: select dropdown
- Title*: 3-300 chars
- Content Type: text/image/video/poll
- Text Content: rich text, max 40,000 chars
- Images: multiple files, max 10
- Videos: max 1 file, max 100MB
- Poll (nếu chọn):
  + Question*: text
  + Options*: min 2, max 10
  + Duration: days
```

**7. Form chỉnh sửa profile**

```
Fields:
- Username: 3-20 chars (không đổi được sau khi tạo)
- Bio: max 500 chars
- Avatar: upload file
- Banner: upload file
- Location: select province
- Interests: multi-select
- Gender: select
```

**8. Form báo cáo**

```
Fields:
- Target: post_id hoặc comment_id
- Reason*: select (spam, harassment, violence, etc.)
- Description: optional, max 500 chars
```

**9. Form tin nhắn**

```
Fields:
- To User*: select from contacts
- Message*: text, max 2000 chars
```

**10. Form ban user (Moderator/Admin)**

```
Fields:
- User ID*: required
- Type*: ban hoặc mute
- Reason*: text, max 500 chars
- Duration: số ngày (0 = vĩnh viễn)
```

### 1.4.2. Quy định nội dung

**1. Quy định về tài khoản**

- Một email chỉ đăng ký được 1 tài khoản
- Username không được trùng, không chứa ký tự đặc biệt (trừ \_)
- Mật khẩu tối thiểu 8 ký tự
- Phải xác thực email bằng OTP trước khi sử dụng

**2. Quy định về cộng đồng**

- Tên cộng đồng không được trùng
- Không chứa từ ngữ nhạy cảm, vi phạm pháp luật
- Owner có quyền xóa cộng đồng
- Mỗi user tối đa tạo 10 cộng đồng

**3. Quy định về bài viết**

- Bài viết phải thuộc về 1 cộng đồng
- Tiêu đề không được để trống
- Nếu cộng đồng yêu cầu duyệt: bài viết chờ moderator phê duyệt
- Bài viết có thể chỉnh sửa (lưu lịch sử)
- Xóa bài = soft delete

**4. Nội dung cấm đăng**

- Spam, quảng cáo không phép
- Nội dung bạo lực, khủng bố
- Nội dung khiêu dâm (trừ cộng đồng 18+)
- Thông tin giả mạo
- Vi phạm bản quyền
- Ngôn từ thù ghét

**5. Quy định vote**

- Mỗi user chỉ vote 1 lần/bài viết hoặc comment
- Có thể đổi vote (upvote ↔ downvote)
- Không hiển thị danh sách user đã vote

**6. Quy định kiểm duyệt**

- Moderator có quyền duyệt/từ chối bài trong cộng đồng mình
- Owner có quyền như moderator
- Admin có quyền trên toàn hệ thống
- Lịch sử kiểm duyệt được lưu lại

**7. Quy định ban/mute**

- Ban user: Không thể post, comment trong cộng đồng
- Mute user: Không thể comment (vẫn xem được)
- Ban có thể vĩnh viễn hoặc có thời hạn
- Owner/Moderator ban trong phạm vi cộng đồng
- Admin ban toàn hệ thống

### 1.4.3. Quy định kỹ thuật

**1. API Response Format**

```json
Success:
{
  "success": true,
  "data": {...},
  "message": "Success message"
}

Error:
{
  "success": false,
  "message": "Error message",
  "error_code": "ERROR_CODE"
}
```

**2. Authentication**

- Access token gửi qua header: `Authorization: Bearer <token>`
- Refresh token: POST `/api/auth/refresh` với refresh_token
- Logout: POST `/api/auth/logout` xóa refresh token

**3. Pagination**

```
Query params:
- page: số trang (default: 1)
- limit: số items/trang (default: 10, max: 50)

Response:
{
  "data": [...],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 95,
    "page_size": 10
  }
}
```

**4. Error Codes**

- `INVALID_ID`: ID không hợp lệ
- `NOT_FOUND`: Không tìm thấy
- `UNAUTHORIZED`: Chưa đăng nhập
- `FORBIDDEN`: Không có quyền
- `VALIDATION_ERROR`: Dữ liệu không hợp lệ
- `DUPLICATE`: Trùng lặp (username, email, community name)

---

**Tóm tắt Chương 1:**

Chương 1 đã trình bày hiện trạng các diễn đàn hiện nay còn nhiều hạn chế về quản lý cộng đồng, trải nghiệm người dùng và bảo mật. Từ đó, đề tài LKForum được xây dựng với nội dung chính là hệ thống diễn đàn cộng đồng sử dụng công nghệ Svelte 5, Go, MongoDB, cung cấp đầy đủ các module: xác thực, quản lý cộng đồng, bài viết, bình luận, tin nhắn real-time, kiểm duyệt và quản trị. Hệ thống có phạm vi rõ ràng với 4 vai trò người dùng và triển khai trên nền tảng web. Các giới hạn của hệ thống được xác định về mặt chức năng (không hỗ trợ video call, tin nhắn nhóm), kỹ thuật (giới hạn file size, token expiry) và bảo mật. Cuối cùng, chương trình bày các biểu mẫu chính (đăng ký, tạo cộng đồng, đăng bài) và quy định liên quan đến nội dung, tài khoản, kiểm duyệt cùng các quy định kỹ thuật về API.

2. **Quản lý cộng đồng**
   - Tạo cộng đồng mới với tên, mô tả, avatar, banner
   - Tham gia/rời khỏi cộng đồng
   - Quản lý thành viên (ban/unban user)
   - Quản lý moderator
   - Cấu hình cộng đồng (private, post approval, join approval)
   - Quy tắc cộng đồng

3. **Quản lý bài viết**
   - Tạo bài viết (text, images, videos, poll)
   - Chỉnh sửa/xóa bài viết
   - Lưu nháp
   - Bình chọn (upvote/downvote)
   - Lưu bài viết
   - Ẩn/hiện bài viết
   - Báo cáo bài viết vi phạm

4. **Bình luận**
   - Bình luận trên bài viết
   - Trả lời bình luận (nested comments)
   - Vote bình luận
   - Xóa bình luận

5. **Tin nhắn**
   - Nhắn tin trực tiếp 1-1
   - Real-time messaging với WebSocket
   - Lịch sử tin nhắn

6. **Kiểm duyệt**
   - Duyệt/từ chối bài viết pending
   - Xem bài viết đã chỉnh sửa
   - Xử lý báo cáo
   - Ban/unban bài viết

7. **Thống kê và quản trị**
   - Thống kê tổng quan hệ thống
   - Quản lý người dùng (ban/unban/delete)
   - Quản lý cộng đồng
   - Xử lý báo cáo

**Phạm vi kỹ thuật:**

- Frontend: Web application (desktop và mobile responsive)
- Backend: RESTful API
- Database: MongoDB Atlas (cloud)
- Cache: Redis
- Deployment: Docker containers
- Version Control: Git

## 1.5. CÁC GIỚI HẠN CỦA HỆ THỐNG

### 1.5.1. Giới hạn chức năng

- Không hỗ trợ video streaming trực tiếp
- Không có chức năng gọi voice/video call
- Không hỗ trợ tin nhắn nhóm (chỉ hỗ trợ 1-1)
- Không có marketplace/e-commerce
- Không tích hợp thanh toán trực tuyến
- Không có mobile app (chỉ có web responsive)

### 1.5.2. Giới hạn kỹ thuật

**Hiệu năng:**

- Upload file: Tối đa 10MB/file, tối đa 10 files/bài viết
- Video: Tối đa 100MB, hỗ trợ định dạng MP4, WebM
- Poll: Tối đa 10 options/poll
- Nested comments: Giới hạn 5 levels
- Số lượng moderator: Tối đa 25/cộng đồng

**Bảo mật:**

- Access token: Thời gian sống 1 giờ
- Refresh token: Thời gian sống 7 ngày
- OTP code: Thời gian sống 10 phút
- Rate limiting: 100 requests/phút/user

**Lưu trữ:**

- Avatar/Cover: Lưu trên server local hoặc cloud storage (chưa triển khai CDN)
- Database: MongoDB Atlas Free Tier (giới hạn 512MB)

### 1.5.3. Giới hạn người dùng

- Chỉ hỗ trợ tiếng Việt và tiếng Anh
- Yêu cầu trình duyệt hiện đại (Chrome 90+, Firefox 88+, Safari 14+)
- Yêu cầu JavaScript được bật
- Cần kết nối internet ổn định cho real-time messaging

## 1.6. QUY ĐỊNH VÀ TIÊU CHUẨN ÁP DỤNG

### 1.6.1. Quy định về nội dung

**Nội dung cấm đăng:**

- Nội dung bạo lực, khủng bố
- Nội dung khiêu dâm (trừ cộng đồng 18+)
- Spam, quảng cáo không phép
- Thông tin sai lệch, giả mạo
- Vi phạm bản quyền
- Ngôn từ thù ghét, phân biệt đối xử

**Quy định về cộng đồng:**

- Tên cộng đồng: 3-21 ký tự, chỉ chứa chữ, số, gạch dưới
- Mô tả cộng đồng: Tối đa 500 ký tự
- Mỗi user có thể tạo tối đa 10 cộng đồng
- Cộng đồng private: Chỉ thành viên mới xem được

**Quy định về bài viết:**

- Tiêu đề: 3-300 ký tự
- Nội dung text: Tối đa 40,000 ký tự (có thể cấu hình/cộng đồng)
- Bài viết cần approval: Tùy theo cài đặt cộng đồng
- Edit history: Lưu lại lịch sử chỉnh sửa

**Quy định về tài khoản:**

- Username: 3-20 ký tự, chỉ chữ, số, gạch dưới
- Email: Phải xác thực qua OTP
- Mật khẩu: Tối thiểu 8 ký tự
- Ban user: Tạm thời hoặc vĩnh viễn
- Xóa tài khoản: Soft delete, có thể khôi phục trong 30 ngày

### 1.6.2. Tiêu chuẩn kỹ thuật

**Backend API:**

- RESTful API design principles
- HTTP status codes chuẩn
- JSON response format
- Error handling chuẩn
- API versioning (/api/v1)

**Frontend:**

- Mobile-first responsive design
- Accessibility (WCAG 2.1 Level AA)
- SEO optimization
- Performance: First Contentful Paint < 2s
- Cross-browser compatibility

**Security:**

- HTTPS/TLS encryption
- JWT authentication
- CORS policy
- Input validation và sanitization
- SQL/NoSQL injection prevention
- XSS protection
- CSRF protection

**Database:**

- MongoDB best practices
- Indexing cho các query thường dùng
- Data normalization khi cần
- Backup định kỳ

## 1.7. LỢI ÍCH CỦA ĐỀ TÀI

### 1.7.1. Lợi ích về mặt nghiên cứu và học tập

- Áp dụng kiến thức về phân tích, thiết kế hệ thống phức tạp
- Thực hành phát triển full-stack với công nghệ hiện đại
- Nghiên cứu kiến trúc microservices và real-time communication
- Tìm hiểu về bảo mật web application
- Học cách tối ưu hiệu năng và scalability

### 1.7.2. Lợi ích thực tiễn

**Cho người dùng:**

- Không gian tự do để tạo và quản lý cộng đồng
- Giao diện thân thiện, dễ sử dụng
- Bảo mật thông tin cá nhân
- Tính năng tương tác đa dạng

**Cho tổ chức/doanh nghiệp:**

- Có thể tự host và tùy chỉnh
- Quản lý cộng đồng khách hàng
- Thu thập feedback và insight
- Xây dựng brand community

**Cho developer:**

- Source code mở để tham khảo
- Kiến trúc rõ ràng, dễ mở rộng
- Tài liệu đầy đủ
- Best practices trong development

### 1.7.3. Tính mới và khả năng áp dụng

**Tính mới:**

- Sử dụng Svelte 5 với Runes API (mới nhất)
- Kiến trúc Go backend hiệu năng cao
- WebSocket cho real-time messaging
- System design theo best practices hiện đại

**Khả năng áp dụng:**

- Cộng đồng học thuật (trường học, lớp học)
- Cộng đồng doanh nghiệp nội bộ
- Cộng đồng sở thích, đam mê
- Forum hỗ trợ khách hàng
- Nền tảng Q&A chuyên ngành

## 1.8. BỐ CỤC ĐỀ TÀI

Đề tài được chia thành 5 chương:

**Chương 1: Tổng quan đề tài**
Giới thiệu, hiện trạng, mục tiêu, nội dung, phạm vi và giới hạn của hệ thống

**Chương 2: Cơ sở lý thuyết**
Trình bày các công nghệ, framework, pattern và kiến trúc được sử dụng

**Chương 3: Phân tích và thiết kế hệ thống**
Phân tích yêu cầu, thiết kế kiến trúc, cơ sở dữ liệu, API và giao diện

**Chương 4: Xây dựng và triển khai hệ thống**
Chi tiết quá trình implement các module, tính năng và deployment

**Chương 5: Kết quả và đánh giá**
Đánh giá kết quả đạt được, kiểm thử, hướng phát triển

---

**Tóm tắt chương 1:**

Chương 1 đã trình bày tổng quan về đề tài LKForum - một nền tảng diễn đàn cộng đồng hiện đại. Đề tài xuất phát từ nhu cầu thực tế về một hệ thống forum linh hoạt, dễ sử dụng và bảo mật. Hệ thống được xây dựng với các công nghệ hiện đại (Svelte 5, Go, MongoDB) và cung cấp đầy đủ các chức năng quản lý cộng đồng, bài viết, bình luận, tin nhắn real-time và kiểm duyệt nội dung. Mặc dù còn một số giới hạn về chức năng và kỹ thuật, LKForum có tiềm năng áp dụng rộng rãi trong nhiều lĩnh vực từ giáo dục đến doanh nghiệp.
