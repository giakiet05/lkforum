<script lang="ts">
  let activeTab: "account" | "privacy" | "email" | "notification" = "account";

  const user = {
    name: "Long",
    email: "long@example.com",
    avatar: "/avatar.jpg",
    displayName: "Long",
    bio: "A passionate developer sharing thoughts and ideas about technology, coding, and life.",
  };

  function handleSaveAccount() {
    // TODO: Implement save account functionality
    console.log("Save account settings");
  }

  function handleChangePassword() {
    // TODO: Implement change password functionality
    console.log("Change password");
  }

  function handleDeleteAccount() {
    // TODO: Implement delete account functionality
    console.log("Delete account");
  }
</script>

<div class="settings-page">
  <div class="settings-container">
    <div class="settings-header">
      <h1>Settings</h1>
      <p class="settings-description">
        Manage your account settings and preferences
      </p>
    </div>

    <div class="settings-content">
      <div class="settings-sidebar">
        <button
          class="settings-tab"
          class:active={activeTab === "account"}
          on:click={() => (activeTab = "account")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <circle
              cx="10"
              cy="6"
              r="3"
              stroke="currentColor"
              stroke-width="1.5"
            />
            <path
              d="M4 18C4 14.6863 6.68629 12 10 12C13.3137 12 16 14.6863 16 18"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
          Account
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "privacy"}
          on:click={() => (activeTab = "privacy")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 2L4 5V9C4 13 6.5 16.5 10 18C13.5 16.5 16 13 16 9V5L10 2Z"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          Privacy
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "email"}
          on:click={() => (activeTab = "email")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <rect
              x="3"
              y="5"
              width="14"
              height="10"
              rx="2"
              stroke="currentColor"
              stroke-width="1.5"
            />
            <path
              d="M3 7L10 11L17 7"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          Email
        </button>

        <button
          class="settings-tab"
          class:active={activeTab === "notification"}
          on:click={() => (activeTab = "notification")}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M10 3C7.79 3 6 4.79 6 7V10L4 12V13H16V12L14 10V7C14 4.79 12.21 3 10 3Z"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M8.5 16C8.5 16.8284 9.17157 17.5 10 17.5C10.8284 17.5 11.5 16.8284 11.5 16"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
          Notification
        </button>
      </div>

      <div class="settings-main">
        {#if activeTab === "account"}
          <div class="settings-section">
            <h2>Account Settings</h2>
            <p class="section-description">
              Manage your account information and preferences
            </p>

            <div class="form-group">
              <label for="display-name">Display Name</label>
              <input
                type="text"
                id="display-name"
                value={user.displayName}
                placeholder="Enter your display name"
              />
            </div>

            <div class="form-group">
              <label for="username">Username</label>
              <div class="input-with-prefix">
                <span class="input-prefix">u/</span>
                <input
                  type="text"
                  id="username"
                  value={user.name}
                  placeholder="Enter your username"
                />
              </div>
              <p class="input-hint">Your username cannot be changed once set</p>
            </div>

            <div class="form-group">
              <label for="bio">Bio</label>
              <textarea
                id="bio"
                rows="4"
                value={user.bio}
                placeholder="Tell us about yourself"
              ></textarea>
              <p class="input-hint">Brief description for your profile</p>
            </div>

            <div class="form-group">
              <label>Profile Picture</label>
              <div class="avatar-upload">
                <div class="avatar-preview">
                  <img src={user.avatar} alt="Profile" />
                </div>
                <div class="avatar-actions">
                  <button class="btn-secondary">Change Avatar</button>
                  <button class="btn-text">Remove</button>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button class="btn-primary" on:click={handleSaveAccount}>
                Save Changes
              </button>
              <button class="btn-secondary">Cancel</button>
            </div>

            <div class="password-section">
              <h3>Password</h3>
              <button class="btn-secondary" on:click={handleChangePassword}>
                Change Password
              </button>
            </div>
          </div>
        {:else if activeTab === "privacy"}
          <div class="settings-section">
            <h2>Privacy Settings</h2>
            <p class="section-description">
              Privacy settings will be implemented here
            </p>
          </div>
        {:else if activeTab === "email"}
          <div class="settings-section">
            <h2>Email Settings</h2>
            <p class="section-description">
              Email settings will be implemented here
            </p>
          </div>
        {:else if activeTab === "notification"}
          <div class="settings-section">
            <h2>Notification Settings</h2>
            <p class="section-description">
              Notification settings will be implemented here
            </p>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .settings-page {
    background-color: white;
    min-height: 100vh;
    padding-top: 72px;
  }

  .settings-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1.5rem 2rem;
  }

  .settings-header {
    margin-bottom: 2rem;
  }

  .settings-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: #1a1a1b;
    margin: 0 0 0.5rem 0;
    font-family: "Roboto", sans-serif;
  }

  .settings-description {
    color: #7c7c7c;
    margin: 0;
    font-size: 0.95rem;
  }

  .settings-content {
    display: flex;
    gap: 2rem;
  }

  .settings-sidebar {
    width: 240px;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .settings-tab {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: none;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    color: #7c7c7c;
    font-size: 0.95rem;
    font-weight: 500;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s ease;
    text-align: left;
  }

  .settings-tab:hover {
    background-color: #f6f7f8;
    color: #1a1a1b;
  }

  .settings-tab.active {
    background-color: #f0f1f2;
    color: #153060;
    font-weight: 600;
  }

  .settings-tab svg {
    flex-shrink: 0;
  }

  .settings-main {
    flex: 1;
    background-color: white;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    padding: 2rem;
  }

  .settings-section {
    margin-bottom: 2rem;
  }

  .settings-section h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 0.5rem 0;
    font-family: "Roboto", sans-serif;
  }

  .section-description {
    color: #7c7c7c;
    margin: 0 0 2rem 0;
    font-size: 0.9rem;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #1a1a1b;
    font-size: 0.9rem;
  }

  .form-group input[type="text"],
  .form-group textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    font-size: 0.95rem;
    color: #1a1a1b;
    font-family: "Roboto", sans-serif;
    transition: all 0.2s ease;
  }

  .form-group input[type="text"]:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #153060;
    box-shadow: 0 0 0 3px rgba(21, 48, 96, 0.1);
  }

  .input-with-prefix {
    display: flex;
    align-items: center;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.2s ease;
  }

  .input-with-prefix:focus-within {
    border-color: #153060;
    box-shadow: 0 0 0 3px rgba(21, 48, 96, 0.1);
  }

  .input-prefix {
    padding: 0.75rem;
    background-color: #f6f7f8;
    color: #7c7c7c;
    font-weight: 500;
    border-right: 1px solid #e6e6e6;
  }

  .input-with-prefix input {
    flex: 1;
    border: none;
    padding: 0.75rem;
  }

  .input-with-prefix input:focus {
    outline: none;
    box-shadow: none;
  }

  .input-hint {
    margin: 0.5rem 0 0 0;
    font-size: 0.85rem;
    color: #7c7c7c;
  }

  .avatar-upload {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .avatar-preview {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    overflow: hidden;
    border: 3px solid #e6e6e6;
  }

  .avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-actions {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-actions {
    display: flex;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 2rem;
    border-top: 1px solid #e6e6e6;
  }

  .btn-primary {
    padding: 0.75rem 2rem;
    background-color: #153060;
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
  }

  .btn-primary:hover {
    background-color: #0d2144;
  }

  .btn-secondary {
    padding: 0.75rem 2rem;
    background-color: #f6f7f8;
    color: #1a1a1b;
    border: none;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    font-size: 0.95rem;
    transition: all 0.2s ease;
  }

  .btn-secondary:hover {
    background-color: #e9e9e9;
  }

  .btn-text {
    padding: 0.5rem 1rem;
    background: none;
    color: #7c7c7c;
    border: none;
    font-weight: 500;
    cursor: pointer;
    font-size: 0.9rem;
    transition: color 0.2s ease;
  }

  .btn-text:hover {
    color: #1a1a1b;
  }

  .password-section {
    margin-top: 3rem;
    padding-top: 2rem;
    border-top: 1px solid #e6e6e6;
  }

  .password-section h3 {
    font-size: 1rem;
    font-weight: 600;
    color: #1a1a1b;
    margin: 0 0 1rem 0;
  }

  .password-section .btn-secondary {
    font-weight: 500;
  }

  @media (max-width: 768px) {
    .settings-content {
      flex-direction: column;
    }

    .settings-sidebar {
      width: 100%;
      flex-direction: row;
      overflow-x: auto;
    }

    .settings-tab {
      white-space: nowrap;
    }
  }
</style>
