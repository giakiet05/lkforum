package service

import (
	"context"
	"errors"
	"testing"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/email"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo" // Mày phải import cái repo thật
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// --- MOCK USER REPO (NÂNG CẤP) ---
type mockUserRepo struct {
	repo.UserRepo     // Nhúng interface thật
	GetByEmailFunc    func(ctx context.Context, email string) (*model.User, error)
	GetByUsernameFunc func(ctx context.Context, username string) (*model.User, error)
	// Thêm hàm Create
	CreateFunc func(ctx context.Context, user *model.User) (*model.User, error)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, errors.New("GetByEmailFunc not implemented")
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.GetByUsernameFunc != nil {
		return m.GetByUsernameFunc(ctx, username)
	}
	return nil, errors.New("GetByUsernameFunc not implemented")
}

// Implement hàm Create
func (m *mockUserRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil, errors.New("CreateFunc not implemented")
}

// --- MOCK EMAIL SENDER ---
// Cần cái này để NewUserService không bị lỗi
type mockEmailSender struct {
	email.Sender              // Nhúng interface
	SendVerificationEmailFunc func(to, otp string) error
}

func (m *mockEmailSender) SendVerificationEmail(to, otp string) error {
	if m.SendVerificationEmailFunc != nil {
		return m.SendVerificationEmailFunc(to, otp)
	}
	// Mặc định là thành công, không làm gì cả
	// để cái goroutine nó chạy mà không bị panic
	return nil
}

func TestUserService_Login(t *testing.T) {
	// --- SETUP CHUNG ---
	correctPassword := "Hello@123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	require.NoError(t, err) // Dùng require để dừng test nếu setup fail

	mockVerifiedUser := &model.User{
		ID:         primitive.NewObjectID(),
		Username:   "kiet",
		Email:      "hello@gmail.com",
		Password:   string(hashedPassword),
		IsVerified: true, // User này đã xác thực
		Role:       model.UserRole,
	}

	mockUnverifiedUser := &model.User{
		ID:         primitive.NewObjectID(),
		Username:   "kiet",                 // Cùng username "kiet"
		Email:      "hello@gmail.com",      // Cùng email
		Password:   string(hashedPassword), // Cùng pass
		IsVerified: false,                  // NHƯNG CHƯA VERIFY
		Role:       model.UserRole,
	}

	// --- BẢNG TEST CASES TỐI ƯU (6 CASES) ---
	testCases := []struct {
		name         string
		identifier   string
		password     string
		mockSetup    func(*mockUserRepo)
		expectErr    error
		expectUser   bool
		expectTokens bool
	}{
		{
			name:       "TC01 - Success (Login bang Username)",
			identifier: "kiet",
			password:   correctPassword,
			mockSetup: func(m *mockUserRepo) {
				m.GetByUsernameFunc = func(ctx context.Context, username string) (*model.User, error) {
					return mockVerifiedUser, nil
				}
			},
			expectErr:    nil,
			expectUser:   true,
			expectTokens: true,
		},
		{
			name:       "TC02 - Success (Login bang Email)",
			identifier: "hello@gmail.com",
			password:   correctPassword,
			mockSetup: func(m *mockUserRepo) {
				m.GetByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
					return mockVerifiedUser, nil
				}
			},
			expectErr:    nil,
			expectUser:   true,
			expectTokens: true,
		},
		{
			name:       "TC03 - Fail (Sai mat khau)",
			identifier: "kiet",
			password:   "wrongpass",
			mockSetup: func(m *mockUserRepo) {
				// Vẫn tìm thấy user, nhưng bcrypt sẽ fail
				m.GetByUsernameFunc = func(ctx context.Context, username string) (*model.User, error) {
					return mockVerifiedUser, nil
				}
			},
			expectErr:    apperror.ErrInvalidCredentials,
			expectUser:   false,
			expectTokens: false,
		},
		{
			name:       "TC04 - Fail (User khong ton tai)",
			identifier: "unknown",
			password:   "123", // Password gì cũng được, vì fail từ trước
			mockSetup: func(m *mockUserRepo) {
				m.GetByUsernameFunc = func(ctx context.Context, username string) (*model.User, error) {
					return nil, mongo.ErrNoDocuments
				}
			},
			expectErr:    apperror.ErrInvalidCredentials,
			expectUser:   false,
			expectTokens: false,
		},
		{
			name:       "TC05 - Fail (User chua verify email)",
			identifier: "kiet",
			password:   correctPassword,
			mockSetup: func(m *mockUserRepo) {
				// Tìm thấy user, pass đúng, NHƯNG user đó chưa verify
				m.GetByUsernameFunc = func(ctx context.Context, username string) (*model.User, error) {
					return mockUnverifiedUser, nil
				}
			},
			expectErr:    apperror.ErrEmailNotVerified,
			expectUser:   false,
			expectTokens: false,
		},
		{
			name:       "TC06 - Fail (Loi database bat ngo)",
			identifier: "kiet",
			password:   correctPassword,
			mockSetup: func(m *mockUserRepo) {
				// Giả lập lỗi DB sập
				m.GetByUsernameFunc = func(ctx context.Context, username string) (*model.User, error) {
					return nil, errors.New("database is down")
				}
			},
			expectErr:    errors.New("database is down"),
			expectUser:   false,
			expectTokens: false,
		},
	}

	// --- CHẠY TEST ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup mock
			mockRepo := &mockUserRepo{}
			tc.mockSetup(mockRepo)

			// 2. Tạo service với mock repo
			s := NewUserService(mockRepo, nil) // emailSender là nil vì Login() ko dùng

			// 3. Gọi hàm cần test
			user, accessToken, refreshToken, err := s.Login(tc.identifier, tc.password)

			// 4. Khẳng định (Assert) kết quả
			if tc.expectErr != nil {
				// Nếu mong đợi có lỗi
				assert.Error(t, err) // Khẳng định là có lỗi

				// Nếu là lỗi động (như DB error), so sánh message
				// Nếu là lỗi tĩnh (như apperror), so sánh thẳng
				if tc.name == "TC06 - Fail (Loi database bat ngo)" {
					assert.EqualError(t, err, tc.expectErr.Error())
				} else {
					assert.Equal(t, tc.expectErr, err) // Khẳng định là ĐÚNG lỗi đó
				}
			} else {
				// Nếu mong đợi thành công
				assert.NoError(t, err) // Khẳng định là không có lỗi
			}

			if tc.expectUser {
				assert.NotNil(t, user) // Khẳng định là có trả về user
			} else {
				assert.Nil(t, user) // Khẳng định là user trả về là nil
			}

			if tc.expectTokens {
				assert.NotEmpty(t, accessToken) // Khẳng định token không rỗng
				assert.NotEmpty(t, refreshToken)
			} else {
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
			}
		})
	}
}

// --- TEST REGISTER (Code mới) ---
func TestUserService_RegisterUser(t *testing.T) {
	// --- SETUP CHUNG ---
	// Giả lập user đã tồn tại (dựa trên precondition của mày)
	preconditionUser := &model.User{
		ID:       primitive.NewObjectID(),
		Username: "kiet",
		Email:    "hello@gmail.com",
	}

	// --- BẢNG TEST CASES TỐI ƯU (4 CASES) ---
	testCases := []struct {
		name       string
		username   string
		email      string
		password   string
		mockSetup  func(*mockUserRepo)
		expectUser bool // Chỉ cần check nil/not nil
		expectErr  error
	}{
		{
			name:     "TC01 - Success (Dang ky thanh cong)",
			username: "long",           // User mới
			email:    "long@gmail.com", // Email mới
			password: "Long@123",
			mockSetup: func(m *mockUserRepo) {
				// 1. Check username -> không thấy
				m.GetByUsernameFunc = func(ctx context.Context, s string) (*model.User, error) {
					return nil, nil
				}
				// 2. Check email -> không thấy
				m.GetByEmailFunc = func(ctx context.Context, s string) (*model.User, error) {
					return nil, nil
				}
				// 3. Tạo user -> thành công
				m.CreateFunc = func(ctx context.Context, u *model.User) (*model.User, error) {
					u.ID = primitive.NewObjectID() // Giả lập DB gán ID
					return u, nil
				}
			},
			expectUser: true,
			expectErr:  nil,
		},
		{
			name:     "TC02 - Fail (Username da ton tai)",
			username: "kiet",           // Dùng username đã tồn tại
			email:    "long@gmail.com", // Email mới
			password: "Long@123",
			mockSetup: func(m *mockUserRepo) {
				// 1. Check username -> THẤY
				m.GetByUsernameFunc = func(ctx context.Context, s string) (*model.User, error) {
					return preconditionUser, nil
				}
				// (Hàm sẽ return ngay, không cần mock GetByEmail hay Create)
			},
			expectUser: false,
			expectErr:  apperror.ErrUsernameExists,
		},
		{
			name:     "TC03 - Fail (Email da ton tai)",
			username: "long",            // User mới
			email:    "hello@gmail.com", // Dùng email đã tồn tại
			password: "Long@123",
			mockSetup: func(m *mockUserRepo) {
				// 1. Check username -> không thấy
				m.GetByUsernameFunc = func(ctx context.Context, s string) (*model.User, error) {
					return nil, nil
				}
				// 2. Check email -> THẤY
				m.GetByEmailFunc = func(ctx context.Context, s string) (*model.User, error) {
					return preconditionUser, nil
				}
				// (Hàm sẽ return ngay, không cần mock Create)
			},
			expectUser: false,
			expectErr:  apperror.ErrEmailExists,
		},
		{
			name:     "TC04 - Fail (Loi DB khi Create)",
			username: "long",           // User mới
			email:    "long@gmail.com", // Email mới
			password: "Long@123",
			mockSetup: func(m *mockUserRepo) {
				// 1. Check username -> không thấy
				m.GetByUsernameFunc = func(ctx context.Context, s string) (*model.User, error) {
					return nil, nil
				}
				// 2. Check email -> không thấy
				m.GetByEmailFunc = func(ctx context.Context, s string) (*model.User, error) {
					return nil, nil
				}
				// 3. Tạo user -> THẤT BẠI
				m.CreateFunc = func(ctx context.Context, u *model.User) (*model.User, error) {
					return nil, errors.New("DB down")
				}
			},
			expectUser: false,
			expectErr:  errors.New("DB down"),
		},
	}

	// --- CHẠY TEST ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Setup mock repo
			mockRepo := &mockUserRepo{}
			tc.mockSetup(mockRepo)

			// 2. Setup mock email sender
			mockSender := &mockEmailSender{}

			// 3. Tạo service
			s := NewUserService(mockRepo, mockSender)

			// 4. Gọi hàm cần test
			user, err := s.RegisterUser(tc.username, tc.email, tc.password)

			// 5. Khẳng định (Assert) kết quả
			if tc.expectErr != nil {
				// Nếu mong đợi có lỗi
				assert.Error(t, err)
				if tc.name == "TC04 - Fail (Loi DB khi Create)" {
					assert.EqualError(t, err, tc.expectErr.Error())
				} else {
					assert.Equal(t, tc.expectErr, err)
				}
			} else {
				// Nếu mong đợi thành công
				assert.NoError(t, err)
			}

			if tc.expectUser {
				assert.NotNil(t, user)
				assert.Equal(t, tc.username, user.Username) // Check kỹ hơn
				assert.Equal(t, tc.email, user.Email)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}
