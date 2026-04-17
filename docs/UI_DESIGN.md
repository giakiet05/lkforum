TÀI LIỆU THIẾT KẾ GIAO DIỆN - LKFORUM

MỤC LỤC

1. Màn hình Trang chủ (Home)
2. Màn hình Đăng nhập (Login)
3. Màn hình Đăng ký (Register)
4. Màn hình Quên mật khẩu (Forgot Password)
5. Màn hình Xác thực Google (Google Setup)
6. Màn hình Chi tiết bài đăng (Post Detail)
7. Màn hình Chỉnh sửa bài đăng (Edit Post)
8. Màn hình Cộng đồng (Community)
9. Màn hình Tạo cộng đồng (Create Community)
10. Màn hình Cài đặt cộng đồng (Community Settings)
11. Màn hình Quản lý cộng đồng (Manage Communities)
12. Màn hình Công cụ quản trị (Mod Tools)
13. Màn hình Trang cá nhân (Profile)
14. Màn hình Cài đặt tài khoản (Settings)
15. Màn hình Tin nhắn (Messages)
16. Màn hình Khám phá (Explore)
17. Màn hình Phổ biến (Popular)

---

1. Màn hình Trang chủ (Home)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Trang chủ là trang đích chính khi người dùng truy cập website diễn đàn. Giao diện hiển thị danh sách bài đăng mới nhất từ các cộng đồng mà người dùng đã tham gia, kèm theo các thông tin về tác giả, cộng đồng, thời gian đăng, số lượt upvote/downvote, và số bình luận. Người dùng có thể tương tác với bài đăng thông qua các nút upvote, downvote, comment, share, save và report. Giao diện bao gồm thanh điều hướng trên cùng với logo, thanh tìm kiếm, icon thông báo và avatar người dùng; thanh menu bên trái với các mục Home, Popular, Explore và danh sách cộng đồng đã tham gia; vùng nội dung chính hiển thị feed bài đăng với các tab lọc (Latest, Hot, Top); và sidebar bên phải hiển thị cộng đồng đề xuất và trending topics.

---

2. Màn hình Đăng nhập (Login)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Đăng nhập cho phép người dùng xác thực tài khoản trước khi truy cập các chức năng của diễn đàn. Hỗ trợ 2 phương thức đăng nhập: đăng nhập bằng Email/Mật khẩu (local) và đăng nhập bằng tài khoản Google OAuth. Giao diện gồm logo LKForum ở trên cùng, form nhập email và mật khẩu với các trường bắt buộc, tùy chọn "Ghi nhớ đăng nhập", liên kết "Quên mật khẩu", nút đăng nhập chính, đường phân cách "Hoặc", nút "Đăng nhập bằng Google", và liên kết "Đăng ký ngay" ở cuối form. Hệ thống validate dữ liệu realtime và hiển thị thông báo lỗi phù hợp khi đăng nhập thất bại hoặc tài khoản bị khóa.

---

3. Màn hình Đăng ký (Register)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Đăng ký cho phép người dùng tạo tài khoản mới bằng email qua quy trình 3 bước. Bước 1: Nhập email và gửi mã OTP xác thực đến email đó. Bước 2: Nhập mã OTP gồm 6 chữ số để xác thực email, với countdown 5 phút và tùy chọn gửi lại mã khi hết hạn. Bước 3: Tạo username (3-20 ký tự) và mật khẩu mạnh (tối thiểu 8 ký tự, có chữ hoa, chữ thường, số và ký tự đặc biệt), xác nhận mật khẩu, sau đó hoàn tất đăng ký. Ngoài ra hỗ trợ đăng ký nhanh bằng Google OAuth. Sau khi đăng ký thành công, người dùng được tự động đăng nhập và chuyển về trang chủ.

---

4. Màn hình Quên mật khẩu (Forgot Password)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Quên mật khẩu cho phép người dùng đặt lại mật khẩu khi quên mật khẩu cũ thông qua quy trình 3 bước. Bước 1: Nhập email đã đăng ký để nhận mã OTP xác thực. Bước 2: Xác thực mã OTP gồm 6 chữ số với thời gian hiệu lực 5 phút, có thể gửi lại mã nếu hết hạn. Bước 3: Nhập mật khẩu mới với yêu cầu bảo mật tương tự đăng ký (8+ ký tự, chữ hoa, chữ thường, số, ký tự đặc biệt) và xác nhận lại mật khẩu. Sau khi đặt lại thành công, người dùng được chuyển về trang đăng nhập và các phiên đăng nhập cũ sẽ tự động logout.

---

5. Màn hình Xác thực Google (Google Setup)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Xác thực Google xuất hiện khi người dùng đăng nhập lần đầu bằng Google OAuth. Màn hình hiển thị avatar và email từ tài khoản Google (read-only), yêu cầu người dùng chọn một username duy nhất để hoàn tất việc tạo tài khoản. Username phải từ 3-20 ký tự, chỉ chứa chữ, số và dấu gạch dưới, chưa được sử dụng bởi người khác. Hệ thống kiểm tra tính khả dụng của username realtime và hiển thị icon xác nhận. Sau khi hoàn tất, tài khoản được tạo với provider="google" và người dùng tự động đăng nhập vào hệ thống.

---

6. Màn hình Chi tiết bài đăng (Post Detail)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Chi tiết bài đăng hiển thị toàn bộ nội dung của một bài đăng bao gồm tiêu đề, nội dung văn bản đầy đủ, hình ảnh/video/poll (nếu có), thông tin tác giả và cộng đồng. Người dùng có thể tương tác thông qua các nút upvote/downvote, comment, share, save và report. Phần bình luận hiển thị dạng nested với khả năng reply, edit, delete và vote cho từng comment. Bài đăng có thể được chỉnh sửa hoặc xóa bởi tác giả, moderator hoặc admin. Sidebar bên phải hiển thị thông tin cộng đồng với nút join/leave và danh sách quy tắc. URL của bài đăng sử dụng slug-based format dạng /post/title-slug-{id} để thân thiện với SEO.

---

7. Màn hình Chỉnh sửa bài đăng (Edit Post)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Chỉnh sửa bài đăng cho phép tác giả hoặc moderator chỉnh sửa tiêu đề và nội dung của bài đăng. Có thể thay đổi văn bản, thêm/xóa hình ảnh, thêm/xóa tags nhưng không thể chỉnh sửa poll sau khi có người vote. Form hiển thị breadcrumb để người dùng biết vị trí hiện tại. Có nút "Hủy" để quay lại và nút "Lưu thay đổi" để cập nhật, với tracking changes để enable/disable nút Save. Sau khi lưu thành công, bài đăng được đánh dấu is_edited = true và người dùng được chuyển về trang chi tiết bài đăng.

---

8. Màn hình Cộng đồng (Community)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Cộng đồng hiển thị thông tin chi tiết về một cộng đồng bao gồm banner, avatar, tên (dạng lk/tên-cộng-đồng), mô tả, và thống kê số thành viên, số bài đăng. Người dùng có thể tham gia/rời cộng đồng bằng nút Join/Joined, tạo bài mới, và truy cập công cụ quản trị nếu là moderator. Giao diện có 4 tabs: Posts hiển thị danh sách bài đăng với filter (Hot, New, Top), About hiển thị mô tả đầy đủ và thông tin chi tiết, Rules liệt kê các quy tắc, và Moderators hiển thị danh sách người quản trị. Cộng đồng có thể là Public/Private và có thể có badge 18+ với xác nhận tuổi khi truy cập.

---

9. Màn hình Tạo cộng đồng (Create Community)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Tạo cộng đồng cho phép người dùng tạo một cộng đồng mới với các thông tin: tên cộng đồng (3-21 ký tự, chỉ chữ thường, số, gạch ngang, unique, không thể đổi sau khi tạo), mô tả ngắn (max 500 ký tự), avatar (circle, max 5MB, tỷ lệ 1:1), banner (rectangle, max 10MB, tỷ lệ 3:1). Có các tùy chọn cài đặt: Private community (chỉ members xem được), 18+ content, bài đăng cần duyệt, thành viên mới cần duyệt. Hệ thống kiểm tra tên cộng đồng realtime để đảm bảo chưa tồn tại. Sau khi tạo thành công, người dùng tự động trở thành creator và moderator của cộng đồng đó.

---

10. Màn hình Cài đặt cộng đồng (Community Settings)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Cài đặt cộng đồng dành cho creator và moderators để quản lý cộng đồng với 5 tabs. Tab General: chỉnh sửa mô tả, tags. Tab Appearance: thay đổi avatar, banner, theme color. Tab Rules: quản lý danh sách quy tắc với khả năng thêm/sửa/xóa và sắp xếp thứ tự. Tab Moderation: cấu hình kiểm duyệt bài đăng, thành viên, content filters, banned words. Tab Advanced: cài đặt privacy (Public/Private/Restricted), content rating 18+, và danger zone với options archive hoặc delete community. Chỉ creator và moderators có quyền truy cập, thay đổi được lưu bằng nút Save changes.

---

11. Màn hình Quản lý cộng đồng (Manage Communities)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Quản lý cộng đồng hiển thị danh sách các cộng đồng mà người dùng là creator hoặc moderator. Mỗi community card hiển thị avatar, tên, role badge (Creator/Moderator), thống kê số members, posts, và số unmoderated posts nếu có. Có quick actions để truy cập nhanh: View community, Mod Tools, Settings. Nếu người dùng chưa quản lý cộng đồng nào, hiển thị empty state với nút "Tạo cộng đồng mới". Màn hình này giúp người quản trị dễ dàng theo dõi và quản lý nhiều cộng đồng cùng lúc.

---

12. Màn hình Công cụ quản trị (Mod Tools)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Công cụ quản trị cung cấp các công cụ cho moderators quản lý cộng đồng với 5 tabs. Tab Queue: duyệt bài đăng chờ approval với actions Approve/Reject/Remove (nếu bật post_require_approval). Tab Reports: xử lý các báo cáo vi phạm về posts/comments/users với actions View/Dismiss/Remove content/Ban user. Tab Banned Users: quản lý danh sách users bị cấm (banned/muted) với thông tin lý do, thời gian, moderator thực hiện, có thể unban hoặc edit ban. Tab Moderators: quản lý danh sách mods, add/remove moderators (chỉ creator). Tab Moderation Log: xem lịch sử các hành động mod với filters theo moderator, action type, date range. Chỉ creator và moderators có quyền truy cập.

---

13. Màn hình Trang cá nhân (Profile)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Trang cá nhân hiển thị thông tin công khai của người dùng bao gồm cover image, avatar, username, bio, stats (số bài đăng, bình luận, reputation, ngày tham gia), thông tin cá nhân (location, gender, birthday, interests), và social links (website, Facebook, YouTube, GitHub). Có 4 tabs: Posts hiển thị danh sách bài đăng với filter (New/Hot/Top), Comments hiển thị danh sách bình luận, Saved hiển thị bài đăng đã lưu (chỉ khi xem profile của mình), About hiển thị thông tin chi tiết. Nếu là profile của mình có nút "Edit Profile", nếu là profile người khác có nút "Send Message" và "Report User". Profile có thể được set private trong settings.

---

14. Màn hình Cài đặt tài khoản (Settings)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Cài đặt tài khoản cho phép người dùng tùy chỉnh thông tin cá nhân và preferences với 6 tabs. Tab Profile: chỉnh sửa avatar, cover, bio, personal info (gender, date of birth, location), interests, social links. Tab Account: quản lý email, password (chỉ local account), deactivate/delete account. Tab Appearance: chọn theme (Light/Dark/Auto), font size (Small/Medium/Large) với live preview. Tab Notifications: cấu hình thông báo in-app và email cho các sự kiện (comment, mention, upvote, message), email digest frequency. Tab Privacy: cài đặt profile visibility, interactions (allow DM, mentions, follow), content preferences (show NSFW). Tab Blocked Users: quản lý danh sách users đã chặn với option unblock.

---

15. Màn hình Tin nhắn (Messages)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Tin nhắn cho phép người dùng nhắn tin trực tiếp với nhau với giao diện 2 cột. Cột trái hiển thị danh sách conversations với search bar, mỗi item show avatar, username, message preview, timestamp và badge số message chưa đọc. Cột phải là chat box hiển thị messages của conversation đang chọn với header (avatar, username, status online/offline, menu), message area (hiển thị messages align left/right theo người gửi, system messages căn giữa), và input area (textarea với emoji picker, nút Send, Enter để gửi). Có modal "New Message" để tìm user và bắt đầu conversation mới. Hỗ trợ realtime messaging với typing indicator và notification sound.

---

16. Màn hình Khám phá (Explore)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Khám phá giúp người dùng tìm kiếm cộng đồng mới và nội dung nổi bật từ toàn bộ diễn đàn. Giao diện bao gồm search bar để tìm cộng đồng và bài đăng với autocomplete suggestions, card "Trending Topics" hiển thị top 10 topics hot với số bài đăng và growth indicator, card "Cộng đồng đề xuất" hiển thị grid các community cards với nút Join, và Explore Feed với 3 tabs (All/Trending/New) hiển thị bài đăng từ tất cả cộng đồng với infinite scroll. Click vào topic để filter posts, click vào community để xem chi tiết, join community để thêm vào danh sách "My communities".

---

17. Màn hình Phổ biến (Popular)

[Hình ảnh minh họa - thêm sau]

Mô tả chức năng:

Màn hình Phổ biến hiển thị các bài đăng có nhiều tương tác nhất với time filter dropdown cho phép chọn khoảng thời gian (Today/This Week/This Month/All Time). Danh sách bài đăng được sắp xếp theo điểm (votes_count.up - votes_count.down), hiển thị ranking number (#1, #2, #3...) và badge "🔥 Hot" cho top 3. Sidebar hiển thị Top Communities (cộng đồng có nhiều thành viên nhất) và Rising Stars (cộng đồng đang tăng trưởng nhanh). Người dùng có thể thay đổi time filter để xem top posts trong các khoảng thời gian khác nhau, tương tác với bài đăng realtime để cập nhật ranking.

---

Tổng số màn hình: 17
Ngày tạo: 20/01/2026
Framework: Svelte 5 + TypeScript
