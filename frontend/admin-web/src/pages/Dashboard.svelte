<script lang="ts">
  import { onMount } from "svelte";
  import { logout } from "../services/auth-service";
  import { setAuthenticated } from "../stores/auth-store";
  import { push } from "svelte-spa-router";
  import { getUsers, banUser, unbanUser } from "../services/user-service";
  import {
    getCommunities,
    deleteCommunity,
  } from "../services/community-service";
  import type { UserResponse } from "../dtos/user-dto";
  import type { CommunityResponse } from "../dtos/community-dto";

  let activeTab: "users" | "communities" = $state("users");
  let users = $state<UserResponse[]>([]);
  let communities = $state<CommunityResponse[]>([]);
  let loading = $state(false);

  async function loadUsers() {
    loading = true;
    try {
      const response = await getUsers({ limit: 100 });
      users = response.users;
    } catch (error) {
      console.error("Failed to load users:", error);
    } finally {
      loading = false;
    }
  }

  async function loadCommunities() {
    loading = true;
    try {
      const response = await getCommunities({ limit: 100 });
      communities = response.communities;
    } catch (error) {
      console.error("Failed to load communities:", error);
    } finally {
      loading = false;
    }
  }

  async function handleBanUser(userId: string) {
    if (!confirm("Are you sure you want to ban this user?")) return;
    try {
      await banUser(userId);
      await loadUsers();
    } catch (error) {
      alert("Failed to ban user");
    }
  }

  async function handleUnbanUser(userId: string) {
    try {
      await unbanUser(userId);
      await loadUsers();
    } catch (error) {
      alert("Failed to unban user");
    }
  }

  async function handleDeleteCommunity(communityId: string) {
    if (!confirm("Are you sure you want to delete this community?")) return;
    try {
      await deleteCommunity(communityId);
      await loadCommunities();
    } catch (error) {
      alert("Failed to delete community");
    }
  }

  function handleLogout() {
    logout();
    setAuthenticated(false);
    push("/login");
  }

  onMount(() => {
    // Only load data if we have auth tokens
    // App.svelte will redirect to login if not authenticated
    const timer = setTimeout(() => {
      loadUsers();
      loadCommunities();
    }, 100);

    return () => clearTimeout(timer);
  });
</script>

<div class="dashboard">
  <header>
    <h1>Admin Dashboard</h1>
    <button on:click={handleLogout} class="logout-btn">Logout</button>
  </header>

  <div class="tabs">
    <button
      class:active={activeTab === "users"}
      on:click={() => (activeTab = "users")}
    >
      Users ({users.length})
    </button>
    <button
      class:active={activeTab === "communities"}
      on:click={() => (activeTab = "communities")}
    >
      Communities ({communities.length})
    </button>
  </div>

  {#if loading}
    <div class="loading">Loading...</div>
  {:else if activeTab === "users"}
    <div class="info-message">
      <strong>Note:</strong> User management requires backend admin routes to be
      registered.
      <br />Please ask backend team to register admin routes in
      <code>init.go</code>
    </div>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Username</th>
            <th>Email</th>
            <th>Created At</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user}
            <tr>
              <td>{user.username}</td>
              <td>{user.email}</td>
              <td>{new Date(user.created_at).toLocaleDateString()}</td>
              <td>
                <span class="status" class:banned={user.is_banned}>
                  {user.is_banned ? "Banned" : "Active"}
                </span>
              </td>
              <td>
                {#if user.is_banned}
                  <button
                    class="action-btn unban"
                    on:click={() => handleUnbanUser(user.id)}
                  >
                    Unban
                  </button>
                {:else}
                  <button
                    class="action-btn ban"
                    on:click={() => handleBanUser(user.id)}
                  >
                    Ban
                  </button>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Members</th>
            <th>Created At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each communities as community}
            <tr>
              <td>c/{community.name}</td>
              <td>{community.member_count}</td>
              <td>{new Date(community.created_at).toLocaleDateString()}</td>
              <td>
                <button
                  class="action-btn delete"
                  on:click={() => handleDeleteCommunity(community.id)}
                >
                  Delete
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .dashboard {
    min-height: 100vh;
    background: #f5f5f5;
    padding: 2rem;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  h1 {
    margin: 0;
    color: #333;
  }

  .logout-btn {
    padding: 0.5rem 1rem;
    background: #f44336;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
  }

  .logout-btn:hover {
    background: #da190b;
  }

  .tabs {
    display: flex;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .tabs button {
    padding: 0.75rem 1.5rem;
    background: white;
    border: 2px solid #ddd;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    color: #666;
  }

  .tabs button.active {
    background: #4caf50;
    border-color: #4caf50;
    color: white;
  }

  .loading {
    text-align: center;
    padding: 3rem;
    color: #666;
  }

  .table-container {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid #eee;
  }

  th {
    background: #f9f9f9;
    font-weight: 600;
    color: #555;
  }

  tbody tr:hover {
    background: #fafafa;
  }

  .status {
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 13px;
    font-weight: 500;
    background: #e8f5e9;
    color: #2e7d32;
  }

  .status.banned {
    background: #ffebee;
    color: #c62828;
  }

  .action-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 500;
  }

  .action-btn.ban {
    background: #ff9800;
    color: white;
  }

  .action-btn.ban:hover {
    background: #f57c00;
  }

  .action-btn.unban {
    background: #4caf50;
    color: white;
  }

  .action-btn.unban:hover {
    background: #45a049;
  }

  .action-btn.delete {
    background: #f44336;
    color: white;
  }

  .action-btn.delete:hover {
    background: #da190b;
  }

  .info-message {
    background: #fff3cd;
    border: 1px solid #ffc107;
    color: #856404;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    font-size: 14px;
  }

  .info-message code {
    background: #fff;
    padding: 2px 6px;
    border-radius: 3px;
    font-family: monospace;
  }
</style>
