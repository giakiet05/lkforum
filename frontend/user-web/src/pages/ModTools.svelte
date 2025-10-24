<script lang="ts">
  import { push } from "svelte-spa-router";
  import { mockQueuePosts } from "../mocks/mod-queue.mock";
  import type { QueuePost } from "../mocks/mod-queue.mock";
  import { mockCommunityRules } from "../mocks/community-rules.mock";
  import type { CommunityRule } from "../mocks/community-rules.mock";
  import { mockRestrictedUsers } from "../mocks/restricted-users.mock";
  import type { RestrictedUser } from "../mocks/restricted-users.mock";
  import { mockModerators, mockApprovedUsers } from "../mocks/moderators.mock";
  import type { Moderator, ApprovedUser } from "../mocks/moderators.mock";

  export interface Props {
    params?: { name?: string };
  }

  let { params }: Props = $props();

  const communityName = params?.name || "";

  type SidebarItem = "queue" | "restricted" | "members" | "rules";
  type QueueTab = "unmoderated" | "edited" | "removed" | "reported";
  type SortOption = "newest" | "oldest" | "most-reported";
  type RestrictedTab = "banned" | "muted";
  type MembersTab = "moderators" | "approved";

  let activeSidebarItem = $state<SidebarItem>("queue");
  let activeQueueTab = $state<QueueTab>("unmoderated");
  let sortBy = $state<SortOption>("newest");

  // Restricted Users state
  let activeRestrictedTab = $state<RestrictedTab>("banned");
  let showBanModal = $state(false);
  let showMuteModal = $state(false);
  let banUsername = $state("");
  let banRule = $state("");
  let banDuration = $state("");
  let banReason = $state("");
  let banNote = $state("");

  // Mod & Members state
  let activeMembersTab = $state<MembersTab>("moderators");
  let showInviteModal = $state(false);
  let inviteUsername = $state("");
  let inviteType = $state<"mod" | "approved">("mod");
  let invitePermission = $state("Everything");
  let inviteCanEdit = $state(true);

  // Rules state
  let showRuleForm = $state(false);
  let editingRuleId = $state<number | null>(null);
  let ruleName = $state("");
  let ruleDescription = $state("");
  let reportReason = $state("");

  const isRuleFormValid = $derived(
    ruleName.trim().length > 0 && ruleDescription.trim().length > 0
  );

  const filteredPosts = $derived(() => {
    let posts = mockQueuePosts.filter((p) => p.queueType === activeQueueTab);

    // Sort posts
    if (sortBy === "newest") {
      posts.sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      );
    } else if (sortBy === "oldest") {
      posts.sort(
        (a, b) =>
          new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      );
    } else if (sortBy === "most-reported") {
      posts.sort((a, b) => (b.reportCount || 0) - (a.reportCount || 0));
    }

    return posts;
  });

  function handleExitModTools() {
    push(`/lk/${communityName}`);
  }

  function handleSidebarClick(item: SidebarItem) {
    activeSidebarItem = item;
  }

  function formatTime(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffHours / 24);

    if (diffHours < 1) return "just now";
    if (diffHours < 24) return `${diffHours} hours ago`;
    if (diffDays < 7) return `${diffDays} days ago`;
    return date.toLocaleDateString();
  }

  function handleApprove(postId: string) {
    console.log("Approve post:", postId);
    // TODO: Implement approve functionality
  }

  function handleRemove(postId: string) {
    console.log("Remove post:", postId);
    // TODO: Implement remove functionality
  }

  function handleSaveRule() {
    if (!isRuleFormValid) return;
    console.log("Save rule:", { ruleName, ruleDescription, reportReason });
    alert("Rule saved!");
    // Reset form and go back to list
    ruleName = "";
    ruleDescription = "";
    reportReason = "";
    showRuleForm = false;
    editingRuleId = null;
  }

  function handleCreateRule() {
    showRuleForm = true;
    editingRuleId = null;
    ruleName = "";
    ruleDescription = "";
    reportReason = "";
  }

  function handleEditRule(ruleId: number) {
    // TODO: Load rule data from mock
    showRuleForm = true;
    editingRuleId = ruleId;
    // Mock data for now
    ruleName = "rule1";
    ruleDescription = "we dont talk about this group";
    reportReason = "";
  }

  function handleDeleteRule(ruleId: number) {
    if (confirm("Are you sure you want to delete this rule?")) {
      console.log("Delete rule:", ruleId);
      alert("Rule deleted!");
    }
  }

  function handleBackToRulesList() {
    showRuleForm = false;
    editingRuleId = null;
    ruleName = "";
    ruleDescription = "";
    reportReason = "";
  }

  // Restricted Users functions
  const filteredRestrictedUsers = $derived(
    mockRestrictedUsers.filter((u) => u.type === activeRestrictedTab)
  );

  function handleOpenBanModal() {
    showBanModal = true;
    banUsername = "";
    banRule = "";
    banDuration = "";
    banReason = "";
    banNote = "";
  }

  function handleOpenMuteModal() {
    showMuteModal = true;
    banUsername = "";
    banRule = "";
    banDuration = "";
    banReason = "";
    banNote = "";
  }

  function handleCloseBanModal() {
    showBanModal = false;
  }

  function handleCloseMuteModal() {
    showMuteModal = false;
  }

  function handleBanUser() {
    console.log("Ban user:", {
      banUsername,
      banRule,
      banDuration,
      banReason,
      banNote,
    });
    alert("User banned!");
    handleCloseBanModal();
  }

  function handleMuteUser() {
    console.log("Mute user:", {
      banUsername,
      banRule,
      banDuration,
      banReason,
      banNote,
    });
    alert("User muted!");
    handleCloseMuteModal();
  }

  // Mod & Members functions
  function handleOpenInviteModal(type: "mod" | "approved") {
    inviteType = type;
    showInviteModal = true;
    inviteUsername = "";
    invitePermission = "Everything";
    inviteCanEdit = true;
  }

  function handleCloseInviteModal() {
    showInviteModal = false;
    inviteUsername = "";
    invitePermission = "Everything";
    inviteCanEdit = true;
  }

  function handleInviteUser() {
    if (!inviteUsername.trim()) {
      alert("Please enter a username!");
      return;
    }
    console.log(`Invite ${inviteType}:`, {
      username: inviteUsername,
      permission: invitePermission,
      canEdit: inviteCanEdit,
    });
    alert(
      `Invitation sent to ${inviteUsername} to become ${inviteType === "mod" ? "moderator" : "approved user"}!`
    );
    handleCloseInviteModal();
  }

  function handleEditMember(id: number, type: "mod" | "approved") {
    console.log("Edit member:", id, type);
    alert("Edit member functionality coming soon!");
  }

  function handleDeleteMember(id: number, type: "mod" | "approved") {
    if (
      confirm(
        `Are you sure you want to remove this ${type === "mod" ? "moderator" : "approved user"}?`
      )
    ) {
      console.log("Delete member:", id, type);
      alert("Member removed!");
    }
  }
</script>

<div class="mod-tools-page">
  <!-- Sidebar -->
  <aside class="mod-sidebar">
    <button class="exit-mod-tools" onclick={handleExitModTools}>
      <!-- Use /arrowback_icon.svg (place file into public folder as arrowback_icon.svg) -->
      <img src="/arrowback_icon.svg" alt="" width="20" height="20" />
      Exit mod tools
    </button>

    <nav class="mod-nav">
      <button
        class="nav-item"
        class:active={activeSidebarItem === "queue"}
        onclick={() => handleSidebarClick("queue")}
      >
        <img src="/queue_icon.svg" alt="" width="20" height="20" />
        <span>Queue</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "restricted"}
        onclick={() => handleSidebarClick("restricted")}
      >
        <img src="/restricted_icon.svg" alt="" width="20" height="20" />
        <span>Restricted Users</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "members"}
        onclick={() => handleSidebarClick("members")}
      >
        <img src="/member_icon.svg" alt="" width="20" height="20" />
        <span>Mods & Members</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "rules"}
        onclick={() => handleSidebarClick("rules")}
      >
        <img src="/rule_icon.svg" alt="" width="20" height="20" />
        <span>Rules</span>
      </button>
    </nav>
  </aside>

  <!-- Main Content -->
  <main class="mod-content">
    {#if activeSidebarItem === "queue"}
      <div class="queue-section">
        <div class="queue-header">
          <h1>Queue</h1>
          <div class="sort-options">
            <select bind:value={sortBy}>
              <option value="newest">Newest</option>
              <option value="oldest">Oldest</option>
              <option value="most-reported">Most Reported</option>
            </select>
          </div>
        </div>

        <!-- Queue Tabs -->
        <div class="queue-tabs">
          <button
            class="tab-btn"
            class:active={activeQueueTab === "unmoderated"}
            onclick={() => (activeQueueTab = "unmoderated")}
          >
            Unmoderated
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "edited"}
            onclick={() => (activeQueueTab = "edited")}
          >
            Edited
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "removed"}
            onclick={() => (activeQueueTab = "removed")}
          >
            Removed
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "reported"}
            onclick={() => (activeQueueTab = "reported")}
          >
            Reported
          </button>
        </div>

        <!-- Posts List -->
        <div class="posts-list">
          {#each filteredPosts() as post (post.id)}
            <div class="post-card">
              <div class="post-header">
                <div class="post-author">
                  {#if post.authorAvatar}
                    <img
                      src={post.authorAvatar}
                      alt={post.author}
                      class="author-avatar"
                    />
                  {:else}
                    <div class="author-avatar-placeholder">
                      {post.author[0].toUpperCase()}
                    </div>
                  {/if}
                  <div class="author-info">
                    <span class="author-name">u/{post.author}</span>
                    <span class="post-time">{formatTime(post.createdAt)}</span>
                  </div>
                </div>
                {#if post.reportCount}
                  <span class="report-badge">{post.reportCount} reports</span>
                {/if}
              </div>

              <h3 class="post-title">{post.title}</h3>
              <p class="post-content">{post.content}</p>

              {#if post.reportReason}
                <div class="report-info">
                  <strong>Report reason:</strong>
                  {post.reportReason}
                </div>
              {/if}

              {#if post.removedReason}
                <div class="removed-info">
                  <strong>Removed by:</strong>
                  {post.removedBy} - {post.removedReason}
                </div>
              {/if}

              <div class="post-actions">
                <button
                  class="action-btn approve"
                  onclick={() => handleApprove(post.id)}
                >
                  Approve
                </button>
                <button
                  class="action-btn remove"
                  onclick={() => handleRemove(post.id)}
                >
                  Remove
                </button>
              </div>
            </div>
          {:else}
            <div class="empty-state">
              <p>No posts in this queue</p>
            </div>
          {/each}
        </div>
      </div>
    {:else if activeSidebarItem === "restricted"}
      <!-- Restricted Users Section -->
      <div class="restricted-section">
        <div class="restricted-header">
          <h1>Restricted Users</h1>
          <button
            class="action-btn-primary"
            onclick={activeRestrictedTab === "banned"
              ? handleOpenBanModal
              : handleOpenMuteModal}
          >
            {activeRestrictedTab === "banned" ? "Ban User" : "Mute User"}
          </button>
        </div>

        <!-- Tabs -->
        <div class="restricted-tabs">
          <button
            class="tab-btn"
            class:active={activeRestrictedTab === "banned"}
            onclick={() => (activeRestrictedTab = "banned")}
          >
            Banned
          </button>
          <button
            class="tab-btn"
            class:active={activeRestrictedTab === "muted"}
            onclick={() => (activeRestrictedTab = "muted")}
          >
            Muted
          </button>
        </div>

        <!-- Table -->
        <div class="restricted-table">
          {#if activeRestrictedTab === "banned"}
            <div class="table-header">
              <div class="col">USERNAME</div>
              <div class="col">DURATION</div>
              <div class="col">DATE</div>
              <div class="col">REASON</div>
              <div class="col">NOTE</div>
            </div>
          {:else}
            <div class="table-header muted">
              <div class="col">USERNAME</div>
              <div class="col">Duration</div>
              <div class="col">NOTE</div>
            </div>
          {/if}

          {#each filteredRestrictedUsers as user}
            <div class="table-row">
              {#if activeRestrictedTab === "banned"}
                <div class="col user-col">
                  <img src={user.avatar} alt="" class="user-avatar" />
                  <span>{user.username}</span>
                </div>
                <div class="col">{user.duration}</div>
                <div class="col">{user.date || "-"}</div>
                <div class="col">{user.reason || "-"}</div>
                <div class="col">{user.note || "-"}</div>
              {:else}
                <div class="col user-col">
                  <img src={user.avatar} alt="" class="user-avatar" />
                  <span>{user.username}</span>
                </div>
                <div class="col">{user.duration}</div>
                <div class="col">{user.note || "-"}</div>
              {/if}
            </div>
          {:else}
            <div class="empty-state">
              <p>
                No {activeRestrictedTab} users
              </p>
            </div>
          {/each}
        </div>
      </div>
    {:else if activeSidebarItem === "members"}
      <!-- Mod & Members Section -->
      <div class="members-section">
        <div class="members-header">
          <h1>Mod & Members</h1>
          <button
            class="action-btn-primary"
            onclick={() =>
              handleOpenInviteModal(
                activeMembersTab === "moderators" ? "mod" : "approved"
              )}
          >
            {activeMembersTab === "moderators"
              ? "Invite Mod"
              : "Add Approved User"}
          </button>
        </div>

        <!-- Tabs -->
        <div class="members-tabs">
          <button
            class="tab-btn"
            class:active={activeMembersTab === "moderators"}
            onclick={() => (activeMembersTab = "moderators")}
          >
            Moderators
          </button>
          <button
            class="tab-btn"
            class:active={activeMembersTab === "approved"}
            onclick={() => (activeMembersTab = "approved")}
          >
            Approved Users
          </button>
        </div>

        <!-- Table -->
        <div class="members-table">
          {#if activeMembersTab === "moderators"}
            <div class="table-header moderators">
              <div class="col">USERNAME</div>
              <div class="col">PERMISSIONS</div>
              <div class="col">You can edit</div>
              <div class="col">JOINED</div>
              <div class="col"></div>
            </div>

            {#each mockModerators as mod}
              <div class="table-row moderators">
                <div class="col user-col">
                  <img src={mod.avatar} alt="" class="user-avatar" />
                  <span>{mod.username}</span>
                </div>
                <div class="col">{mod.permissions}</div>
                <div class="col">{mod.canEdit ? "Yes" : "No"}</div>
                <div class="col joined-col">{mod.joinedDate}</div>
                <div class="col actions-col">
                  <button
                    class="icon-btn edit"
                    onclick={() => handleEditMember(mod.id, "mod")}
                    title="Edit moderator"
                  >
                    <img
                      src="/write_icon.svg"
                      alt="Edit"
                      width="20"
                      height="20"
                    />
                  </button>
                  <button
                    class="icon-btn delete"
                    onclick={() => handleDeleteMember(mod.id, "mod")}
                    title="Remove moderator"
                  >
                    <img
                      src="/bin_icon.svg"
                      alt="Delete"
                      width="20"
                      height="20"
                    />
                  </button>
                </div>
              </div>
            {:else}
              <div class="empty-state">
                <p>No moderators yet</p>
              </div>
            {/each}
          {:else}
            <div class="table-header approved">
              <div class="col">USERNAME</div>
              <div class="col">Date</div>
              <div class="col"></div>
            </div>

            {#each mockApprovedUsers as user}
              <div class="table-row approved">
                <div class="col user-col">
                  <img src={user.avatar} alt="" class="user-avatar" />
                  <span>{user.username}</span>
                </div>
                <div class="col">{user.joinedDate}</div>
                <div class="col actions-col">
                  <button
                    class="icon-btn edit"
                    onclick={() => handleEditMember(user.id, "approved")}
                    title="Edit user"
                  >
                    <img
                      src="/write_icon.svg"
                      alt="Edit"
                      width="20"
                      height="20"
                    />
                  </button>
                  <button
                    class="icon-btn delete"
                    onclick={() => handleDeleteMember(user.id, "approved")}
                    title="Remove user"
                  >
                    <img
                      src="/bin_icon.svg"
                      alt="Delete"
                      width="20"
                      height="20"
                    />
                  </button>
                </div>
              </div>
            {:else}
              <div class="empty-state">
                <p>No approved users yet</p>
              </div>
            {/each}
          {/if}
        </div>
      </div>
    {:else if activeSidebarItem === "rules"}
      {#if !showRuleForm}
        <!-- Rules List View -->
        <div class="rules-list-section">
          <div class="rules-list-header">
            <h1>Community Rules</h1>
            <button class="create-rule-btn" onclick={handleCreateRule}>
              Create Rule
            </button>
          </div>

          <div class="rules-table">
            <div class="rules-table-header">
              <div class="col-name">NAME</div>
              <div class="col-created">CREATED</div>
            </div>

            {#each mockCommunityRules as rule, index}
              <div class="rule-row" onclick={() => handleEditRule(rule.id)}>
                <div class="rule-info">
                  <span class="rule-number">{index + 1}</span>
                  <div class="rule-details">
                    <h3 class="rule-name">{rule.name}</h3>
                    <p class="rule-desc">{rule.description}</p>
                  </div>
                </div>
                <div class="rule-actions">
                  <button
                    class="icon-btn edit"
                    onclick={(e) => {
                      e.stopPropagation();
                      handleEditRule(rule.id);
                    }}
                    title="Edit rule"
                  >
                    <img
                      src="/write_icon.svg"
                      alt="Edit"
                      width="20"
                      height="20"
                    />
                  </button>
                  <button
                    class="icon-btn delete"
                    onclick={(e) => {
                      e.stopPropagation();
                      handleDeleteRule(rule.id);
                    }}
                    title="Delete rule"
                  >
                    <img
                      src="/bin_icon.svg"
                      alt="Delete"
                      width="20"
                      height="20"
                    />
                  </button>
                  <button
                    class="icon-btn more"
                    onclick={(e) => e.stopPropagation()}
                    title="More options"
                  >
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                    >
                      <circle cx="10" cy="4" r="1.5" />
                      <circle cx="10" cy="10" r="1.5" />
                      <circle cx="10" cy="16" r="1.5" />
                    </svg>
                  </button>
                </div>
              </div>
            {:else}
              <div class="empty-rules">
                <p>No rules yet. Create your first rule!</p>
              </div>
            {/each}
          </div>
        </div>
      {:else}
        <!-- Create/Edit Rule Form -->
        <div class="rules-section">
          <div class="rules-header">
            <div class="rules-title">
              <button class="back-btn" onclick={handleBackToRulesList}>
                <img src="/arrowback_icon.svg" alt="" width="20" height="20" />
              </button>
              <div>
                <h2>Name and describe your rule</h2>
                <p class="sub">
                  Rules set the expectations for members and redditors visiting
                  your community
                </p>
              </div>
            </div>
          </div>

          <form class="rule-form">
            <div class="form-row">
              <label>Rule name<span class="required">*</span></label>
              <div class="pill-input">
                <input
                  type="text"
                  placeholder="Rule name"
                  maxlength="100"
                  bind:value={ruleName}
                />
                <span class="char-count">{100 - ruleName.length}</span>
              </div>
            </div>

            <div class="form-row">
              <label>Description<span class="required">*</span></label>
              <div class="pill-input textarea">
                <textarea
                  placeholder="Description"
                  maxlength="500"
                  bind:value={ruleDescription}
                ></textarea>
                <span class="char-count">{500 - ruleDescription.length}</span>
              </div>
            </div>

            <h3 class="section-heading">Reporting</h3>
            <p class="small">
              Users or mods can select a report reason when reporting content
            </p>

            <div class="form-row">
              <label>Report reason</label>
              <div class="pill-input">
                <input
                  type="text"
                  placeholder="Report reason"
                  maxlength="100"
                  bind:value={reportReason}
                />
                <span class="char-count">{100 - reportReason.length}</span>
              </div>
              <p class="hint">
                By default, this is the same as your rule name. Max characters
                100
              </p>
            </div>

            <div class="form-row save-row">
              <button
                class="save-btn"
                class:enabled={isRuleFormValid}
                disabled={!isRuleFormValid}
                onclick={handleSaveRule}
                type="button"
              >
                Save
              </button>
            </div>
          </form>
        </div>
      {/if}
    {/if}
  </main>
</div>

<!-- Ban User Modal -->
{#if showBanModal}
  <div class="modal-overlay" onclick={handleCloseBanModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>Ban user</h2>

      <div class="form-group">
        <div class="search-input-wrapper">
          <img
            src="/searchuser_icon.svg"
            alt="Search"
            class="search-icon"
            width="20"
            height="20"
          />
          <input
            type="text"
            placeholder="Search for user"
            bind:value={banUsername}
            class="search-input"
          />
        </div>
        <p class="hint">Enter to find user*</p>
      </div>

      <div class="form-group">
        <select bind:value={banRule} class="modal-select">
          <option value="">Rules</option>
          {#each mockCommunityRules as rule}
            <option value={rule.name}>{rule.name}</option>
          {/each}
        </select>
      </div>

      <div class="form-group">
        <select bind:value={banDuration} class="modal-select">
          <option value="">Duration</option>
          <option value="1h">1 hour</option>
          <option value="24h">24 hours</option>
          <option value="7d">7 days</option>
          <option value="30d">30 days</option>
          <option value="permanent">Permanent</option>
        </select>
      </div>

      <div class="form-group">
        <textarea
          placeholder="Reason"
          bind:value={banReason}
          class="modal-textarea"
        ></textarea>
      </div>

      <div class="form-group">
        <textarea placeholder="Note" bind:value={banNote} class="modal-textarea"
        ></textarea>
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseBanModal}>
          Cancel
        </button>
        <button class="btn-danger" onclick={handleBanUser}>Banned</button>
      </div>
    </div>
  </div>
{/if}

<!-- Mute User Modal -->
{#if showMuteModal}
  <div class="modal-overlay" onclick={handleCloseMuteModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>Mute user</h2>

      <div class="form-group">
        <div class="search-input-wrapper">
          <img
            src="/searchuser_icon.svg"
            alt="Search"
            class="search-icon"
            width="20"
            height="20"
          />
          <input
            type="text"
            placeholder="Search for user"
            bind:value={banUsername}
            class="search-input"
          />
        </div>
        <p class="hint">Enter to find user*</p>
      </div>

      <div class="form-group">
        <select bind:value={banRule} class="modal-select">
          <option value="">Rules</option>
          {#each mockCommunityRules as rule}
            <option value={rule.name}>{rule.name}</option>
          {/each}
        </select>
      </div>

      <div class="form-group">
        <select bind:value={banDuration} class="modal-select">
          <option value="">Duration</option>
          <option value="1h">1 hour</option>
          <option value="24h">24 hours</option>
          <option value="7d">7 days</option>
          <option value="30d">30 days</option>
          <option value="permanent">Permanent</option>
        </select>
      </div>

      <div class="form-group">
        <textarea
          placeholder="Reason"
          bind:value={banReason}
          class="modal-textarea"
        ></textarea>
      </div>

      <div class="form-group">
        <textarea placeholder="Note" bind:value={banNote} class="modal-textarea"
        ></textarea>
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseMuteModal}>
          Cancel
        </button>
        <button class="btn-danger" onclick={handleMuteUser}>Muted</button>
      </div>
    </div>
  </div>
{/if}

<!-- Invite/Add User Modal -->
{#if showInviteModal}
  <div class="modal-overlay" onclick={handleCloseInviteModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>
        {inviteType === "mod" ? "Invite Moderator" : "Add Approved User"}
      </h2>

      <div class="form-group">
        <div class="search-input-wrapper">
          <img
            src="/searchuser_icon.svg"
            alt="Search"
            class="search-icon"
            width="20"
            height="20"
          />
          <input
            type="text"
            placeholder="Search for user"
            bind:value={inviteUsername}
            class="search-input"
          />
        </div>
        <p class="hint">Enter username to search</p>
      </div>

      {#if inviteType === "mod"}
        <div class="form-group">
          <label>Permissions</label>
          <select bind:value={invitePermission} class="permission-select">
            <option value="Everything">Everything</option>
            <option value="Manage Posts & Comments"
              >Manage Posts & Comments</option
            >
            <option value="Manage Users">Manage Users</option>
            <option value="Manage Settings">Manage Settings</option>
          </select>
        </div>

        <div class="form-group">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={inviteCanEdit} />
            <span>You can edit this moderator</span>
          </label>
        </div>
      {/if}

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseInviteModal}>
          Cancel
        </button>
        <button class="action-btn-primary" onclick={handleInviteUser}>
          {inviteType === "mod" ? "Invite" : "Add"}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .mod-tools-page {
    display: flex;
    min-height: 100vh;
    background: #f6f7f8;
  }

  /* Sidebar */
  .mod-sidebar {
    width: 240px;
    background: white;
    border-right: 1px solid #edeff1;
    padding: 20px 0;
    position: sticky;
    top: 0;
    height: 100vh;
    overflow-y: auto;
  }

  .exit-mod-tools {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    width: 100%;
    background: transparent;
    border: none;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    transition: background 0.2s;
  }

  .exit-mod-tools:hover {
    background: #f6f7f8;
  }

  .mod-nav {
    margin-top: 20px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 20px;
    width: 100%;
    background: transparent;
    border: none;
    border-left: 3px solid transparent;
    font-size: 14px;
    font-weight: 500;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
  }

  .nav-item:hover {
    background: #f6f7f8;
  }

  .nav-item.active {
    background: #f6f7f8;
    border-left-color: var(--blue--);
    color: var(--blue--);
    font-weight: 600;
  }

  .nav-item img {
    flex-shrink: 0;
  }

  /* Main Content */
  .mod-content {
    flex: 1;
    padding: 24px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
  }

  .queue-section {
    background: white;
    border-radius: 8px;
    overflow: hidden;
  }

  .queue-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 24px;
    border-bottom: 1px solid #edeff1;
  }

  .queue-header h1 {
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0;
  }

  /* Sort Options - Same style as Profile page */
  .sort-options {
    padding: 0 1rem;
    position: relative;
  }

  .sort-options::after {
    content: "";
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    width: 20px;
    height: 20px;
    background-image: url("/Sort.jpg");
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    pointer-events: none;
    opacity: 0.8;
  }

  .sort-options select {
    padding: 0.5rem 2rem 0.5rem 2.75rem;
    border: none;
    border-radius: 4px;
    background-color: #f8f9fa;
    background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='1' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 1em;
    color: #1a1a1b;
    font-size: 0.9rem;
    cursor: pointer;
    font-family: "Roboto", sans-serif;
    font-weight: 400;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
    transition: all 0.2s ease;
  }

  .sort-options select:not(:focus) {
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .sort-options select:hover {
    background-color: #f0f1f2;
  }

  .sort-options select:focus {
    outline: none;
    background-color: #fff;
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.06);
  }

  .sort-options select option {
    padding: 0.75rem 1rem;
    background-color: white;
    color: #1a1a1b;
    font-size: 0.9rem;
    font-weight: 400;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .sort-options select option:hover {
    background-color: #f8f9fa;
    color: #00008b;
  }

  .sort-options select option:checked {
    background-color: #f0f1f2;
    font-weight: 500;
  }

  .sort-options select:focus {
    border-radius: 4px;
  }

  @media screen and (-webkit-min-device-pixel-ratio: 0) {
    .sort-options select {
      border-radius: 4px !important;
    }

    .sort-options select:focus {
      border: none;
    }

    .sort-options select option:checked {
      background: #f0f1f2 linear-gradient(0deg, #f0f1f2 0%, #f0f1f2 100%);
      font-weight: 500;
    }

    .sort-options select option:hover {
      background: #e8f0fe linear-gradient(0deg, #e8f0fe 0%, #e8f0fe 100%);
      color: #00008b;
    }
  }

  .sort-select {
    padding: 8px 12px;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 14px;
    background: white;
    cursor: pointer;
  }

  /* Tabs */
  .queue-tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid #edeff1;
    padding: 0 24px;
    background: white;
  }

  .tab-btn {
    padding: 12px 16px;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 14px;
    font-weight: 600;
    color: #7c7c7c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .tab-btn.active {
    color: #1c1c1c;
    border-bottom-color: var(--blue--);
  }

  .tab-btn:hover {
    color: #1c1c1c;
  }

  /* Posts List */
  .posts-list {
    padding: 24px;
  }

  .post-card {
    background: #f6f7f8;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
  }

  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .post-author {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .author-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }

  .author-avatar-placeholder {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--blue--);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
  }

  .author-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .author-name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .post-time {
    font-size: 12px;
    color: var(--grayfont);
  }

  .report-badge {
    padding: 4px 8px;
    background: #ff4444;
    color: white;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
  }

  .post-title {
    font-size: 16px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0 0 8px 0;
  }

  .post-content {
    font-size: 14px;
    color: #1c1c1c;
    margin: 0 0 12px 0;
    line-height: 1.5;
  }

  .report-info,
  .removed-info {
    padding: 12px;
    background: #fff3cd;
    border-left: 3px solid #ffc107;
    border-radius: 4px;
    font-size: 13px;
    margin-bottom: 12px;
  }

  .removed-info {
    background: #f8d7da;
    border-left-color: #dc3545;
  }

  .post-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .action-btn {
    padding: 8px 16px;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 4px;
    font-size: 13px;
    font-weight: 600;
    color: #1c1c1c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn:hover {
    background: #f6f7f8;
  }

  .action-btn.approve {
    background: var(--blue--);
    color: white;
    border: 1px solid var(--blue--);
    border-radius: 16px;
  }

  .action-btn.approve:hover {
    background: #0000cd;
    border-color: #0000cd;
  }

  .action-btn.remove {
    background: var(--button-secondary-bg);
    color: #1c1c1c;
    border: 1px solid transparent;
    border-radius: 16px;
  }

  .action-btn.remove:hover {
    background: rgba(214, 216, 222, 0.6);
  }

  .empty-state {
    text-align: center;
    padding: 48px 24px;
    color: var(--grayfont);
  }

  .placeholder-section {
    background: white;
    border-radius: 8px;
    padding: 48px 24px;
    text-align: center;
  }

  /* Rules List View */
  .rules-list-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .rules-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .rules-list-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .create-rule-btn {
    padding: 10px 20px;
    background: var(--blue--);
    color: white;
    border: none;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .create-rule-btn:hover {
    background: #0000cd;
  }

  .rules-table {
    border: 1px solid #edeff1;
    border-radius: 8px;
    overflow: hidden;
  }

  .rules-table-header {
    display: grid;
    grid-template-columns: 1fr 150px;
    background: #f6f7f8;
    padding: 12px 16px;
    font-size: 12px;
    font-weight: 700;
    color: var(--grayfont);
    text-transform: uppercase;
  }

  .rule-row {
    display: grid;
    grid-template-columns: 1fr 150px;
    padding: 16px;
    border-top: 1px solid #edeff1;
    align-items: center;
    cursor: pointer;
    transition: background 0.2s;
  }

  .rule-row:first-child {
    border-top: none;
  }

  .rule-row:hover {
    background: #f9f9f9;
  }

  .rule-info {
    display: flex;
    gap: 16px;
    align-items: flex-start;
  }

  .rule-number {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
    min-width: 20px;
  }

  .rule-details {
    flex: 1;
  }

  .rule-name {
    margin: 0 0 4px 0;
    font-size: 15px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .rule-desc {
    margin: 0;
    font-size: 13px;
    color: var(--grayfont);
  }

  .rule-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }

  .icon-btn {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.2s;
    color: #1c1c1c;
  }

  .icon-btn:hover {
    background: #f6f7f8;
  }

  .icon-btn.delete img {
    filter: brightness(0) saturate(100%) invert(27%) sepia(93%) saturate(4373%)
      hue-rotate(348deg) brightness(88%) contrast(88%);
  }

  .icon-btn.delete:hover {
    background: #fff0f0;
  }

  .empty-rules {
    padding: 48px 24px;
    text-align: center;
    color: var(--grayfont);
  }

  /* Rules section styles */
  .rules-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .rules-header {
    margin-bottom: 24px;
  }

  .rules-title {
    display: flex;
    align-items: flex-start;
    gap: 16px;
  }

  .back-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    transition: background 0.2s;
  }

  .back-btn:hover {
    background: #f6f7f8;
  }

  .rules-title h2 {
    margin: 0;
    font-size: 22px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .rules-title .sub {
    margin: 6px 0 0 0;
    color: var(--grayfont);
    font-size: 14px;
  }

  .rule-form {
    margin-top: 12px;
  }

  .form-row {
    margin-bottom: 18px;
  }

  .form-row label {
    display: block;
    font-size: 14px;
    color: #1c1c1c;
    margin-bottom: 8px;
    font-weight: 600;
  }

  .required {
    color: #d9534f;
    margin-left: 6px;
  }

  .pill-input {
    display: flex;
    align-items: center;
    background: #eef2f4;
    padding: 12px 16px;
    border-radius: 16px;
    position: relative;
  }

  .pill-input input {
    border: none;
    background: transparent;
    width: 100%;
    font-size: 14px;
    color: #1c1c1c;
    outline: none;
  }

  .pill-input.textarea {
    align-items: flex-start;
    padding-bottom: 36px;
  }

  .pill-input.textarea textarea {
    width: 100%;
    border: none;
    background: transparent;
    min-height: 120px;
    resize: vertical;
    font-size: 14px;
    color: #1c1c1c;
    outline: none;
    font-family: inherit;
  }

  .char-count {
    position: absolute;
    right: 16px;
    bottom: 12px;
    font-size: 12px;
    color: var(--grayfont);
  }

  .section-heading {
    margin-top: 16px;
    margin-bottom: 8px;
    font-size: 16px;
    font-weight: 700;
  }

  .small {
    color: var(--grayfont);
    font-size: 13px;
  }

  .hint {
    font-size: 12px;
    color: var(--grayfont);
    margin-top: 6px;
  }

  .save-row {
    display: flex;
    justify-content: flex-end;
    margin-top: 24px;
  }

  .save-btn {
    padding: 10px 24px;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    border: 1px solid transparent;
  }

  .save-btn:disabled {
    background: var(--button-secondary-bg);
    color: #1c1c1c;
    cursor: not-allowed;
  }

  .save-btn.enabled {
    background: var(--blue--);
    color: white;
    border-color: var(--blue--);
    cursor: pointer;
  }

  .save-btn.enabled:hover {
    background: #0000cd;
    border-color: #0000cd;
  }

  .placeholder-section h1 {
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
    margin: 0 0 16px 0;
  }

  .placeholder-section p {
    font-size: 16px;
    color: var(--grayfont);
    margin: 0;
  }

  /* Restricted Users Section */
  .restricted-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .restricted-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .restricted-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .action-btn-primary {
    padding: 10px 20px;
    background: var(--blue--);
    color: white;
    border: none;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn-primary:hover {
    background: #00008b;
    filter: brightness(0.85);
  }

  .restricted-tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid #edeff1;
    margin-bottom: 24px;
  }

  .restricted-tabs .tab-btn {
    padding: 12px 24px;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 14px;
    font-weight: 600;
    color: #7c7c7c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .restricted-tabs .tab-btn.active {
    color: #1c1c1c;
    border-bottom-color: var(--blue--);
  }

  .restricted-tabs .tab-btn:hover {
    color: #1c1c1c;
  }

  .restricted-table {
    border: 1px solid #edeff1;
    border-radius: 8px;
    overflow: hidden;
  }

  .table-header {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
    background: #f6f7f8;
    padding: 12px 16px;
    font-size: 12px;
    font-weight: 700;
    color: var(--grayfont);
    text-transform: uppercase;
  }

  .table-header.muted {
    grid-template-columns: 1fr 1fr 1fr;
  }

  .table-row {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr 1fr;
    padding: 16px;
    border-top: 1px solid #edeff1;
    align-items: center;
    font-size: 14px;
    color: #1c1c1c;
  }

  .table-row .col {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .restricted-table .table-row:has(+ .table-row.muted),
  .table-row.muted {
    grid-template-columns: 1fr 1fr 1fr;
  }

  /* Modal Styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    background: white;
    border-radius: 12px;
    padding: 24px;
    width: 90%;
    max-width: 440px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  .modal-content h2 {
    margin: 0 0 20px 0;
    font-size: 20px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .form-group {
    margin-bottom: 16px;
  }

  .search-input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-input-wrapper .search-icon {
    position: absolute;
    left: 16px;
    pointer-events: none;
    opacity: 0.6;
  }

  .search-input,
  .modal-select {
    width: 100%;
    padding: 12px 16px;
  }

  .search-input-wrapper .search-input {
    padding-left: 48px;
  }

  .search-input,
  .modal-select {
    border: 1px solid #edeff1;
    border-radius: 8px;
    font-size: 14px;
    background: #f6f7f8;
    color: #1c1c1c;
  }

  .search-input:focus,
  .modal-select:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .modal-textarea {
    width: 100%;
    min-height: 80px;
    padding: 12px 16px;
    border: 1px solid #edeff1;
    border-radius: 8px;
    font-size: 14px;
    background: #f6f7f8;
    color: #1c1c1c;
    resize: vertical;
    font-family: inherit;
  }

  .modal-textarea:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .hint {
    font-size: 12px;
    color: var(--grayfont);
    margin: 6px 0 0 0;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
  }

  .btn-cancel {
    padding: 10px 20px;
    background: var(--button-secondary-bg);
    color: var(--blue--);
    border: none;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-cancel:hover {
    background: rgba(214, 216, 222, 0.6);
  }

  .btn-danger {
    padding: 10px 20px;
    background: #ff4444;
    color: white;
    border: none;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-danger:hover {
    background: #ff0000;
  }

  .permission-select {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #edeff1;
    border-radius: 8px;
    font-size: 14px;
    background: #f6f7f8;
    color: #1c1c1c;
    cursor: pointer;
  }

  .permission-select:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    user-select: none;
  }

  .checkbox-label input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
    accent-color: var(--blue--);
  }

  .checkbox-label span {
    font-size: 14px;
    color: #1c1c1c;
  }

  /* Mod & Members Section */
  .members-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .members-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .members-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .members-tabs {
    display: flex;
    gap: 0;
    border-bottom: 1px solid #edeff1;
    margin-bottom: 24px;
  }

  .members-tabs .tab-btn {
    padding: 12px 24px;
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 14px;
    font-weight: 600;
    color: #7c7c7c;
    cursor: pointer;
    transition: all 0.2s;
  }

  .members-tabs .tab-btn.active {
    color: #1c1c1c;
    border-bottom-color: var(--blue--);
  }

  .members-tabs .tab-btn:hover {
    color: #1c1c1c;
  }

  .members-table {
    border: 1px solid #edeff1;
    border-radius: 8px;
    overflow: hidden;
  }

  .members-table .table-header {
    display: grid;
    background: #f6f7f8;
    padding: 12px 16px;
    font-size: 12px;
    font-weight: 700;
    color: var(--grayfont);
    text-transform: uppercase;
  }

  .members-table .table-header.moderators {
    grid-template-columns: 2fr 1fr 1fr 1fr 120px;
  }

  .members-table .table-header.approved {
    grid-template-columns: 2fr 1fr 120px;
  }

  .members-table .table-row {
    display: grid;
    padding: 16px;
    border-top: 1px solid #edeff1;
    align-items: center;
    font-size: 14px;
    color: #1c1c1c;
  }

  .members-table .table-row.moderators {
    grid-template-columns: 2fr 1fr 1fr 1fr 120px;
  }

  .members-table .table-row.approved {
    grid-template-columns: 2fr 1fr 120px;
  }

  .user-col {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }

  .joined-col {
    white-space: pre-line;
    font-size: 13px;
    line-height: 1.4;
  }

  .actions-col {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
  }
</style>
