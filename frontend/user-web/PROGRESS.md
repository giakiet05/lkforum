# Frontend Development Progress

## ✅ ĐÃ HOÀN THÀNH

### Profile Feature

- ✅ User service (13 functions): getMyProfile, getUserByUsername, updateProfile, changePassword, uploadAvatar/Cover, deleteAvatar/Cover, getSettings, updateSettings, getProvinces/Interests/Genders
- ✅ Profile.svelte: Load profile từ API, upload avatar/cover, display stats
- ✅ Tabs: Posts (load từ API), Comments (load từ API), Upvoted (coming soon), Downvoted (coming soon)
- ✅ Error handling với retry button
- ✅ Loading states
- ✅ Null-safe rendering

### Messaging Feature

- ✅ Channel service (6 functions): createChannel, getChannelById, getChannelsByUser, getChannelBetweenUsers, updateChannel, deleteChannel
- ✅ Message service (3 functions): getMessageById, getMessages, deleteMessage
- ✅ WebSocket service: Singleton với auto-reconnect, sendMessage, event listeners
- ✅ Chat store: Global state management với derived stores
- ✅ Messages.svelte: Full integration với real API
- ✅ ChatPopup.svelte: Full integration với real API
- ⚠️ **BLOCKED**: WebSocket connection (Backend middleware cần fix - token từ query param)

### Post Feature (Mới)

- ✅ Post service (5 functions): getPosts, getPostsByUserId, createPost, voteOnPost, uploadPostImages
- ✅ Comment service (2 functions): getCommentsFilter, getCommentsByUserId
- ✅ DTOs: CreatePostRequest, CreatePollRequest, PostResponse, CommentResponse
- ✅ Mock data ở Home.svelte để xem UI (4 posts: text, image, poll, video)

### UI Components

- ✅ CreatePostModal.svelte: UI hoàn chỉnh (Text, Images, Link tabs, community selector, drafts)
- ✅ DraftsModal.svelte: UI quản lý drafts
- ✅ Post.svelte: Component hiển thị post với poll support
- ✅ CommentSection.svelte: Component comments

---

## ❌ CHƯA HOÀN THÀNH / CẦN LÀM

### Profile Feature

- ❌ **Backend lỗi**: `/api/posts?author_id=xxx` và `/api/comments/filter?user_id=xxx` trả về 500 error khi user chưa có posts/comments (backend cần fix để trả về array rỗng thay vì crash)
- ❌ **Backend thiếu API**: Upvoted posts (`/api/users/me/upvoted`)
- ❌ **Backend thiếu API**: Downvoted posts (`/api/users/me/downvoted`)
- ❌ Settings.svelte chưa refactor (vẫn dùng mock data)

### Post Creation

- ❌ **CreatePostModal chưa tích hợp API**:
  - Vẫn dùng `mockJoinedCommunities` thay vì gọi API lấy danh sách communities
  - Vẫn dùng `mockDraftsDetails` thay vì API drafts
  - `handlePost()` chỉ console.log, chưa gọi `createPost()` service
- ❌ **Image/Video upload workflow**: Backend yêu cầu 2 bước (tạo post text → upload images riêng), cần implement logic này
- ❌ **Community API thiếu**: Cần API `/api/communities/joined` hoặc `/api/users/me/memberships` để lấy danh sách communities user đã join
- ❌ **Drafts API**: Chưa có backend API cho drafts (hoặc chưa biết)

### Messaging Feature

- ⚠️ **WebSocket blocked**: Backend middleware `RequireAuth()` chỉ đọc token từ Authorization header, không hỗ trợ query param. Browser WebSocket API không thể gửi custom headers → Connection failed Error 1006
  - **Cần backend dev fix**: Middleware phải đọc token từ query param hoặc tạo WebSocketAuth middleware riêng

### Other Features

- ❌ Register.svelte: Có 6 critical bugs (wrong API endpoints, missing verification flow)
- ❌ Settings.svelte: Chưa tích hợp API
- ❌ Vote system: Chưa implement (có service function `voteOnPost()` nhưng chưa tích hợp vào UI)
- ❌ Save post: Chưa có API và service

---

## 🔧 CẦN BACKEND FIX NGAY

### Cao

1. **WebSocket middleware**: Hỗ trợ token từ query param cho WebSocket connections
2. **GetPosts với filter author_id**: Fix lỗi 500 khi user chưa có posts → trả về `[]`
3. **GetComments với filter user_id**: Fix lỗi 500 khi user chưa có comments → trả về `[]`

### Trung bình

4. API lấy upvoted posts: `GET /api/users/me/upvoted`
5. API lấy downvoted posts: `GET /api/users/me/downvoted`
6. API lấy joined communities: `GET /api/communities/joined` hoặc `/api/users/me/memberships`

### Thấp (có thể làm sau)

7. Drafts API (nếu cần persist drafts)
8. Saved posts API (save/unsave posts)

---

## 📋 NEXT STEPS (Theo thứ tự ưu tiên)

### 1. Đợi Backend Fix (BLOCKING)

- WebSocket middleware
- GetPosts/GetComments 500 error

### 2. Refactor CreatePostModal (Sau khi backend fix)

- Tích hợp `createPost()` service
- Gọi API lấy joined communities (thay mock)
- Implement 2-step upload cho images (create post → upload images)
- Handle errors và loading states

### 3. Refactor Settings.svelte

- Load settings từ API
- Update settings
- Change password flow
- Delete account confirmation

### 4. Implement Vote System

- Integrate `voteOnPost()` vào Post component
- Update UI khi vote (optimistic updates)
- Handle vote errors

### 5. Fix Register.svelte bugs

- Update API endpoints
- Implement verification token flow
- Fix mock data issues

---

## 📝 NOTES

- **Mock data**: Đang giữ mock data ở Home.svelte để test UI trong khi đợi backend fix
- **Service pattern**: Tất cả API calls đi qua service layer, không gọi trực tiếp trong components
- **Error handling**: Đã implement error states với retry buttons
- **TypeScript**: Strict mode, sử dụng DTOs chính xác theo backend
- **Svelte 5**: Sử dụng runes ($state, $derived, $effect)

---

**Last updated**: November 12, 2025
