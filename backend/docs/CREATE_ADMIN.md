# Hướng dẫn tạo Admin User

## Cách 1: Sử dụng script Go (Đã có MongoDB)

Chạy lệnh sau từ thư mục `backend`:

```bash
go run cmd/seed_admin/main.go
```

Script sẽ tạo admin user với thông tin:

- **Email**: admin@lkforum.com
- **Username**: admin
- **Password**: admin123

## Cách 2: Thêm admin user trực tiếp vào MongoDB

### Bước 1: Cài đặt MongoDB

- Download: https://www.mongodb.com/try/download/community
- Hoặc dùng MongoDB Atlas (cloud): https://www.mongodb.com/cloud/atlas/register

### Bước 2: Kết nối MongoDB

Mở MongoDB Compass hoặc mongosh và kết nối đến database

### Bước 3: Chạy query sau

```javascript
use lkforum

db.users.insertOne({
  _id: ObjectId(),
  username: "admin",
  email: "admin@lkforum.com",
  password: "$2a$10$rN8qVqZxGqDWLpLgKZ.YaOGK7YqZ5X5QxWxR6RnZT0pQXGjJVxEqy", // admin123
  provider: "local",
  role: "admin",
  reputation: 0,
  is_verified: true,
  is_banned: false,
  created_at: new Date(),
  role_content: {
    as_admin: {
      permissions: ["all"]
    }
  }
})
```

## Cách 3: Sửa user hiện tại thành admin

Nếu đã có user trong database, update role:

```javascript
db.users.updateOne(
  { email: "your-email@example.com" },
  {
    $set: {
      role: "admin",
      role_content: {
        as_admin: {
          permissions: ["all"],
        },
      },
    },
  }
);
```

## Kiểm tra

Sau khi tạo admin, test login:

```bash
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "admin@lkforum.com",
    "password": "admin123"
  }'
```

Hoặc dùng frontend admin-web tại `http://localhost:5174`
