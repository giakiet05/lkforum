# TÀI LIỆU THIẾT KẾ GIAO DIỆN - LKFORUM

## MỤC LỤC

- [1. Màn hình Trang chủ (Home)](#1-màn-hình-trang-chủ-home)
- [2. Màn hình Đăng nhập (Login)](#2-màn-hình-đăng-nhập-login)
- [3. Màn hình Đăng ký (Register)](#3-màn-hình-đăng-ký-register)
- [4. Màn hình Quên mật khẩu (Forgot Password)](#4-màn-hình-quên-mật-khẩu-forgot-password)
- [5. Màn hình Xác thực Google (Google Setup)](#5-màn-hình-xác-thực-google-google-setup)
- [6. Màn hình Chi tiết bài đăng (Post Detail)](#6-màn-hình-chi-tiết-bài-đăng-post-detail)
- [7. Màn hình Chỉnh sửa bài đăng (Edit Post)](#7-màn-hình-chỉnh-sửa-bài-đăng-edit-post)
- [8. Màn hình Cộng đồng (Community)](#8-màn-hình-cộng-đồng-community)
- [9. Màn hình Tạo cộng đồng (Create Community)](#9-màn-hình-tạo-cộng-đồng-create-community)
- [10. Màn hình Cài đặt cộng đồng (Community Settings)](#10-màn-hình-cài-đặt-cộng-đồng-community-settings)
- [11. Màn hình Quản lý cộng đồng (Manage Communities)](#11-màn-hình-quản-lý-cộng-đồng-manage-communities)
- [12. Màn hình Công cụ quản trị (Mod Tools)](#12-màn-hình-công-cụ-quản-trị-mod-tools)
- [13. Màn hình Trang cá nhân (Profile)](#13-màn-hình-trang-cá-nhân-profile)
- [14. Màn hình Cài đặt tài khoản (Settings)](#14-màn-hình-cài-đặt-tài-khoản-settings)
- [15. Màn hình Tin nhắn (Messages)](#15-màn-hình-tin-nhắn-messages)
- [16. Màn hình Khám phá (Explore)](#16-màn-hình-khám-phá-explore)
- [17. Màn hình Phổ biến (Popular)](#17-màn-hình-phổ-biến-popular)

---

## 1. Màn hình Trang chủ (Home)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Trang chủ là trang đích chính khi người dùng truy cập website diễn đàn. Giao diện hiển thị danh sách bài đăng mới nhất từ các cộng đồng mà người dùng đã tham gia, kèm theo các thông tin về tác giả, cộng đồng, thời gian đăng, số lượt upvote/downvote, và số bình luận.

### Các thành phần giao diện

1.1. Thanh điều hướng trên (Topbar)

- **Logo LKForum**: Biểu tượng trang web ở góc trái
- **Thanh tìm kiếm**: Ô tìm kiếm với placeholder "Tìm kiếm bài đăng, cộng đồng..."
  - Màu nền: Trắng (mặc định), xanh khi focus
  - Icon kính lúp bên trái
- **Các nút chức năng** (góc phải):
  - Icon tạo bài mới (dấu +)
  - Icon thông báo (có badge số lượng thông báo chưa đọc)
  - Avatar người dùng (dropdown menu khi click)

1.2. Thanh điều hướng bên trái (Sidebar)

- **Menu chính**:
  - Trang chủ (Home)
  - Phổ biến (Popular)
  - Khám phá (Explore)
  - Tất cả (All)
- **Cộng đồng của bạn**:
  - Danh sách các cộng đồng đã tham gia
  - Nút "Tạo cộng đồng"
  - Nút "Quản lý cộng đồng" (nếu là creator/moderator)

1.3. Vùng nội dung chính (Main Feed)

- **Tabs lọc bài đăng**:
  - Mới nhất (Latest) - mặc định
  - Hot (bài đăng có nhiều tương tác)
  - Top (bài đăng có điểm cao nhất)
- **Danh sách bài đăng**: Hiển thị dạng card, mỗi card bao gồm:
  - Avatar và tên tác giả
  - Tên cộng đồng (dạng "lk/tên-cộng-đồng")
  - Thời gian đăng (ví dụ: "2 giờ trước")
  - Tiêu đề bài đăng
  - Nội dung preview (văn bản/hình ảnh/video/poll)
  - Thanh công cụ tương tác:
    - Nút Upvote (mũi tên lên) + số điểm + nút Downvote (mũi tên xuống)
    - Nút Comment (biểu tượng chat) + số lượng bình luận
    - Nút Share (biểu tượng chia sẻ)
    - Nút Save (biểu tượng bookmark)
    - Menu 3 chấm (Report/Hide)

1.4. Thanh bên phải (Right Sidebar)

- **Card "Tạo bài đăng"**: Nút nhanh để tạo bài mới
- **Cộng đồng đề xuất**: Danh sách các cộng đồng có thể quan tâm
- **Trending Topics**: Các chủ đề đang hot

### Tương tác người dùng

- Click vào bài đăng → Chuyển đến trang chi tiết bài đăng
- Click vào tên cộng đồng → Chuyển đến trang cộng đồng
- Click vào avatar/tên tác giả → Chuyển đến trang cá nhân
- Upvote/Downvote → Cập nhật điểm số realtime
- Comment → Chuyển đến chi tiết bài đăng, focus vào ô comment
- Share → Hiện popup chia sẻ (Copy link, Facebook, Twitter...)
- Save → Lưu bài vào danh sách "Đã lưu"
- Report → Hiện modal báo cáo vi phạm
- Hide → Ẩn bài đăng khỏi feed

### Trạng thái đặc biệt

- **Chưa đăng nhập**: Hiển thị banner mời đăng nhập/đăng ký ở sidebar
- **Không có bài đăng**: Hiển thị thông báo "Chưa có bài đăng nào. Hãy tham gia các cộng đồng để xem nội dung!"
- **Loading**: Hiển thị skeleton loading cho danh sách bài đăng

---

## 2. Màn hình Đăng nhập (Login)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Đăng nhập cho phép người dùng xác thực tài khoản trước khi truy cập các chức năng của diễn đàn. Hỗ trợ 2 phương thức đăng nhập: đăng nhập bằng Email/Mật khẩu (local) và đăng nhập bằng tài khoản Google OAuth.

### Các thành phần giao diện

2.1. Vùng chứa form (Container)

- **Logo LKForum**: Hiển thị ở trên cùng, căn giữa
- **Tiêu đề**: "Đăng nhập vào LKForum"
- **Form đăng nhập**:

2.2. Trường nhập liệu

- **Email** (bắt buộc):
  - Label: "Email"
  - Type: email
  - Placeholder: "Nhập email của bạn"
  - Validation: Format email hợp lệ
  - Hiển thị lỗi: "Email không hợp lệ" (nếu sai format)

- **Mật khẩu** (bắt buộc):
  - Label: "Mật khẩu"
  - Type: password (có nút toggle show/hide)
  - Placeholder: "Nhập mật khẩu"
  - Icon mắt bên phải để hiển thị/ẩn mật khẩu
  - Validation: Không được để trống
  - Hiển thị lỗi: "Mật khẩu không được để trống"

2.3. Các tùy chọn bổ sung

- **Checkbox "Ghi nhớ đăng nhập"**:
  - Khi chọn, lưu phiên đăng nhập lâu dài (refresh token)
  - Mặc định: không check

- **Liên kết "Quên mật khẩu?"**:
  - Hiển thị bên phải, màu xanh
  - Click → Chuyển đến màn hình Quên mật khẩu

2.4. Nút hành động

- **Nút "Đăng nhập"**:
  - Màu xanh chủ đạo, full-width
  - Disabled nếu form chưa valid
  - Click → Gửi request xác thực đến API
  - Loading state: Hiển thị spinner + text "Đang đăng nhập..."
  - Thành công → Redirect về trang trước đó hoặc Home
  - Thất bại → Hiển thị toast lỗi "Email hoặc mật khẩu không đúng"

- **Hoặc**: Đường phân cách với text "Hoặc"

- **Nút "Đăng nhập bằng Google"**:
  - Màu trắng, viền xám, có icon Google
  - Full-width
  - Click → Chuyển hướng đến Google OAuth

2.5. Liên kết đăng ký

- **Text**: "Chưa có tài khoản?"
- **Link "Đăng ký ngay"**: Màu xanh, bold
- Click → Chuyển đến màn hình Đăng ký

### Tương tác người dùng

- Nhập email/password → Validate realtime
- Click "Đăng nhập" → Xác thực tài khoản
- Click "Quên mật khẩu" → Chuyển trang reset password
- Click "Đăng nhập bằng Google" → OAuth flow
- Click "Đăng ký ngay" → Chuyển trang đăng ký

### Trạng thái đặc biệt

- **Đăng nhập thất bại**: Toast màu đỏ "Email hoặc mật khẩu không đúng"
- **Tài khoản bị ban**: Toast "Tài khoản của bạn đã bị khóa. Lý do: [ban_reason]"
- **Email chưa xác thực**: (Không áp dụng - local users luôn verified sau registration)
- **Loading**: Button disabled + spinner

---

## 3. Màn hình Đăng ký (Register)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Đăng ký cho phép người dùng tạo tài khoản mới bằng email. Quy trình gồm 3 bước: (1) Nhập email và gửi OTP, (2) Xác thực mã OTP, (3) Hoàn tất đăng ký với username và password.

### Các thành phần giao diện

3.1. Bước 1: Nhập Email

- **Tiêu đề**: "Tạo tài khoản LKForum"
- **Mô tả**: "Bước 1/3: Xác thực email"

**Trường nhập liệu:**

- **Email** (bắt buộc):
  - Label: "Email"
  - Type: email
  - Placeholder: "Nhập email của bạn"
  - Validation:
    - Format email hợp lệ
    - Email chưa được đăng ký
  - Hiển thị lỗi:
    - "Email không hợp lệ" (format sai)
    - "Email đã được sử dụng" (trùng trong DB)

**Nút hành động:**

- **Nút "Gửi mã xác thực"**:
  - Màu xanh, full-width
  - Disabled nếu email chưa valid
  - Click → Gửi OTP đến email
  - Loading: "Đang gửi mã..."
  - Thành công → Chuyển sang Bước 2
  - Thất bại → Toast lỗi

3.2. Bước 2: Xác thực OTP

- **Tiêu đề**: "Xác thực email"
- **Mô tả**: "Bước 2/3: Nhập mã OTP đã gửi đến [email]"

**Trường nhập liệu:**

- **Mã OTP** (bắt buộc):
  - 6 ô input riêng biệt, mỗi ô 1 chữ số
  - Auto-focus và tự động chuyển ô khi nhập
  - Validation: Đúng 6 chữ số
  - Hiển thị lỗi: "Mã OTP không đúng"

**Tùy chọn:**

- **Đếm ngược**: "Mã OTP có hiệu lực trong 05:00" (countdown từ 5 phút)
- **Link "Gửi lại mã"**: Chỉ hiện khi countdown = 0
  - Click → Gửi lại OTP mới, reset countdown

**Nút hành động:**

- **Nút "Xác thực"**:
  - Màu xanh, full-width
  - Disabled nếu chưa nhập đủ 6 số
  - Click → Xác thực OTP với server
  - Thành công → Chuyển sang Bước 3
  - Thất bại → Toast "Mã OTP không đúng"

3.3. Bước 3: Hoàn tất đăng ký

- **Tiêu đề**: "Hoàn tất đăng ký"
- **Mô tả**: "Bước 3/3: Tạo username và mật khẩu"

**Trường nhập liệu:**

- **Username** (bắt buộc):
  - Label: "Tên người dùng"
  - Placeholder: "Chọn username của bạn"
  - Validation:
    - 3-20 ký tự
    - Chỉ chữ, số, dấu gạch dưới
    - Chưa được sử dụng
  - Hiển thị lỗi:
    - "Username phải từ 3-20 ký tự"
    - "Username đã được sử dụng"

- **Mật khẩu** (bắt buộc):
  - Label: "Mật khẩu"
  - Type: password (có nút toggle show/hide)
  - Placeholder: "Tạo mật khẩu mạnh"
  - Validation:
    - Tối thiểu 8 ký tự
    - Có ít nhất 1 chữ hoa
    - Có ít nhất 1 chữ thường
    - Có ít nhất 1 số
    - Có ít nhất 1 ký tự đặc biệt (@$!%\*?&)
  - Hiển thị yêu cầu dưới ô input:
    - ✓ Ít nhất 8 ký tự
    - ✓ Có chữ hoa và chữ thường
    - ✓ Có số và ký tự đặc biệt
  - Hiển thị lỗi: "Mật khẩu chưa đủ mạnh"

- **Xác nhận mật khẩu** (bắt buộc):
  - Label: "Xác nhận mật khẩu"
  - Type: password
  - Placeholder: "Nhập lại mật khẩu"
  - Validation: Trùng với mật khẩu đã nhập
  - Hiển thị lỗi: "Mật khẩu không khớp"

**Nút hành động:**

- **Nút "Hoàn tất đăng ký"**:
  - Màu xanh, full-width
  - Disabled nếu form chưa valid
  - Click → Tạo tài khoản
  - Loading: "Đang tạo tài khoản..."
  - Thành công → Tự động đăng nhập + redirect Home
  - Thất bại → Toast lỗi

3.4. Phương thức đăng ký khác

- **Đường phân cách**: "Hoặc"
- **Nút "Đăng ký bằng Google"**:
  - Màu trắng, viền xám, icon Google
  - Click → Google OAuth

3.5. Link đăng nhập

- **Text**: "Đã có tài khoản?"
- **Link "Đăng nhập ngay"**: Màu xanh
- Click → Chuyển về màn hình Đăng nhập

### Tương tác người dùng

- Bước 1: Nhập email → Gửi OTP
- Bước 2: Nhập OTP → Xác thực
- Bước 3: Nhập username/password → Hoàn tất
- Có thể back về bước trước bằng nút "Quay lại"

### Trạng thái đặc biệt

- **Email đã tồn tại**: Toast "Email đã được đăng ký"
- **OTP hết hạn**: Toast "Mã OTP đã hết hạn, vui lòng gửi lại"
- **Username đã tồn tại**: Toast "Username đã được sử dụng"
- **Mật khẩu yếu**: Hiện danh sách yêu cầu chưa đạt

---

## 4. Màn hình Quên mật khẩu (Forgot Password)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Quên mật khẩu cho phép người dùng đặt lại mật khẩu khi quên mật khẩu cũ. Quy trình gồm 3 bước: (1) Nhập email, (2) Xác thực OTP, (3) Đặt mật khẩu mới.

### Các thành phần giao diện

4.1. Bước 1: Nhập Email

- **Tiêu đề**: "Quên mật khẩu"
- **Mô tả**: "Nhập email đã đăng ký để nhận mã xác thực"

**Trường nhập:**

- **Email** (bắt buộc):
  - Placeholder: "Nhập email của bạn"
  - Validation: Email hợp lệ và đã đăng ký
  - Lỗi: "Email không tồn tại trong hệ thống"

**Nút:**

- **"Gửi mã xác thực"**: Click → Gửi OTP, chuyển Bước 2

4.2. Bước 2: Xác thực OTP

- **Mô tả**: "Nhập mã OTP đã gửi đến [email]"
- **6 ô nhập OTP**
- **Countdown 5 phút**
- **Link "Gửi lại mã"**
- **Nút "Xác thực"**: Click → Verify OTP, chuyển Bước 3

4.3. Bước 3: Đặt mật khẩu mới

- **Tiêu đề**: "Đặt mật khẩu mới"

**Trường nhập:**

- **Mật khẩu mới** (bắt buộc):
  - Validation: Giống đăng ký (8+ ký tự, chữ hoa, chữ thường, số, ký tự đặc biệt)
  - Show/hide password

- **Xác nhận mật khẩu** (bắt buộc):
  - Validation: Trùng với mật khẩu mới
  - Lỗi: "Mật khẩu không khớp"

**Nút:**

- **"Đặt lại mật khẩu"**:
  - Click → Cập nhật mật khẩu
  - Thành công → Toast "Đặt lại mật khẩu thành công" + redirect Login

### Tương tác

- OTP expires sau 5 phút
- Có thể gửi lại OTP tối đa 3 lần
- Sau khi reset thành công → Tự động logout các session khác

---

## 5. Màn hình Xác thực Google (Google Setup)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình này xuất hiện khi người dùng đăng nhập lần đầu bằng Google OAuth. Yêu cầu người dùng chọn username duy nhất để hoàn tất việc tạo tài khoản.

### Các thành phần giao diện

5.1. Thông tin Google đã lấy được

- **Avatar Google**: Hiển thị ảnh đại diện từ Google
- **Email Google**: Hiển thị email (read-only, không thể chỉnh sửa)
- **Text**: "Chào mừng! Vui lòng chọn username để hoàn tất đăng ký"

5.2. Trường nhập username

- **Username** (bắt buộc):
  - Label: "Tên người dùng"
  - Placeholder: "Chọn username của bạn"
  - Validation:
    - 3-20 ký tự
    - Chỉ chữ, số, gạch dưới
    - Chưa được sử dụng
  - Real-time check: Icon ✓ hoặc ✗ bên cạnh
  - Lỗi: "Username đã được sử dụng"

5.3. Nút hành động

- **Nút "Hoàn tất"**:
  - Màu xanh, full-width
  - Disabled nếu username chưa valid
  - Click → Tạo tài khoản Google + đăng nhập
  - Thành công → Redirect Home

### Tương tác

- User không thể quay lại (vì đã OAuth)
- Username phải unique
- Sau khi hoàn tất → Tạo user với provider="google"

---

## 6. Màn hình Chi tiết bài đăng (Post Detail)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Chi tiết bài đăng hiển thị toàn bộ nội dung của một bài đăng và danh sách bình luận. Cho phép người dùng tương tác (upvote, comment, share, save) và xem các bình luận theo cấu trúc nested (reply).

### Các thành phần giao diện

6.1. Thông tin bài đăng (Post Card)

- **Header**:
  - Avatar và username tác giả
  - Tên cộng đồng (link "lk/tên-cộng-đồng")
  - Thời gian đăng (ví dụ: "3 giờ trước")
  - Menu 3 chấm (nếu là chủ bài hoặc mod):
    - Edit (Chỉnh sửa)
    - Delete (Xóa bài)
    - Report (Báo cáo - cho user khác)

- **Nội dung**:
  - **Tiêu đề**: Font lớn, bold (h1)
  - **Nội dung văn bản**: Hiển thị đầy đủ (không truncate)
  - **Hình ảnh**:
    - Hiển thị dạng gallery nếu nhiều ảnh
    - Click để xem fullscreen
    - Có nút Previous/Next nếu > 1 ảnh
  - **Video**:
    - Embedded player với controls
    - Autoplay tắt
  - **Poll** (nếu là bài poll):
    - Câu hỏi poll
    - Danh sách options với radio/checkbox
    - Hiển thị % vote và số vote của mỗi option
    - Nút "Vote" (nếu chưa vote)
    - Disabled nếu đã vote hoặc poll hết hạn
    - Text "Poll ends in [time]" hoặc "Poll ended"

- **Thanh tương tác**:
  - **Upvote/Downvote**: Mũi tên lên/xuống + số điểm
    - Màu: Upvoted = xanh, Downvoted = đỏ, chưa vote = xám
  - **Comment**: Icon chat + số lượng comment
  - **Share**: Icon chia sẻ
    - Click → Popup với options: Copy link, Facebook, Twitter
    - Copy link tự động copy URL dạng slug: `/post/title-slug-{id}`
  - **Save**: Icon bookmark
    - Filled nếu đã save, outline nếu chưa
  - **Menu 3 chấm**: Report / Hide

6.2. Phần bình luận (Comment Section)

- **Header**:
  - Text "Bình luận ([số lượng])"
  - Dropdown sắp xếp: "Mới nhất", "Cũ nhất", "Hot nhất"

- **Ô nhập bình luận**:
  - Avatar người dùng hiện tại
  - Textarea: "Viết bình luận..."
  - Nút "Bình luận" (màu xanh)
  - Disabled nếu chưa đăng nhập → Hiện text "Đăng nhập để bình luận"

- **Danh sách bình luận**:
  - Hiển thị theo cấu trúc nested (parent-child)
  - Mỗi comment bao gồm:
    - Avatar và username
    - Thời gian comment
    - Nội dung comment
    - Upvote/Downvote + điểm
    - Nút "Reply" → Hiện ô reply nested
    - Menu 3 chấm (nếu là chủ comment hoặc mod):
      - Edit
      - Delete
      - Report

  - **Reply comments**:
    - Indent sang phải
    - Hiển thị tối đa 2-3 levels
    - Nếu có nhiều replies → Link "Xem thêm [n] replies"

6.3. Sidebar bên phải

- **Thông tin cộng đồng**:
  - Avatar cộng đồng
  - Tên cộng đồng
  - Mô tả ngắn
  - Số thành viên
  - Nút "Join" / "Joined" (toggle)
  - Link "Xem cộng đồng"

- **Quy tắc cộng đồng**:
  - Danh sách các rules (nếu có)

### Tương tác người dùng

- Click Upvote/Downvote → Cập nhật điểm realtime
- Click Comment → Focus vào ô comment
- Click Reply → Hiện nested comment box
- Click Share → Popup chia sẻ
- Click Save → Toggle save state
- Submit comment → Thêm vào danh sách realtime
- Click Edit → Chuyển sang edit mode
- Click Delete → Confirm popup → Xóa bài

### Trạng thái đặc biệt

- **Bài đã xóa**: Hiển thị "[Đã xóa]" thay vì nội dung
- **Chưa đăng nhập**: Các nút vote/comment disabled
- **Đã vote poll**: Hiển thị kết quả, không thể vote lại
- **Comment đã xóa**: "[Đã xóa]" nhưng vẫn giữ replies

---

## 7. Màn hình Chỉnh sửa bài đăng (Edit Post)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Chỉnh sửa cho phép tác giả bài đăng chỉnh sửa tiêu đề và nội dung. Chỉ cho phép chỉnh sửa trong thời gian nhất định (ví dụ: 24h sau khi đăng).

### Các thành phần giao diện

7.1. Form chỉnh sửa

- **Tiêu đề**: "Chỉnh sửa bài đăng"
- **Breadcrumb**: lk/tên-cộng-đồng > Bài đăng gốc

**Trường nhập:**

- **Tiêu đề** (bắt buộc):
  - Input text
  - Max length: 300 ký tự
  - Counter: "250/300"
  - Validation: Không để trống

- **Nội dung**:
  - **Text**: Textarea với rich text editor (nếu có)
  - **Hình ảnh**:
    - Hiển thị ảnh hiện tại
    - Nút "Xóa ảnh"
    - Nút "Thêm ảnh" (tối đa 10)
  - **Video**: Tương tự
  - **Poll**: Không cho phép chỉnh sửa (disable)

- **Tags**:
  - Input chip, có thể thêm/xóa tags
  - Max 5 tags

7.2. Nút hành động

- **Nút "Hủy"**:
  - Màu xám
  - Click → Confirm "Bỏ các thay đổi?" → Back

- **Nút "Lưu thay đổi"**:
  - Màu xanh
  - Disabled nếu không có thay đổi
  - Click → Cập nhật bài đăng
  - Thành công → Toast "Đã cập nhật" + redirect về PostDetail
  - Set `is_edited = true`

### Tương tác

- Chỉ tác giả hoặc mod mới có quyền edit
- Không cho edit poll sau khi có người vote
- Track changes để enable/disable nút Save

---

## 8. Màn hình Cộng đồng (Community)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Cộng đồng hiển thị thông tin chi tiết về một cộng đồng và danh sách bài đăng trong cộng đồng đó. Cho phép người dùng tham gia/rời cộng đồng, tạo bài mới, và truy cập công cụ quản trị (nếu là moderator).

### Các thành phần giao diện

8.1. Header cộng đồng

- **Banner**: Ảnh banner full-width (nếu có)
- **Thông tin chính**:
  - Avatar cộng đồng (góc trái)
  - Tên cộng đồng: "lk/tên-cộng-đồng"
  - Mô tả ngắn
  - **Thống kê**:
    - [số] thành viên
    - [số] bài đăng
    - Ngày tạo
  - **Nút hành động** (góc phải):
    - **Nút "Join" / "Joined"**:
      - Xanh nếu chưa join
      - Xám + checkmark nếu đã join
      - Click → Toggle membership
    - **Nút "Create Post"**: Tạo bài mới trong community này
    - **Menu 3 chấm** (nếu là creator/moderator):
      - Mod Tools (Công cụ quản trị)
      - Community Settings (Cài đặt)

8.2. Tabs điều hướng

- **Posts** (mặc định): Danh sách bài đăng
- **About**: Thông tin cộng đồng
- **Rules**: Quy tắc cộng đồng
- **Moderators**: Danh sách moderators

8.3. Tab Posts

- **Filter dropdown**:
  - Hot (mặc định)
  - New
  - Top (Today, This Week, This Month, All Time)

- **Danh sách bài đăng**: Giống Home feed

8.4. Tab About

- **Mô tả đầy đủ**: Markdown content
- **Thông tin chi tiết**:
  - Ngày tạo
  - Người tạo (link profile)
  - Loại cộng đồng: Public/Private, 18+
- **Tags**: Các tags của community

8.5. Tab Rules

- **Danh sách quy tắc**:
  - Mỗi rule: Số thứ tự, tiêu đề, mô tả
  - Ví dụ:
    - 1. Lịch sự và tôn trọng
    - 2. Không spam
    - 3. Nội dung liên quan đến chủ đề

8.6. Tab Moderators

- **Danh sách moderators**:
  - Avatar, username
  - Badge "Creator" hoặc "Moderator"
  - Thời gian bổ nhiệm
- **Nút "Become a moderator"** (nếu đủ điều kiện)

8.7. Sidebar

- **Quick actions**:
  - Create Post (nổi bật)
  - Community Settings (nếu là mod)
  - Report Community

- **Community stats**:
  - Members online: [số] online
  - Rank: Top [số]% communities

### Tương tác

- Click "Join" → Tham gia cộng đồng, feed sẽ hiện bài đăng từ community này
- Click "Joined" → Rời cộng đồng (confirm popup)
- Click "Create Post" → Modal/page tạo bài mới
- Click "Mod Tools" → Trang quản trị (chỉ mod)
- Click "Community Settings" → Cài đặt community (chỉ creator/mod)

### Trạng thái đặc biệt

- **Private community**:
  - Chưa join → Không xem được bài đăng
  - Hiển thị "Community này là riêng tư. Join để xem nội dung"
- **18+ community**:
  - Hiển thị badge "18+"
  - Confirm popup "Bạn có đủ 18 tuổi?"
- **Banned community**: Hiển thị "Community này đã bị cấm. Lý do: [reason]"

---

## 9. Màn hình Tạo cộng đồng (Create Community)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Tạo cộng đồng cho phép người dùng tạo một cộng đồng mới với tên, mô tả, avatar, banner và các cài đặt ban đầu.

### Các thành phần giao diện

9.1. Form tạo cộng đồng

- **Tiêu đề**: "Tạo cộng đồng mới"
- **Mô tả**: "Tạo không gian cho cộng đồng của bạn"

9.2. Thông tin cơ bản

**Trường nhập:**

- **Tên cộng đồng** (bắt buộc):
  - Label: "Tên cộng đồng"
  - Placeholder: "svelte, javascript, vietnam-travel..."
  - Prefix: "lk/"
  - Validation:
    - 3-21 ký tự
    - Chỉ chữ thường, số, gạch ngang
    - Unique (chưa tồn tại)
  - Real-time check availability
  - Lỗi: "Tên cộng đồng đã tồn tại"
  - Helper text: "Tên không thể thay đổi sau khi tạo"

- **Mô tả** (tùy chọn):
  - Label: "Mô tả ngắn"
  - Textarea, max 500 ký tự
  - Placeholder: "Mô tả về cộng đồng của bạn..."
  - Counter: "250/500"

9.3. Hình ảnh

- **Avatar cộng đồng** (tùy chọn):
  - Upload area với drag & drop
  - Preview avatar (circle)
  - Nút "Tải lên"
  - Yêu cầu: JPG/PNG, max 5MB, tỷ lệ 1:1

- **Banner cộng đồng** (tùy chọn):
  - Upload area
  - Preview banner (rectangle)
  - Yêu cầu: JPG/PNG, max 10MB, tỷ lệ 3:1

9.4. Cài đặt

**Checkboxes:**

- **Community type**:
  - [ ] Private community (Chỉ thành viên mới xem được bài đăng)
  - [ ] 18+ content (Nội dung người lớn)

- **Moderation**:
  - [ ] Bài đăng cần duyệt trước khi hiển thị
  - [ ] Thành viên mới cần duyệt trước khi tham gia

9.5. Nút hành động

- **Nút "Hủy"**:
  - Màu xám
  - Click → Back hoặc close modal

- **Nút "Tạo cộng đồng"**:
  - Màu xanh, nổi bật
  - Disabled nếu tên chưa valid
  - Click → Tạo community
  - Loading: "Đang tạo..."
  - Thành công → Toast "Đã tạo cộng đồng lk/[name]" + redirect đến community page

### Tương tác

- Upload ảnh → Preview realtime
- Check tên community → Hiện icon ✓/✗
- Tạo thành công → User tự động là creator + moderator

### Trạng thái đặc biệt

- **Tên trùng**: Toast "Tên cộng đồng đã tồn tại"
- **Upload failed**: Toast "Lỗi upload ảnh"

---

## 10. Màn hình Cài đặt cộng đồng (Community Settings)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Cài đặt cộng đồng cho phép creator và moderators chỉnh sửa thông tin, cài đặt quyền riêng tư, và quản lý quy tắc của cộng đồng. Chỉ creator và moderators mới có quyền truy cập.

### Các thành phần giao diện

10.1. Sidebar menu

- **General** (Chung)
- **Appearance** (Giao diện)
- **Rules** (Quy tắc)
- **Moderation** (Kiểm duyệt)
- **Advanced** (Nâng cao)

10.2. Tab General (Chung)

**Trường chỉnh sửa:**

- **Tên cộng đồng**:
  - Display: "lk/[name]"
  - Không thể chỉnh sửa (read-only)
  - Helper: "Tên cộng đồng không thể thay đổi"

- **Mô tả**:
  - Textarea, max 500 ký tự
  - Counter

- **Tags**:
  - Input chip
  - Max 5 tags
  - Suggestions: programming, lifestyle, gaming...

10.3. Tab Appearance (Giao diện)

- **Avatar cộng đồng**:
  - Preview hiện tại
  - Nút "Change avatar"
  - Upload modal

- **Banner cộng đồng**:
  - Preview hiện tại
  - Nút "Change banner"
  - Upload modal

- **Theme color**:
  - Color picker
  - Preview theme

10.4. Tab Rules (Quy tắc)

- **Danh sách quy tắc hiện tại**:
  - Mỗi rule: Tiêu đề + Mô tả + Nút Edit/Delete
  - Drag & drop để sắp xếp thứ tự

- **Nút "Add rule"**:
  - Click → Modal thêm rule mới
  - Form: Title (bắt buộc), Description (tùy chọn)

10.5. Tab Moderation (Kiểm duyệt)

**Cài đặt kiểm duyệt:**

- **Post moderation**:
  - [ ] Require approval for new posts
  - [ ] Auto-remove posts with [n] reports

- **Member moderation**:
  - [ ] Require approval for new members
  - [ ] Minimum account age: [n] days

- **Content filters**:
  - [ ] Block links
  - [ ] Block images
  - [ ] Filter banned words
  - Textarea: "Danh sách từ cấm (mỗi dòng 1 từ)"

10.6. Tab Advanced (Nâng cao)

**Cài đặt nâng cao:**

- **Privacy**:
  - Radio buttons:
    - ( ) Public (Ai cũng xem được)
    - ( ) Private (Chỉ members)
    - ( ) Restricted (Chỉ moderators duyệt)

- **Content rating**:
  - [ ] 18+ content

- **Danger zone**:
  - **Nút "Archive community"**: Ẩn khỏi tìm kiếm, không cho post mới
  - **Nút "Delete community"**: Xóa vĩnh viễn (chỉ creator)
    - Confirm popup: "Nhập tên community để xác nhận xóa"

10.7. Footer actions

- **Nút "Cancel"**: Hủy thay đổi, back
- **Nút "Save changes"**:
  - Màu xanh
  - Disabled nếu không có thay đổi
  - Click → Cập nhật settings
  - Toast "Đã cập nhật"

### Tương tác

- Chỉ creator/moderators mới truy cập được
- Thay đổi settings → Enable nút Save
- Delete community → Confirm rất kỹ

---

## 11. Màn hình Quản lý cộng đồng (Manage Communities)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Quản lý cộng đồng hiển thị danh sách các cộng đồng mà người dùng là creator hoặc moderator. Cho phép truy cập nhanh vào các công cụ quản trị.

### Các thành phần giao diện

11.1. Header

- **Tiêu đề**: "Quản lý cộng đồng"
- **Mô tả**: "Các cộng đồng bạn đang quản lý"

11.2. Danh sách cộng đồng

**Mỗi community card:**

- **Thông tin**:
  - Avatar cộng đồng
  - Tên: "lk/[name]"
  - Role badge: "Creator" hoặc "Moderator"
  - Stats:
    - [số] members
    - [số] posts
    - [số] unmoderated posts (nếu require approval)

- **Quick actions**:
  - Nút "View community"
  - Nút "Mod Tools"
  - Nút "Settings" (chỉ creator/mod)

11.3. Empty state

- **Nếu không có community**:
  - Icon + Text "Bạn chưa quản lý cộng đồng nào"
  - Nút "Tạo cộng đồng mới"

### Tương tác

- Click community card → View community
- Click "Mod Tools" → Trang mod tools
- Click "Settings" → Community settings

---

## 12. Màn hình Công cụ quản trị (Mod Tools)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Công cụ quản trị cung cấp các công cụ để moderators quản lý bài đăng, members, reports và bans trong cộng đồng. Chỉ creator và moderators mới có quyền truy cập.

### Các thành phần giao diện

12.1. Sidebar menu

- **Queue** (Hàng đợi duyệt)
- **Reports** (Báo cáo)
- **Banned Users** (Users bị cấm)
- **Moderators** (Quản trị viên)
- **Moderation Log** (Lịch sử)

12.2. Tab Queue (Hàng đợi)

**Nếu post_require_approval = true:**

- **Danh sách bài đăng chờ duyệt**:
  - Preview bài đăng
  - Thông tin: Tác giả, thời gian submit
  - **Actions**:
    - Nút "Approve" (màu xanh)
    - Nút "Reject" (màu đỏ)
    - Nút "Remove" (xóa)

- **Empty state**: "Không có bài đăng chờ duyệt"

12.3. Tab Reports (Báo cáo)

- **Danh sách reports**:
  - **Mỗi report card**:
    - Type: Post / Comment / User
    - Target: Link đến đối tượng bị report
    - Reporter: Username
    - Reason: Spam, Harassment, Violence...
    - Description: Mô tả chi tiết
    - Thời gian report

  - **Actions**:
    - Nút "View" (xem đối tượng bị report)
    - Nút "Dismiss" (bỏ qua report)
    - Nút "Remove content" (xóa nội dung)
    - Nút "Ban user" (cấm user)

12.4. Tab Banned Users

- **Danh sách users bị cấm**:
  - **Mỗi ban entry**:
    - Avatar + Username
    - Ban type: Banned (không vào được) / Muted (không post được)
    - Reason: Lý do cấm
    - Banned by: Mod thực hiện
    - Expires: Thời gian hết hạn hoặc "Permanent"

  - **Actions**:
    - Nút "Unban" (gỡ cấm)
    - Nút "Edit ban" (thay đổi thời gian/lý do)

- **Nút "Ban a user"**:
  - Click → Modal nhập username, lý do, thời gian

12.5. Tab Moderators

- **Danh sách moderators**:
  - Avatar + Username
  - Role: Creator / Moderator
  - Assigned date
  - Status: Active / Inactive

- **Actions** (chỉ creator):
  - Nút "Add moderator"
  - Nút "Remove" (xóa mod)
  - Toggle "Active" (tạm ngưng quyền mod)

12.6. Tab Moderation Log

- **Lịch sử các hành động mod**:
  - Timestamp
  - Moderator: Ai thực hiện
  - Action: Approve/Reject/Remove/Ban...
  - Target: Đối tượng bị tác động
  - Reason: Lý do (nếu có)

- **Filters**:
  - By moderator
  - By action type
  - By date range

### Tương tác

- Approve post → Bài đăng hiển thị công khai
- Reject post → Notify tác giả
- Dismiss report → Xóa report khỏi queue
- Ban user → User không thể post/comment trong community

---

## 13. Màn hình Trang cá nhân (Profile)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Trang cá nhân hiển thị thông tin công khai của người dùng, bao gồm avatar, bio, stats, và danh sách bài đăng/bình luận. Có 2 chế độ xem: profile của mình (có nút Edit) và profile người khác.

### Các thành phần giao diện

13.1. Header Profile

- **Cover image**: Ảnh bìa full-width (nếu có)
- **Avatar**: Ảnh đại diện (circle), overlap với cover
- **Thông tin chính**:
  - **Username**: @username
  - **Badge**: Admin / Moderator (nếu có)
  - **Bio**: Mô tả ngắn về bản thân
  - **Stats**:
    - [số] bài đăng
    - [số] bình luận
    - [số] điểm uy tín (reputation)
    - Tham gia từ [date]

- **Thông tin bổ sung** (nếu public):
  - Location: [tỉnh/thành phố]
  - Gender: [nam/nữ/không tiết lộ]
  - Birthday: [ngày/tháng/năm]
  - Interests: [tags]

- **Social links**:
  - Icon website, Facebook, YouTube, GitHub
  - Click → Mở link mới

- **Nút hành động** (góc phải):
  - **Nếu là profile của mình**:
    - Nút "Edit Profile" (chuyển Settings)
  - **Nếu là profile người khác**:
    - Nút "Send Message" (mở chat)
    - Nút "Report User" (báo cáo)

13.2. Tabs

- **Posts** (mặc định): Danh sách bài đăng của user
- **Comments**: Danh sách bình luận
- **Saved** (chỉ hiện khi xem profile của mình): Bài đăng đã lưu
- **About**: Thông tin chi tiết

13.3. Tab Posts

- **Filter**:
  - New (mới nhất)
  - Hot (nhiều tương tác)
  - Top (điểm cao)

- **Danh sách bài đăng**: Giống Home feed
- **Empty state**: "Chưa có bài đăng nào"

13.4. Tab Comments

- **Danh sách comments**:
  - Mỗi comment hiển thị:
    - Nội dung comment
    - Link đến bài đăng gốc
    - Upvotes
    - Thời gian

13.5. Tab Saved (chỉ profile của mình)

- **Danh sách bài đăng đã save**
- **Nút unsave** trên mỗi bài

13.6. Tab About

- **Thông tin chi tiết**:
  - Email (nếu public)
  - Date of birth
  - Gender
  - Location
  - Interests
  - Social links

- **Activity stats**:
  - Total posts
  - Total comments
  - Total upvotes received
  - Karma breakdown

### Tương tác

- Click "Edit Profile" → Settings page
- Click "Send Message" → Messages page
- Click post → Post detail
- Click community → Community page

### Trạng thái đặc biệt

- **Private profile**: "Profile này là riêng tư"
- **Banned user**: Banner "User này đã bị cấm"
- **Deleted account**: "Tài khoản đã bị xóa"

---

## 14. Màn hình Cài đặt tài khoản (Settings)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Cài đặt cho phép người dùng tùy chỉnh thông tin cá nhân, preferences, privacy, và các cài đặt tài khoản khác.

### Các thành phần giao diện

14.1. Sidebar menu

- **Profile** (Hồ sơ)
- **Account** (Tài khoản)
- **Appearance** (Giao diện)
- **Notifications** (Thông báo)
- **Privacy** (Riêng tư)
- **Blocked Users** (Chặn)

14.2. Tab Profile (Hồ sơ)

**Upload ảnh:**

- **Avatar**:
  - Preview circle
  - Nút "Change avatar"
  - Upload modal với crop tool

- **Cover image**:
  - Preview rectangle
  - Nút "Change cover"

**Thông tin cơ bản:**

- **Username**:
  - Read-only (không thể đổi)
  - Helper: "Username không thể thay đổi"

- **Bio**:
  - Textarea, max 200 ký tự
  - Counter

- **Personal info**:
  - Gender: Dropdown (Nam/Nữ/Không tiết lộ)
  - Date of birth: Date picker
  - Location: Dropdown (63 tỉnh/thành VN)

- **Interests**:
  - Multi-select chips
  - Max 10 interests

**Social links:**

- Website: URL input
- Facebook: URL input
- YouTube: URL input
- GitHub: URL input

Nút "Save changes": Màu xanh

14.3. Tab Account (Tài khoản)

**Email:**

- Display: email@example.com
- Nút "Change email" (nếu local account)
  - Flow: Verify OTP → Update email

**Password:**

- **Nếu local account**:
  - Nút "Change password"
  - Modal:
    - Current password
    - New password
    - Confirm password

- **Nếu Google account**:
  - Text: "Đăng nhập bằng Google"
  - Không có password

**Danger zone:**

- **Deactivate account**:
  - Nút "Deactivate account"
  - Confirm: "Tài khoản sẽ bị ẩn trong 30 ngày trước khi xóa vĩnh viễn"

- **Delete account**:
  - Nút "Delete account" (màu đỏ)
  - Confirm popup: Nhập password + "DELETE" để xác nhận

14.4. Tab Appearance (Giao diện)

**Theme:**

- Radio buttons:
  - ( ) Light
  - ( ) Dark
  - ( ) Auto (theo hệ thống)

**Font size:**

- Radio buttons:
  - ( ) Small
  - ( ) Medium (mặc định)
  - ( ) Large

Preview: Live preview khi thay đổi

14.5. Tab Notifications (Thông báo)

**Channels:**

- [ ] In-app notifications (thông báo trên website)
- [ ] Email notifications (gửi email)

**Notify me when:**

- [ ] Someone comments on my post
- [ ] Someone mentions me
- [ ] Someone upvotes my post/comment
- [ ] Someone sends me a message

**Email digest:**

- Dropdown: Off / Daily / Weekly
- Text: "Nhận tổng hợp hoạt động qua email"

14.6. Tab Privacy (Riêng tư)

**Profile visibility:**

- [ ] Show my profile publicly
- [ ] Show my email on profile
- [ ] Show my post history
- [ ] Show my comment history

**Interactions:**

- [ ] Allow direct messages from anyone
- [ ] Allow mentions from anyone
- [ ] Allow others to follow me (future feature)

**Content:**

- [ ] Show NSFW content
- [ ] Blur NSFW images

14.7. Tab Blocked Users

**Danh sách users đã chặn:**

- Avatar + Username
- Ngày chặn
- Nút "Unblock"

Nút "Block a user":

- Input username
- Nút "Block"

Empty state: "Bạn chưa chặn ai"

### Tương tác

- Thay đổi settings → Auto-save hoặc nút Save
- Upload avatar → Crop tool → Preview → Save
- Change email → OTP flow
- Change password → Verify old password
- Block user → Không thấy bài/comment của họ

---

## 15. Màn hình Tin nhắn (Messages)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Tin nhắn cho phép người dùng nhắn tin trực tiếp với nhau. Giao diện chia làm 2 phần: danh sách conversations bên trái và chat box bên phải.

### Các thành phần giao diện

15.1. Sidebar - Danh sách conversations

- **Header**:
  - Text "Messages"
  - Nút "New message" (icon +)

- **Search bar**:
  - Placeholder "Tìm cuộc trò chuyện..."

- **Danh sách conversations**:
  - Mỗi conversation item:
    - Avatar người chat
    - Username
    - Message preview (truncate)
    - Timestamp
    - Badge số message chưa đọc (nếu có)
    - Active state nếu đang chọn

15.2. Chat box - Vùng chat chính

**Nếu chưa chọn conversation:**

- Empty state: "Chọn một cuộc trò chuyện hoặc bắt đầu mới"

**Nếu đã chọn conversation:**

- **Header**:
  - Avatar + Username người chat
  - Status: Online/Offline (nếu có)
  - Menu 3 chấm:
    - View profile
    - Block user
    - Delete conversation

- **Message area**:
  - Danh sách messages (scroll to bottom)
  - Mỗi message:
    - Avatar (nếu là message của người khác)
    - Nội dung message
    - Timestamp (khi hover)
    - Align left (người khác) / right (mình)
    - Màu: Xám (người khác), xanh (mình)

  - **System messages**:
    - Căn giữa, màu xám nhạt
    - Ví dụ: "[User] đã rời cuộc trò chuyện"

- **Input area** (bottom):
  - Textarea: "Nhập tin nhắn..."
  - Auto-resize khi nhập nhiều dòng
  - Icon emoji picker
  - Nút "Send" (icon máy bay giấy)
  - Enter để gửi, Shift+Enter để xuống dòng

15.3. Modal "New Message"

- **Tiêu đề**: "Tin nhắn mới"
- **Search user**:
  - Input "Tìm người dùng..."
  - Autocomplete dropdown
  - Hiển thị: Avatar + Username
- **Select user** → Tạo conversation mới

### Tương tác

- Click conversation → Load messages
- Type message → Realtime typing indicator
- Send message → Append to chat realtime
- New message received → Badge + notification sound
- Scroll up → Load older messages (pagination)

### Trạng thái đặc biệt

- **No conversations**: "Bạn chưa có cuộc trò chuyện nào"
- **User blocked you**: "Bạn không thể nhắn tin cho người dùng này"
- **User deactivated**: "Tài khoản này đã bị vô hiệu hóa"

---

## 16. Màn hình Khám phá (Explore)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Khám phá giúp người dùng tìm kiếm cộng đồng mới, trending topics, và các bài đăng nổi bật từ toàn bộ diễn đàn (không chỉ communities đã join).

### Các thành phần giao diện

16.1. Header

- **Tiêu đề**: "Khám phá"
- **Search bar**:
  - Placeholder "Tìm cộng đồng, bài đăng..."
  - Search suggestions khi typing

16.2. Trending Topics

- **Card "Trending Topics"**:
  - Danh sách top 10 topics hot
  - Mỗi topic:
    - Hashtag: #topic-name
    - Số bài đăng: [số] posts today
    - Growth indicator: ↑ [số]%

16.3. Recommended Communities

- **Card "Cộng đồng đề xuất"**:
  - Grid layout (3 columns)
  - Mỗi community card:
    - Avatar
    - Name: lk/name
    - Description (truncate)
    - [số] members
    - Nút "Join"

16.4. Explore Feed

- **Tabs**:
  - All (tất cả bài đăng)
  - Trending (đang hot)
  - New (mới nhất)

- **Danh sách bài đăng**: Giống Home feed
- **Infinite scroll**: Load more khi scroll xuống

### Tương tác

- Click topic → Filter posts by topic
- Click community → View community
- Join community → Thêm vào "My communities"

---

## 17. Màn hình Phổ biến (Popular)

**[Hình ảnh minh họa - thêm sau]**

### Mô tả chức năng

Màn hình Phổ biến hiển thị các bài đăng có nhiều tương tác nhất trong khoảng thời gian được chọn (hôm nay, tuần này, tháng này, mọi lúc).

### Các thành phần giao diện

17.1. Header

- **Tiêu đề**: "Phổ biến"
- **Time filter dropdown**:
  - Today (hôm nay)
  - This Week (tuần này)
  - This Month (tháng này)
  - All Time (mọi lúc)

17.2. Top Posts

- **Danh sách bài đăng**:
  - Sắp xếp theo điểm (votes_count.up - votes_count.down)
  - Hiển thị ranking number (#1, #2, #3...)
  - Badge "🔥 Hot" cho top 3

- **Post cards**: Giống Home feed

17.3. Sidebar

- **Top Communities**: Cộng đồng có nhiều thành viên nhất
- **Rising Stars**: Cộng đồng đang tăng trưởng nhanh

### Tương tác

- Change time filter → Reload posts
- Click post → View detail
- Upvote/Downvote → Update ranking realtime

---

## THÀNH PHẦN CHUNG (SHARED COMPONENTS)

### C1. Component Topbar (Thanh điều hướng trên)

Xuất hiện ở: Tất cả các màn hình khi đã đăng nhập

**Thành phần:**

- Logo LKForum (góc trái)
- Search bar (giữa)
- Notifications icon (có badge)
- Create post icon (+)
- User avatar dropdown (góc phải)

**Dropdown menu:**

- Profile
- Settings
- My Communities
- Messages

---

- Logout

### C2. Component Sidebar (Menu bên trái)

Xuất hiện ở: Home, Popular, Explore, All

**Thành phần:**

- Menu items (Home, Popular, Explore, All)
- Communities list
- Create Community button
- Manage Communities (nếu là mod)

### C3. Component Post Card

Xuất hiện ở: Home, Community, Profile, Popular...

**Thành phần:**

- Author info (avatar + username + community + time)
- Title (h2)
- Content preview
- Media (images/video/poll)
- Interaction bar (upvote, comment, share, save, menu)

### C4. Component Comment

Xuất hiện ở: PostDetail

**Thành phần:**

- Author info
- Content
- Upvote/downvote
- Reply button
- Nested replies (recursive)

### C5. Component Modal

**Loại modals:**

- Create Post Modal
- Report Modal
- Confirm Delete Modal
- Share Modal
- Upload Image Modal

### C6. Component Toast Notification

**Các loại:**

- Success (xanh): "Đã lưu thay đổi"
- Error (đỏ): "Có lỗi xảy ra"
- Warning (vàng): "Vui lòng đăng nhập"
- Info (xanh nhạt): "Email đã được gửi"

---

**Tổng số màn hình: 17**
**Ngày tạo: 20/01/2026**
**Framework: Svelte 5 + TypeScript**
