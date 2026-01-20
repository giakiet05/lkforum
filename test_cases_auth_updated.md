# Test Cases - Module Authentication (Cập nhật)

## A. Đăng ký tài khoản (Registration)

| STT | Tên Test Case            | Mô tả                                     | Điều kiện tiên quyết      | Các bước thực hiện                                                        | Kết quả mong đợi                      | Mức độ ưu tiên |
| --- | ------------------------ | ----------------------------------------- | ------------------------- | ------------------------------------------------------------------------- | ------------------------------------- | -------------- |
| 1   | Đăng ký email hợp lệ     | Kiểm tra đăng ký tài khoản mới bằng email | Chưa có tài khoản         | 1. Nhập email hợp lệ<br>2. Nhập OTP<br>3. Nhập username và password       | Tạo tài khoản thành công              | Cao            |
| 2   | Đăng ký email đã tồn tại | Kiểm tra đăng ký email đã sử dụng         | Email đã đăng ký trước    | 1. Nhập email đã tồn tại<br>2. Bấm "Tiếp tục"                             | Hiển thị lỗi "Email đã được sử dụng"  | Cao            |
| 3   | Đăng ký với OTP sai      | Kiểm tra xác thực OTP sai                 | Đã nhận OTP               | 1. Nhập OTP sai<br>2. Bấm "Xác nhận"                                      | Hiển thị lỗi "Mã OTP không chính xác" | Cao            |
| 4   | Đăng ký với OTP hết hạn  | Kiểm tra OTP quá thời gian                | Đã nhận OTP, chờ > 5 phút | 1. Nhập OTP đúng nhưng hết hạn<br>2. Bấm "Xác nhận"                       | Hiển thị lỗi "Mã OTP đã hết hạn"      | Trung bình     |
| 5   | Gửi lại OTP              | Kiểm tra gửi lại OTP khi chưa nhận được   | Đã gửi OTP lần đầu        | 1. Đợi 60s<br>2. Bấm "Gửi lại"                                            | Gửi OTP mới thành công                | Trung bình     |
| 6   | Đăng ký bằng Google      | Kiểm tra đăng ký qua Google OAuth         | Có tài khoản Google       | 1. Bấm "Đăng nhập với Google"<br>2. Chọn tài khoản Google<br>3. Cấp quyền | Đăng ký/đăng nhập thành công          | Cao            |

## B. Đăng nhập (Login)

| STT | Tên Test Case                           | Mô tả                                         | Điều kiện tiên quyết   | Các bước thực hiện                                                       | Kết quả mong đợi                                      | Mức độ ưu tiên |
| --- | --------------------------------------- | --------------------------------------------- | ---------------------- | ------------------------------------------------------------------------ | ----------------------------------------------------- | -------------- |
| 7   | Đăng nhập thành công với email/password | Kiểm tra đăng nhập hợp lệ                     | Đã có tài khoản local  | 1. Nhập username/email<br>2. Nhập password đúng<br>3. Bấm "Đăng nhập"    | Đăng nhập thành công, chuyển trang chủ                | Cao            |
| 8   | Đăng nhập sai mật khẩu                  | Kiểm tra đăng nhập với password sai           | Đã có tài khoản        | 1. Nhập username đúng<br>2. Nhập password sai<br>3. Bấm "Đăng nhập"      | Hiển thị lỗi "Tên đăng nhập hoặc mật khẩu không đúng" | Cao            |
| 9   | Đăng nhập tài khoản không tồn tại       | Kiểm tra đăng nhập với username không tồn tại | Không có tài khoản     | 1. Nhập username không tồn tại<br>2. Nhập password<br>3. Bấm "Đăng nhập" | Hiển thị lỗi "Tên đăng nhập hoặc mật khẩu không đúng" | Cao            |
| 10  | Đăng nhập bằng Google                   | Kiểm tra đăng nhập qua Google                 | Đã có tài khoản Google | 1. Bấm "Đăng nhập với Google"<br>2. Chọn tài khoản                       | Đăng nhập thành công                                  | Cao            |
| 11  | Hiển thị/ẩn mật khẩu                    | Kiểm tra toggle hiển thị mật khẩu             | Đang ở form login      | 1. Nhập password<br>2. Bấm icon "eye"                                    | Password hiển thị dạng text/ẩn                        | Thấp           |

## C. Quên mật khẩu (Forgot Password)

| STT | Tên Test Case                                      | Mô tả                                                             | Điều kiện tiên quyết                      | Các bước thực hiện                                                                          | Kết quả mong đợi                                                                                               | Mức độ ưu tiên |
| --- | -------------------------------------------------- | ----------------------------------------------------------------- | ----------------------------------------- | ------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- | -------------- |
| 12  | Quên mật khẩu - tài khoản local                    | Kiểm tra reset password cho tài khoản local                       | Đã có tài khoản local                     | 1. Bấm "Quên mật khẩu"<br>2. Nhập email<br>3. Nhập OTP<br>4. Nhập mật khẩu mới              | Reset password thành công                                                                                      | Cao            |
| 13  | Quên mật khẩu - email không tồn tại                | Kiểm tra với email chưa đăng ký                                   | Email không tồn tại                       | 1. Bấm "Quên mật khẩu"<br>2. Nhập email không tồn tại<br>3. Bấm "Gửi mã OTP"                | Hiển thị lỗi "Email không tồn tại"                                                                             | Cao            |
| 14  | Quên mật khẩu - tài khoản Google chưa set password | Kiểm tra với tài khoản Google                                     | Đã đăng ký bằng Google, chưa set password | 1. Bấm "Quên mật khẩu"<br>2. Nhập email Google<br>3. Bấm "Gửi mã OTP"                       | Hiển thị lỗi "Tài khoản này đã đăng ký bằng Google. Vui lòng đăng nhập bằng Google hoặc liên hệ quản trị viên" | Cao            |
| 15  | Quên mật khẩu - tài khoản Google đã set password   | Kiểm tra reset password cho tài khoản Google đã liên kết mật khẩu | Đã đăng ký bằng Google và đã set password | 1. Bấm "Quên mật khẩu"<br>2. Nhập email<br>3. Nhập OTP<br>4. Nhập mật khẩu mới              | Reset password thành công                                                                                      | Cao            |
| 16  | Quên mật khẩu - OTP sai                            | Kiểm tra xác thực OTP sai khi reset                               | Đã nhận OTP reset                         | 1. Nhập OTP sai<br>2. Bấm "Xác nhận"                                                        | Hiển thị lỗi "Mã OTP không chính xác"                                                                          | Cao            |
| 17  | Quên mật khẩu - password không khớp                | Kiểm tra xác nhận mật khẩu mới                                    | Đã verify OTP                             | 1. Nhập mật khẩu mới<br>2. Nhập xác nhận khác với mật khẩu mới<br>3. Bấm "Đặt lại mật khẩu" | Hiển thị lỗi "Mật khẩu xác nhận không khớp"                                                                    | Trung bình     |
| 18  | Quên mật khẩu - gửi lại OTP                        | Kiểm tra gửi lại OTP khi reset                                    | Đã gửi OTP reset                          | 1. Đợi 60s<br>2. Bấm "Gửi lại"                                                              | Gửi OTP mới thành công                                                                                         | Trung bình     |

## D. Đặt mật khẩu cho tài khoản Google (Set Password)

| STT | Tên Test Case                               | Mô tả                                      | Điều kiện tiên quyết                     | Các bước thực hiện                                                                                  | Kết quả mong đợi                                              | Mức độ ưu tiên |
| --- | ------------------------------------------- | ------------------------------------------ | ---------------------------------------- | --------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- | -------------- |
| 19  | Hiển thị nút "Đặt mật khẩu"                 | Kiểm tra hiển thị option đặt mật khẩu      | Đăng nhập bằng Google, chưa set password | 1. Vào trang "Cài đặt tài khoản"                                                                    | Hiển thị nút "Đặt mật khẩu" hoặc thông báo khuyến khích       | Cao            |
| 20  | Đặt mật khẩu thành công                     | Kiểm tra đặt mật khẩu cho tài khoản Google | Đăng nhập bằng Google, chưa set password | 1. Bấm "Đặt mật khẩu"<br>2. Nhập mật khẩu mới<br>3. Xác nhận mật khẩu<br>4. Bấm "Lưu"               | Đặt mật khẩu thành công, có thể dùng email/password đăng nhập | Cao            |
| 21  | Đặt mật khẩu - password không khớp          | Kiểm tra xác nhận mật khẩu                 | Đang ở form đặt mật khẩu                 | 1. Nhập mật khẩu mới<br>2. Nhập xác nhận khác với mật khẩu<br>3. Bấm "Lưu"                          | Hiển thị lỗi "Mật khẩu xác nhận không khớp"                   | Trung bình     |
| 22  | Đặt mật khẩu - password quá ngắn            | Kiểm tra validate độ dài mật khẩu          | Đang ở form đặt mật khẩu                 | 1. Nhập mật khẩu < 6 ký tự<br>2. Bấm "Lưu"                                                          | Hiển thị lỗi "Mật khẩu phải có ít nhất 6 ký tự"               | Trung bình     |
| 23  | Đăng nhập bằng email sau khi set password   | Kiểm tra đăng nhập bằng email/password     | Đã set password cho tài khoản Google     | 1. Đăng xuất<br>2. Nhập email<br>3. Nhập password đã set<br>4. Bấm "Đăng nhập"                      | Đăng nhập thành công                                          | Cao            |
| 24  | Đăng nhập bằng Google sau khi set password  | Kiểm tra vẫn đăng nhập được bằng Google    | Đã set password cho tài khoản Google     | 1. Đăng xuất<br>2. Bấm "Đăng nhập với Google"<br>3. Chọn tài khoản                                  | Đăng nhập thành công                                          | Cao            |
| 25  | Không hiển thị nút "Đặt mật khẩu" nếu đã có | Kiểm tra UI khi đã set password            | Đã set password cho tài khoản Google     | 1. Vào trang "Cài đặt tài khoản"                                                                    | Không hiển thị nút "Đặt mật khẩu"                             | Thấp           |
| 26  | Đổi mật khẩu sau khi đã set                 | Kiểm tra đổi mật khẩu                      | Đã set password, đang đăng nhập          | 1. Vào "Đổi mật khẩu"<br>2. Nhập mật khẩu cũ<br>3. Nhập mật khẩu mới<br>4. Xác nhận<br>5. Bấm "Lưu" | Đổi mật khẩu thành công                                       | Trung bình     |

## E. Đăng xuất (Logout)

| STT | Tên Test Case                           | Mô tả                    | Điều kiện tiên quyết | Các bước thực hiện                      | Kết quả mong đợi                            | Mức độ ưu tiên |
| --- | --------------------------------------- | ------------------------ | -------------------- | --------------------------------------- | ------------------------------------------- | -------------- |
| 27  | Đăng xuất thành công                    | Kiểm tra đăng xuất       | Đã đăng nhập         | 1. Bấm "Đăng xuất"                      | Đăng xuất thành công, chuyển về trang login | Cao            |
| 28  | Không truy cập được trang khi đã logout | Kiểm tra phiên đăng nhập | Đã đăng xuất         | 1. Truy cập URL trang yêu cầu đăng nhập | Chuyển về trang login                       | Cao            |

## F. Session & Token Management

| STT | Tên Test Case            | Mô tả                          | Điều kiện tiên quyết               | Các bước thực hiện                      | Kết quả mong đợi                      | Mức độ ưu tiên |
| --- | ------------------------ | ------------------------------ | ---------------------------------- | --------------------------------------- | ------------------------------------- | -------------- |
| 29  | Access token hết hạn     | Kiểm tra refresh token tự động | Đã đăng nhập, access token hết hạn | 1. Gọi API sau khi access token hết hạn | Tự động refresh token, API thành công | Cao            |
| 30  | Refresh token hết hạn    | Kiểm tra yêu cầu đăng nhập lại | Refresh token hết hạn              | 1. Gọi API khi refresh token hết hạn    | Yêu cầu đăng nhập lại                 | Cao            |
| 31  | Lưu trạng thái đăng nhập | Kiểm tra remember me           | Đã đăng nhập                       | 1. Đóng browser<br>2. Mở lại            | Vẫn đăng nhập (nếu remember me)       | Trung bình     |

---

**Tổng số test cases: 31**

**Phân loại theo mức độ ưu tiên:**

- Cao: 22 test cases
- Trung bình: 8 test cases
- Thấp: 1 test case

**Ghi chú:**

- Test cases mới (19-26) tập trung vào tính năng đặt mật khẩu cho tài khoản Google
- Test case 14-15 kiểm tra luồng quên mật khẩu với tài khoản Google (có/chưa set password)
- Test cases đảm bảo cả hai phương thức đăng nhập (Google OAuth và email/password) hoạt động độc lập và kết hợp
