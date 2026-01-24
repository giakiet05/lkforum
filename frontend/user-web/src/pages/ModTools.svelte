<script lang="ts">
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import { mockQueuePosts } from "../mocks/mod-queue.mock";
  import type { QueuePost } from "../mocks/mod-queue.mock";
  import { mockRestrictedUsers } from "../mocks/restricted-users.mock";
  import type { RestrictedUser } from "../mocks/restricted-users.mock";
  import { mockModerators } from "../mocks/moderators.mock";
  import type { Moderator } from "../mocks/moderators.mock";
  import ConfirmModal from "../components/ConfirmModal.svelte";
  import {
    getCommunities,
    updateCommunity,
    banUser,
    unbanUser,
    unmuteUser,
    getBannedUsers,
    addModerators,
    removeModerators,
    getPendingPosts,
    getEditedPosts,
    moderatePost,
  } from "../services/community-service";
  import {
    getMembershipsByCommunityId,
    kickMember,
    getPendingMembers,
    approveMember,
    rejectMember,
    getApprovedMembers,
  } from "../services/membership-service";
  import { getUserByUsername } from "../services/user-service";
  import { toastStore } from "../stores/toast-store";
  import { authStore } from "../stores/auth-store";
  import type {
    CommunityResponse,
    CommunityRule,
    UserResponse,
    Moderator as ModeratorType,
  } from "../dtos/community-dto";

  export interface Props {
    params?: { name?: string };
  }

  let { params }: Props = $props();

  const communityName = params?.name || "";

  const currentUser = $derived($authStore.user);
  const isCreator = $derived(() => {
    if (!community || !currentUser) return false;
    return community.create_by_id === currentUser.id;
  });

  // Community data
  let community = $state<CommunityResponse | null>(null);
  let isLoadingCommunity = $state(true);
  let communityRules = $state<CommunityRule[]>([]);
  let bannedUsers = $state<UserResponse[]>([]);
  let mutedUsers = $state<UserResponse[]>([]);
  let moderators = $state<ModeratorType[]>([]);
  let isLoadingRestricted = $state(false);
  let isLoadingModerators = $state(false);

  type SidebarItem =
    | "queue"
    | "restricted"
    | "members"
    | "all-members"
    | "pending-members"
    | "rules";
  type QueueTab = "unmoderated" | "edited" | "removed" | "reported";
  type SortOption = "newest" | "oldest" | "most-reported";
  type RestrictedTab = "banned" | "muted";
  type MembersTab = "moderators" | "approved";

  let activeSidebarItem = $state<SidebarItem>("queue");
  let activeQueueTab = $state<QueueTab>("unmoderated");
  let sortBy = $state<SortOption>("newest");

  // All members state
  let allMembers = $state<any[]>([]);
  let isLoadingAllMembers = $state(false);
  let membersPage = $state(1);
  let membersPageSize = $state(20);
  let isKickingMember = $state<string | null>(null);

  // Pending members state (for communities with join_require_approval)
  let pendingMembers = $state<any[]>([]);
  let isLoadingPendingMembers = $state(false);
  let pendingMembersPage = $state(1);
  let pendingMembersPageSize = $state(20);
  let pendingMembersCount = $state(0);
  let isProcessingMember = $state<string | null>(null);

  // Approved members state
  let approvedMembers = $state<any[]>([]);
  let isLoadingApprovedMembers = $state(false);
  let approvedMembersPage = $state(1);
  let approvedMembersPageSize = $state(20);
  let approvedMembersCount = $state(0);

  // Queue/Pending Posts state
  let pendingPosts = $state<any[]>([]);
  let editedPosts = $state<any[]>([]);
  let isLoadingPosts = $state(false);
  let postsPage = $state(1);
  let postsPageSize = $state(10);
  let totalPendingPosts = $state(0);
  let totalEditedPosts = $state(0);

  // Restricted Users state
  let activeRestrictedTab = $state<RestrictedTab>("banned");
  let showBanModal = $state(false);
  let showMuteModal = $state(false);
  let banUsername = $state("");
  let banRule = $state("");
  let banDuration = $state("");
  let banReason = $state("");
  let banNote = $state("");
  let showUserSuggestions = $state(false);
  let selectedSuggestionIndex = $state(-1);

  // Filter members for suggestions based on input
  const filteredMemberSuggestions = $derived(() => {
    if (!banUsername.trim() || !allMembers.length) return [];
    const search = banUsername.toLowerCase().trim();
    return allMembers
      .filter((m) => m.user?.username?.toLowerCase().includes(search))
      .slice(0, 5); // Limit to 5 suggestions
  });

  // Mod & Members state
  let activeMembersTab = $state<MembersTab>("moderators");
  let showInviteModal = $state(false);
  let inviteUsername = $state("");
  let inviteType = $state<"mod" | "approved">("mod");
  let invitePermission = $state("Everything");
  let inviteCanEdit = $state(true);

  // Rules state
  let showRuleForm = $state(false);
  let editingRuleIndex = $state<number | null>(null);
  let ruleName = $state("");
  let ruleDescription = $state("");

  // Confirm modal states
  let showKickConfirm = $state(false);
  let showRejectConfirm = $state(false);
  let showRemoveApprovedConfirm = $state(false);
  let showDeleteRuleConfirm = $state(false);
  let showUnbanConfirm = $state(false);
  let showUnmuteConfirm = $state(false);
  let showDeleteMemberConfirm = $state(false);
  let showRemovePostModal = $state(false);
  let removePostId = $state<string | null>(null);
  let removePostReason = $state("");
  let confirmTargetUser = $state<{ id: string; username: string } | null>(null);
  let confirmTargetIndex = $state<number | null>(null);
  let confirmMemberType = $state<"mod" | "approved">("mod");

  const isRuleFormValid = $derived(
    ruleName.trim().length > 0 && ruleDescription.trim().length > 0,
  );

  // Load community data
  onMount(async () => {
    await loadCommunity();
    await loadRestrictedUsers();
    await loadModerators();
    await loadPendingPosts();
    await loadAllMembers();
    await loadApprovedMembersList();
  });

  // Watch for tab changes and load appropriate posts
  $effect(() => {
    if (community) {
      if (activeQueueTab === "unmoderated") {
        loadPendingPosts();
      } else if (activeQueueTab === "edited") {
        loadEditedPosts();
      }
    }
  });

  // Watch for sidebar changes and load all members
  $effect(() => {
    if (community && activeSidebarItem === "all-members") {
      loadAllMembers();
    }
  });

  // Watch for sidebar changes and load pending members
  $effect(() => {
    if (community && activeSidebarItem === "pending-members") {
      loadPendingMembersList();
    }
  });

  // Load pending members count on community load (for badge)
  $effect(() => {
    if (community?.setting?.join_require_approval) {
      loadPendingMembersList();
    }
  });

  async function loadCommunity() {
    try {
      isLoadingCommunity = true;
      const response = await getCommunities({ limit: 100 });
      const found = response.communities.find((c) => c.name === communityName);

      if (found) {
        community = found;
        communityRules = found.rules || [];
        moderators = found.moderators || [];
      } else {
        console.error("Community not found");
      }
    } catch (error) {
      console.error("Failed to load community:", error);
    } finally {
      isLoadingCommunity = false;
    }
  }

  async function loadRestrictedUsers() {
    if (!community) return;

    try {
      isLoadingRestricted = true;
      const [banned, muted] = await Promise.all([
        getBannedUsers(community.id, "banned"),
        getBannedUsers(community.id, "muted"),
      ]);
      bannedUsers = banned;
      mutedUsers = muted;
    } catch (error) {
      console.error("Failed to load restricted users:", error);
    } finally {
      isLoadingRestricted = false;
    }
  }

  async function loadModerators() {
    if (!community) return;

    try {
      isLoadingModerators = true;
      // Moderators are already loaded in community response
      moderators = community.moderators || [];
    } catch (error) {
      console.error("Failed to load moderators:", error);
    } finally {
      isLoadingModerators = false;
    }
  }

  async function loadAllMembers() {
    if (!community) return;

    try {
      isLoadingAllMembers = true;
      const members = await getMembershipsByCommunityId(
        community.id,
        membersPage,
        membersPageSize,
      );
      allMembers = members || [];
    } catch (error) {
      console.error("Failed to load all members:", error);
      allMembers = [];
    } finally {
      isLoadingAllMembers = false;
    }
  }

  async function handleKickMember(userId: string, username: string) {
    if (!community) return;

    // Prevent kicking creator
    if (userId === community.create_by_id) {
      toastStore.error("Không thể xóa người tạo cộng đồng!");
      return;
    }

    confirmTargetUser = { id: userId, username };
    showKickConfirm = true;
  }

  async function confirmKickMember() {
    if (!community || !confirmTargetUser) return;
    showKickConfirm = false;

    try {
      isKickingMember = confirmTargetUser.id;
      await kickMember(community.id, confirmTargetUser.id);
      toastStore.success(
        `Đã xóa ${confirmTargetUser.username} khỏi cộng đồng!`,
      );
      await loadAllMembers(); // Reload the list
    } catch (error: any) {
      console.error("Failed to kick member:", error);
      toastStore.error(
        error.message || "Không thể xóa thành viên. Vui lòng thử lại.",
      );
    } finally {
      isKickingMember = null;
      confirmTargetUser = null;
    }
  }

  // Pending members functions
  async function loadPendingMembersList() {
    if (!community) return;

    try {
      isLoadingPendingMembers = true;
      const response = await getPendingMembers(
        community.id,
        pendingMembersPage,
        pendingMembersPageSize,
      );
      console.log("🔍 Pending members response:", response);
      console.log("🔍 First pending member:", response.memberships?.[0]);
      pendingMembers = response.memberships || [];
      pendingMembersCount = response.pagination?.total || 0;
    } catch (error) {
      console.error("Failed to load pending members:", error);
      pendingMembers = [];
      pendingMembersCount = 0;
    } finally {
      isLoadingPendingMembers = false;
    }
  }

  // Approved members functions
  async function loadApprovedMembersList() {
    if (!community) return;

    try {
      isLoadingApprovedMembers = true;
      const response = await getApprovedMembers(
        community.id,
        approvedMembersPage,
        approvedMembersPageSize,
      );
      console.log("🔍 Approved members response:", response);
      console.log("🔍 First member:", response.memberships?.[0]);
      approvedMembers = response.memberships || [];
      approvedMembersCount = response.pagination?.total || 0;
    } catch (error) {
      console.error("Failed to load approved members:", error);
      approvedMembers = [];
      approvedMembersCount = 0;
    } finally {
      isLoadingApprovedMembers = false;
    }
  }

  async function handleApproveMember(membershipId: string, username: string) {
    if (!community) return;

    try {
      isProcessingMember = membershipId;
      await approveMember(community.id, membershipId);
      toastStore.success(`Đã duyệt ${username} vào cộng đồng!`);
      await loadPendingMembersList();
      await loadApprovedMembersList();
      await loadAllMembers();
    } catch (error: any) {
      console.error("Failed to approve member:", error);
      toastStore.error(
        error.message || "Không thể duyệt thành viên. Vui lòng thử lại.",
      );
    } finally {
      isProcessingMember = null;
    }
  }

  async function handleRejectMember(membershipId: string, username: string) {
    if (!community) return;

    confirmTargetUser = { id: membershipId, username };
    showRejectConfirm = true;
  }

  async function confirmRejectMember() {
    if (!community || !confirmTargetUser) return;
    showRejectConfirm = false;

    try {
      isProcessingMember = confirmTargetUser.id;
      await rejectMember(community.id, confirmTargetUser.id);
      toastStore.success(
        `Đã từ chối yêu cầu của ${confirmTargetUser.username}!`,
      );
      await loadPendingMembersList();
    } catch (error: any) {
      console.error("Failed to reject member:", error);
      toastStore.error(
        error.message || "Không thể từ chối yêu cầu. Vui lòng thử lại.",
      );
    } finally {
      isProcessingMember = null;
      confirmTargetUser = null;
    }
  }

  async function handleRemoveApprovedUser(userId: string, username: string) {
    if (!community) return;

    confirmTargetUser = { id: userId, username };
    showRemoveApprovedConfirm = true;
  }

  async function confirmRemoveApprovedUser() {
    if (!community || !confirmTargetUser) return;
    showRemoveApprovedConfirm = false;

    try {
      isProcessingMember = confirmTargetUser.id;
      await kickMember(community.id, confirmTargetUser.id);
      toastStore.success(
        `Đã xóa ${confirmTargetUser.username} khỏi danh sách!`,
      );
      await loadApprovedMembersList();
      await loadAllMembers();
    } catch (error: any) {
      console.error("Failed to remove approved user:", error);
      toastStore.error(
        error.message || "Không thể xóa người dùng. Vui lòng thử lại.",
      );
    } finally {
      isProcessingMember = null;
      confirmTargetUser = null;
    }
  }

  async function loadPendingPosts() {
    if (!community) return;

    try {
      isLoadingPosts = true;
      const response = await getPendingPosts(
        community.id,
        postsPage,
        postsPageSize,
      );
      pendingPosts = response.posts || [];
      totalPendingPosts = response.pagination?.total_items || 0;
    } catch (error) {
      console.error("Failed to load pending posts:", error);
      pendingPosts = [];
    } finally {
      isLoadingPosts = false;
    }
  }

  async function loadEditedPosts() {
    if (!community) return;

    try {
      isLoadingPosts = true;
      const response = await getEditedPosts(
        community.id,
        postsPage,
        postsPageSize,
      );
      editedPosts = response.posts || [];
      totalEditedPosts = response.pagination?.total_items || 0;
    } catch (error) {
      console.error("Failed to load edited posts:", error);
      editedPosts = [];
    } finally {
      isLoadingPosts = false;
    }
  }

  const filteredPosts = $derived(() => {
    let posts: any[] = [];
    if (activeQueueTab === "unmoderated") {
      posts = [...pendingPosts];
    } else if (activeQueueTab === "edited") {
      posts = [...editedPosts];
    }

    // Sort posts
    if (sortBy === "newest") {
      posts.sort(
        (a, b) =>
          new Date(b.created_at || b.createdAt).getTime() -
          new Date(a.created_at || a.createdAt).getTime(),
      );
    } else if (sortBy === "oldest") {
      posts.sort(
        (a, b) =>
          new Date(a.created_at || a.createdAt).getTime() -
          new Date(b.created_at || b.createdAt).getTime(),
      );
    }
    // most-reported sorting not available yet from backend

    return posts;
  });

  function handleExitModTools() {
    push(`/lk/${communityName}`);
  }

  function handleSidebarClick(item: SidebarItem) {
    activeSidebarItem = item;
  }

  function formatTime(dateString: string): string {
    if (!dateString) return "Không xác định";
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return "Không xác định";

    // Check for zero value time (0001-01-01 in Go)
    if (date.getFullYear() <= 1) return "Không xác định";

    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffHours / 24);

    if (diffHours < 1) return "Vừa xong";
    if (diffHours < 24) return `${diffHours} giờ trước`;
    if (diffDays < 7) return `${diffDays} ngày trước`;
    return date.toLocaleDateString("vi-VN");
  }

  async function handleApprove(postId: string) {
    if (!community) return;

    try {
      await moderatePost(community.id, postId, true);
      toastStore.success("Post approved successfully!");
      await loadPendingPosts(); // Reload the list
    } catch (error) {
      console.error("Failed to approve post:", error);
      toastStore.error("Failed to approve post. Please try again.");
    }
  }

  async function handleRemove(postId: string, reason?: string) {
    if (!community) return;

    if (reason) {
      // If reason is already provided, proceed directly
      try {
        await moderatePost(community.id, postId, false, reason);
        toastStore.success("Đã xóa bài viết!");
        await loadPendingPosts();
      } catch (error) {
        console.error("Failed to remove post:", error);
        toastStore.error("Không thể xóa bài viết. Vui lòng thử lại.");
      }
    } else {
      // Show modal to get reason
      removePostId = postId;
      removePostReason = "";
      showRemovePostModal = true;
    }
  }

  async function confirmRemovePost() {
    if (!community || !removePostId) return;
    showRemovePostModal = false;

    try {
      await moderatePost(
        community.id,
        removePostId,
        false,
        removePostReason.trim() || undefined,
      );
      toastStore.success("Đã xóa bài viết!");
      await loadPendingPosts();
    } catch (error) {
      console.error("Failed to remove post:", error);
      toastStore.error("Không thể xóa bài viết. Vui lòng thử lại.");
    } finally {
      removePostId = null;
      removePostReason = "";
    }
  }

  function handleSaveRule() {
    if (!isRuleFormValid || !community) return;

    const newRule: CommunityRule = {
      title: ruleName.trim(),
      description: ruleDescription.trim(),
    };

    let updatedRules: CommunityRule[];
    if (editingRuleIndex !== null) {
      // Edit existing rule
      updatedRules = [...communityRules];
      updatedRules[editingRuleIndex] = newRule;
    } else {
      // Create new rule
      updatedRules = [...communityRules, newRule];
    }

    // Update community via API
    updateCommunity({
      id: community.id,
      rules: updatedRules,
    })
      .then(() => {
        communityRules = updatedRules;
        toastStore.success(
          editingRuleIndex !== null ? "Rule updated!" : "Rule created!",
        );
        handleBackToRulesList();
      })
      .catch((error) => {
        console.error("Failed to save rule:", error);
        toastStore.error("Failed to save rule. Please try again.");
      });
  }

  function handleCreateRule() {
    showRuleForm = true;
    editingRuleIndex = null;
    ruleName = "";
    ruleDescription = "";
  }

  function handleEditRule(index: number) {
    const rule = communityRules[index];
    if (rule) {
      showRuleForm = true;
      editingRuleIndex = index;
      ruleName = rule.title;
      ruleDescription = rule.description;
    }
  }

  function handleDeleteRule(index: number) {
    if (!community) return;

    confirmTargetIndex = index;
    showDeleteRuleConfirm = true;
  }

  async function confirmDeleteRule() {
    if (!community || confirmTargetIndex === null) return;
    showDeleteRuleConfirm = false;

    const updatedRules = communityRules.filter(
      (_, i) => i !== confirmTargetIndex,
    );

    updateCommunity({
      id: community.id,
      rules: updatedRules,
    })
      .then(() => {
        communityRules = updatedRules;
        toastStore.success("Đã xóa quy tắc!");
      })
      .catch((error) => {
        console.error("Failed to delete rule:", error);
        toastStore.error("Không thể xóa quy tắc. Vui lòng thử lại.");
      })
      .finally(() => {
        confirmTargetIndex = null;
      });
  }

  function handleBackToRulesList() {
    showRuleForm = false;
    editingRuleIndex = null;
    ruleName = "";
    ruleDescription = "";
  }

  // Restricted Users functions
  const filteredRestrictedUsers = $derived(
    activeRestrictedTab === "banned" ? bannedUsers : mutedUsers,
  );

  function handleOpenBanModal() {
    showBanModal = true;
    banUsername = "";
    banRule = "";
    banDuration = "";
    banReason = "";
    banNote = "";
    showUserSuggestions = false;
    selectedSuggestionIndex = -1;
  }

  function handleOpenMuteModal() {
    showMuteModal = true;
    banUsername = "";
    banRule = "";
    banDuration = "";
    banReason = "";
    banNote = "";
    showUserSuggestions = false;
    selectedSuggestionIndex = -1;
  }

  function handleCloseBanModal() {
    showBanModal = false;
    showUserSuggestions = false;
  }

  function handleCloseMuteModal() {
    showMuteModal = false;
    showUserSuggestions = false;
  }

  function handleSelectUserSuggestion(username: string) {
    banUsername = username;
    showUserSuggestions = false;
    selectedSuggestionIndex = -1;
  }

  function handleUserInputKeydown(e: KeyboardEvent) {
    const suggestions = filteredMemberSuggestions();
    if (!suggestions.length) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      selectedSuggestionIndex = Math.min(
        selectedSuggestionIndex + 1,
        suggestions.length - 1,
      );
      showUserSuggestions = true;
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      selectedSuggestionIndex = Math.max(selectedSuggestionIndex - 1, 0);
    } else if (e.key === "Enter" && selectedSuggestionIndex >= 0) {
      e.preventDefault();
      const selected = suggestions[selectedSuggestionIndex];
      if (selected) {
        handleSelectUserSuggestion(selected.user?.username || "");
      }
    } else if (e.key === "Escape") {
      showUserSuggestions = false;
      selectedSuggestionIndex = -1;
    }
  }

  function handleUserInputFocus() {
    if (banUsername.trim() && filteredMemberSuggestions().length > 0) {
      showUserSuggestions = true;
    }
  }

  function handleUserInputBlur() {
    // Delay to allow click on suggestion
    setTimeout(() => {
      showUserSuggestions = false;
    }, 200);
  }

  async function handleBanUser() {
    if (!community || !banUsername.trim()) {
      toastStore.warning("Please enter a username!");
      return;
    }

    const lengthDays = parseInt(banDuration) || 30;

    try {
      // Get user ID from username
      const user = await getUserByUsername(banUsername.trim());

      await banUser({
        community_id: community.id,
        user_id: user.id,
        type: "banned",
        reason: banReason || "Violation of community rules",
        length_days: lengthDays,
      });

      toastStore.success("User banned successfully!");
      await loadRestrictedUsers(); // Reload the list
      handleCloseBanModal();
    } catch (error) {
      console.error("Failed to ban user:", error);
      toastStore.error(
        "Failed to ban user. Please check the username and try again.",
      );
    }
  }

  async function handleMuteUser() {
    if (!community || !banUsername.trim()) {
      toastStore.warning("Please enter a username!");
      return;
    }

    const lengthDays = parseInt(banDuration) || 30;

    try {
      // Get user ID from username
      const user = await getUserByUsername(banUsername.trim());

      await banUser({
        community_id: community.id,
        user_id: user.id,
        type: "muted",
        reason: banReason || "Violation of community rules",
        length_days: lengthDays,
      });

      toastStore.success("Đã tắt tiếng người dùng thành công!");
      await loadRestrictedUsers(); // Reload the list
      handleCloseMuteModal();
    } catch (error) {
      console.error("Failed to mute user:", error);
      toastStore.error(
        "Không thể tắt tiếng người dùng. Vui lòng kiểm tra tên người dùng và thử lại.",
      );
    }
  }

  async function handleUnbanUser(userId: string) {
    if (!community) return;

    confirmTargetUser = { id: userId, username: "" };
    showUnbanConfirm = true;
  }

  async function confirmUnbanUser() {
    if (!community || !confirmTargetUser) return;
    showUnbanConfirm = false;

    try {
      await unbanUser({
        community_id: community.id,
        user_id: confirmTargetUser.id,
      });

      toastStore.success("Đã bỏ cấm người dùng thành công!");
      await loadRestrictedUsers();
    } catch (error) {
      console.error("Failed to unban user:", error);
      toastStore.error("Không thể bỏ cấm người dùng. Vui lòng thử lại.");
    } finally {
      confirmTargetUser = null;
    }
  }

  async function handleUnmuteUser(userId: string) {
    if (!community) return;

    confirmTargetUser = { id: userId, username: "" };
    showUnmuteConfirm = true;
  }

  async function confirmUnmuteUser() {
    if (!community || !confirmTargetUser) return;
    showUnmuteConfirm = false;

    try {
      await unmuteUser({
        community_id: community.id,
        user_id: confirmTargetUser.id,
      });

      toastStore.success("Đã bỏ tắt tiếng người dùng thành công!");
      await loadRestrictedUsers();
    } catch (error) {
      console.error("Failed to unmute user:", error);
      toastStore.error("Không thể bỏ tắt tiếng người dùng. Vui lòng thử lại.");
    } finally {
      confirmTargetUser = null;
    }
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

  async function handleInviteUser() {
    if (!community || !inviteUsername.trim()) {
      toastStore.warning("Please enter a username!");
      return;
    }

    console.log("Invite mod:", {
      username: inviteUsername,
      permission: invitePermission,
      canEdit: inviteCanEdit,
    });

    if (inviteType === "mod") {
      try {
        console.log("🔍 Looking up user:", inviteUsername.trim());
        // Get user ID from username
        const user = await getUserByUsername(inviteUsername.trim());
        console.log("✅ User found:", user);

        console.log("📤 Adding moderator to community:", community.id);
        await addModerators({
          id: community.id,
          added_moderator: [user.id], // Just send user ID
        });

        console.log("✅ Moderator added successfully!");
        toastStore.success(`Đã thêm ${inviteUsername} làm quản trị viên!`);
        await loadCommunity(); // Reload to get updated moderators
        handleCloseInviteModal();
      } catch (error) {
        console.error("❌ Failed to add moderator:", error);
        toastStore.error(
          "Không thể thêm quản trị viên. Vui lòng kiểm tra tên người dùng và thử lại.",
        );
      }
    } else {
      // Approved users - not yet implemented in backend
      toastStore.info("Approved users feature coming soon!");
      handleCloseInviteModal();
    }
  }

  function handleEditMember(userId: string, type: "mod" | "approved") {
    console.log("Edit member:", userId, type);
    toastStore.info("Tính năng chỉnh sửa thành viên sẽ có sớm!");
  }

  async function handleDeleteMember(userId: string, type: "mod" | "approved") {
    if (!community) return;

    confirmTargetUser = { id: userId, username: "" };
    confirmMemberType = type;
    showDeleteMemberConfirm = true;
  }

  async function confirmDeleteMember() {
    if (!community || !confirmTargetUser) return;
    showDeleteMemberConfirm = false;

    if (confirmMemberType === "mod") {
      try {
        await removeModerators({
          id: community.id,
          removed_moderator: [confirmTargetUser.id],
        });

        toastStore.success("Đã xóa quản trị viên!");
        await loadCommunity(); // Reload to get updated moderators
      } catch (error) {
        console.error("Failed to remove moderator:", error);
        toastStore.error("Không thể xóa quản trị viên. Vui lòng thử lại.");
      }
    } else {
      // Approved users - not yet implemented
      toastStore.info("Approved users feature coming soon!");
    }
    confirmTargetUser = null;
  }
</script>

<div class="mod-tools-page">
  <!-- Sidebar -->
  <aside class="mod-sidebar">
    <button class="exit-mod-tools" onclick={handleExitModTools}>
      <!-- Use /arrowback_icon.svg (place file into public folder as arrowback_icon.svg) -->
      <img src="/arrowback_icon.svg" alt="" width="20" height="20" />
      Thoát quản trị
    </button>

    <nav class="mod-nav">
      <button
        class="nav-item"
        class:active={activeSidebarItem === "queue"}
        onclick={() => handleSidebarClick("queue")}
      >
        <img src="/queue_icon.svg" alt="" width="20" height="20" />
        <span>Hàng chờ</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "restricted"}
        onclick={() => handleSidebarClick("restricted")}
      >
        <img src="/restricted_icon.svg" alt="" width="20" height="20" />
        <span>Người dùng bị giới hạn</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "members"}
        onclick={() => handleSidebarClick("members")}
      >
        <img src="/member_icon.svg" alt="" width="20" height="20" />
        <span>Quản trị viên</span>
      </button>

      <button
        class="nav-item"
        class:active={activeSidebarItem === "all-members"}
        onclick={() => handleSidebarClick("all-members")}
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
          <circle cx="9" cy="7" r="4" />
          <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
          <path d="M16 3.13a4 4 0 0 1 0 7.75" />
        </svg>
        <span>Tất cả thành viên</span>
      </button>

      {#if community?.setting?.join_require_approval}
        <button
          class="nav-item"
          class:active={activeSidebarItem === "pending-members"}
          onclick={() => handleSidebarClick("pending-members")}
        >
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
            <circle cx="8.5" cy="7" r="4" />
            <line x1="20" y1="8" x2="20" y2="14" />
            <line x1="23" y1="11" x2="17" y2="11" />
          </svg>
          <span>Yêu cầu tham gia</span>
          {#if pendingMembersCount > 0}
            <span class="pending-badge">{pendingMembersCount}</span>
          {/if}
        </button>
      {/if}

      <button
        class="nav-item"
        class:active={activeSidebarItem === "rules"}
        onclick={() => handleSidebarClick("rules")}
      >
        <img src="/rule_icon.svg" alt="" width="20" height="20" />
        <span>Nội quy</span>
      </button>
    </nav>
  </aside>

  <!-- Main Content -->
  <main class="mod-content">
    {#if activeSidebarItem === "queue"}
      <div class="queue-section">
        <div class="queue-header">
          <h1>Hàng chờ</h1>
          <div class="sort-options">
            <select bind:value={sortBy}>
              <option value="newest">Mới nhất</option>
              <option value="oldest">Cũ nhất</option>
              <option value="most-reported">Nhiều báo cáo nhất</option>
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
            Chưa kiểm duyệt
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "edited"}
            onclick={() => (activeQueueTab = "edited")}
          >
            Đã chỉnh sửa
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "removed"}
            onclick={() => (activeQueueTab = "removed")}
          >
            Đã xóa
          </button>
          <button
            class="tab-btn"
            class:active={activeQueueTab === "reported"}
            onclick={() => (activeQueueTab = "reported")}
          >
            Bị báo cáo
          </button>
        </div>

        <!-- Posts List -->
        <div class="posts-list">
          {#each filteredPosts() as post (post.id)}
            <div class="post-card">
              <div class="post-header">
                <div class="post-author">
                  {#if post.author?.avatar?.url}
                    <img
                      src={post.author.avatar.url}
                      alt={post.author.username}
                      class="author-avatar"
                    />
                  {:else}
                    <div class="author-avatar-placeholder">
                      {post.author?.username?.[0]?.toUpperCase() || "?"}
                    </div>
                  {/if}
                  <div class="author-info">
                    <span class="author-name"
                      >u/{post.author?.username || "[deleted]"}</span
                    >
                    <span class="post-time"
                      >{formatTime(post.created_at || post.createdAt)}</span
                    >
                  </div>
                </div>
                {#if post.reportCount}
                  <span class="report-badge">{post.reportCount} báo cáo</span>
                {/if}
              </div>

              <h3 class="post-title">{post.title}</h3>
              <p class="post-content">{post.content?.text || ""}</p>

              {#if post.reportReason}
                <div class="report-info">
                  <strong>Lý do báo cáo:</strong>
                  {post.reportReason}
                </div>
              {/if}

              {#if post.removedReason}
                <div class="removed-info">
                  <strong>Xóa bởi:</strong>
                  {post.removedBy} - {post.removedReason}
                </div>
              {/if}

              <div class="post-actions">
                <button
                  class="action-btn approve"
                  onclick={() => handleApprove(post.id)}
                >
                  Duyệt
                </button>
                <button
                  class="action-btn remove"
                  onclick={() => handleRemove(post.id)}
                >
                  Xóa
                </button>
              </div>
            </div>
          {:else}
            <div class="empty-state">
              <p>Không có bài viết trong hàng chờ này</p>
            </div>
          {/each}
        </div>
      </div>
    {:else if activeSidebarItem === "restricted"}
      <!-- Restricted Users Section -->
      <div class="restricted-section">
        <div class="restricted-header">
          <h1>Người dùng bị giới hạn</h1>
          <button
            class="action-btn-primary"
            onclick={activeRestrictedTab === "banned"
              ? handleOpenBanModal
              : handleOpenMuteModal}
          >
            {activeRestrictedTab === "banned"
              ? "Cấm người dùng"
              : "Tắt tiếng người dùng"}
          </button>
        </div>

        <!-- Tabs -->
        <div class="restricted-tabs">
          <button
            class="tab-btn"
            class:active={activeRestrictedTab === "banned"}
            onclick={() => (activeRestrictedTab = "banned")}
          >
            Bị cấm
          </button>
          <button
            class="tab-btn"
            class:active={activeRestrictedTab === "muted"}
            onclick={() => (activeRestrictedTab = "muted")}
          >
            Bị tắt tiếng
          </button>
        </div>

        <!-- Table -->
        <div class="restricted-table">
          {#if activeRestrictedTab === "banned"}
            <div class="table-header">
              <div class="col">TÊN NGƯỜI DÙNG</div>
              <div class="col">THỜI HẠN</div>
              <div class="col">NGÀY</div>
              <div class="col">LÝ DO</div>
              <div class="col">GHI CHÚ</div>
            </div>
          {:else}
            <div class="table-header muted">
              <div class="col">TÊN NGƯỜI DÙNG</div>
              <div class="col">Thời hạn</div>
              <div class="col">GHI CHÚ</div>
            </div>
          {/if}

          {#each filteredRestrictedUsers as user}
            <div class="table-row">
              {#if activeRestrictedTab === "banned"}
                <div class="col user-col">
                  <img
                    src={user.profile?.avatar?.url || "/user.jpg"}
                    alt=""
                    class="user-avatar"
                  />
                  <span>{user.username}</span>
                </div>
                <div class="col">Vĩnh viễn</div>
                <div class="col">-</div>
                <div class="col">-</div>
                <div class="col">
                  <button
                    class="unban-btn"
                    onclick={() => handleUnbanUser(user.id)}
                  >
                    Bỏ cấm
                  </button>
                </div>
              {:else}
                <div class="col user-col">
                  <img
                    src={user.profile?.avatar?.url || "/user.jpg"}
                    alt=""
                    class="user-avatar"
                  />
                  <span>{user.username}</span>
                </div>
                <div class="col">Vĩnh viễn</div>
                <div class="col">
                  <button
                    class="unban-btn"
                    onclick={() => handleUnmuteUser(user.id)}
                  >
                    Bỏ tắt tiếng
                  </button>
                </div>
              {/if}
            </div>
          {:else}
            <div class="empty-state">
              <p>
                Không có người dùng bị {activeRestrictedTab === "banned"
                  ? "cấm"
                  : "tắt tiếng"}
              </p>
            </div>
          {/each}
        </div>
      </div>
    {:else if activeSidebarItem === "members"}
      <!-- Mod & Members Section -->
      <div class="members-section">
        <div class="members-header">
          <h1>Quản trị & Thành viên</h1>
          <button
            class="action-btn-primary"
            onclick={() =>
              handleOpenInviteModal(
                activeMembersTab === "moderators" ? "mod" : "approved",
              )}
          >
            {activeMembersTab === "moderators"
              ? "Mời quản trị viên"
              : "Thêm người dùng được duyệt"}
          </button>
        </div>

        <!-- Tabs -->
        <div class="members-tabs">
          <button
            class="tab-btn"
            class:active={activeMembersTab === "moderators"}
            onclick={() => (activeMembersTab = "moderators")}
          >
            Quản trị viên
          </button>
          <button
            class="tab-btn"
            class:active={activeMembersTab === "approved"}
            onclick={() => (activeMembersTab = "approved")}
          >
            Người dùng được duyệt
          </button>
        </div>

        <!-- Table -->
        <div class="members-table">
          {#if activeMembersTab === "moderators"}
            <div class="table-header moderators">
              <div class="col">TÊN NGƯỜI DÙNG</div>
              <div class="col">QUYỀN HẠN</div>
              <div class="col">Có thể chỉnh sửa</div>
              <div class="col">THAM GIA</div>
              <div class="col"></div>
            </div>

            {#each moderators as mod}
              <div class="table-row moderators">
                <div class="col user-col">
                  <img
                    src={mod.avatar?.url || "/user.jpg"}
                    alt=""
                    class="user-avatar"
                  />
                  <span>{mod.username}</span>
                </div>
                <div class="col">Toàn bộ</div>
                <div class="col">Có</div>
                <div class="col joined-col">-</div>
                <div class="col actions-col">
                  <button
                    class="icon-btn edit"
                    onclick={() => handleEditMember(mod.user_id, "mod")}
                    title="Sửa quản trị viên"
                  >
                    <img
                      src="/write_icon.svg"
                      alt="Edit"
                      width="20"
                      height="20"
                    />
                  </button>
                  {#if isCreator()}
                    <button
                      class="icon-btn delete"
                      onclick={() => handleDeleteMember(mod.user_id, "mod")}
                      title="Xóa quản trị viên"
                    >
                      <img
                        src="/bin_icon.svg"
                        alt="Delete"
                        width="20"
                        height="20"
                      />
                    </button>
                  {/if}
                </div>
              </div>
            {:else}
              <div class="empty-state">
                <p>Chưa có quản trị viên</p>
              </div>
            {/each}
          {:else}
            <div class="table-header approved">
              <div class="col">TÊN NGƯỜI DÙNG</div>
              <div class="col">Ngày tham gia</div>
              <div class="col"></div>
            </div>

            {#if isLoadingApprovedMembers}
              <div class="loading-state">
                <p>Đang tải danh sách người dùng được duyệt...</p>
              </div>
            {:else}
              {#each approvedMembers as member}
                <div class="table-row approved">
                  <div class="col user-col">
                    <img
                      src={member.user?.avatar?.url || "/user.jpg"}
                      alt=""
                      class="user-avatar"
                    />
                    <span>{member.user?.username || "Không rõ"}</span>
                  </div>
                  <div class="col">{formatTime(member.created_at)}</div>
                  <div class="col actions-col">
                    <button
                      class="icon-btn delete"
                      onclick={() =>
                        handleRemoveApprovedUser(
                          member.user_id,
                          member.user?.username || "người dùng",
                        )}
                      title="Xóa người dùng"
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
                  <p>Chưa có người dùng được duyệt</p>
                </div>
              {/each}
            {/if}
          {/if}
        </div>
      </div>
    {:else if activeSidebarItem === "all-members"}
      <!-- All Members Section -->
      <div class="all-members-section">
        <div class="all-members-header">
          <h1>Tất cả thành viên ({allMembers.length})</h1>
          <button
            class="action-btn refresh"
            onclick={() => loadAllMembers()}
            disabled={isLoadingAllMembers}
          >
            {isLoadingAllMembers ? "Đang tải..." : "Làm mới"}
          </button>
        </div>

        <!-- Members Table -->
        <div class="all-members-table">
          <div class="table-header all-members">
            <div class="col">TÊN NGƯỜI DÙNG</div>
            <div class="col">VAI TRÒ</div>
            <div class="col">NGÀY THAM GIA</div>
            <div class="col"></div>
          </div>

          {#if isLoadingAllMembers}
            <div class="loading-state">
              <p>Đang tải danh sách thành viên...</p>
            </div>
          {:else}
            {#each allMembers as member}
              {@const isCreatorMember =
                member.user_id === community?.create_by_id}
              {@const isMod = moderators.some(
                (m) => m.user_id === member.user_id,
              )}
              <div class="table-row all-members">
                <div class="col user-col">
                  <img
                    src={member.user?.avatar?.url || "/user.jpg"}
                    alt=""
                    class="user-avatar"
                  />
                  <span>{member.user?.username || "Unknown"}</span>
                </div>
                <div class="col role-col">
                  {#if isCreatorMember}
                    <span class="role-badge creator">Người tạo</span>
                  {:else if isMod}
                    <span class="role-badge mod">Quản trị viên</span>
                  {:else}
                    <span class="role-badge member">Thành viên</span>
                  {/if}
                </div>
                <div class="col">
                  {member.created_at
                    ? new Date(member.created_at).toLocaleDateString("vi-VN")
                    : "-"}
                </div>
                <div class="col actions-col">
                  {#if !isCreatorMember && (isCreator() || !isMod)}
                    <button
                      class="icon-btn delete kick-btn"
                      onclick={() =>
                        handleKickMember(
                          member.user_id,
                          member.user?.username || "Unknown",
                        )}
                      disabled={isKickingMember === member.user_id}
                      title="Xóa khỏi cộng đồng"
                    >
                      {#if isKickingMember === member.user_id}
                        <span class="loading-spinner"></span>
                      {:else}
                        <img
                          src="/bin_icon.svg"
                          alt="Kick"
                          width="20"
                          height="20"
                        />
                      {/if}
                    </button>
                  {/if}
                </div>
              </div>
            {:else}
              <div class="empty-state">
                <p>Chưa có thành viên nào</p>
              </div>
            {/each}
          {/if}
        </div>
      </div>
    {:else if activeSidebarItem === "pending-members"}
      <!-- Pending Members Section -->
      <div class="pending-members-section">
        <div class="pending-members-header">
          <h1>Yêu cầu tham gia ({pendingMembersCount})</h1>
          <button
            class="action-btn refresh"
            onclick={() => loadPendingMembersList()}
            disabled={isLoadingPendingMembers}
          >
            {isLoadingPendingMembers ? "Đang tải..." : "Làm mới"}
          </button>
        </div>

        <p class="pending-members-description">
          Các yêu cầu tham gia cộng đồng đang chờ duyệt. Bạn có thể duyệt hoặc
          từ chối từng yêu cầu.
        </p>

        <!-- Pending Members Table -->
        <div class="pending-members-table">
          <div class="table-header pending-members">
            <div class="col">TÊN NGƯỜI DÙNG</div>
            <div class="col">NGÀY YÊU CẦU</div>
            <div class="col">HÀNH ĐỘNG</div>
          </div>

          {#if isLoadingPendingMembers}
            <div class="loading-state">
              <p>Đang tải danh sách yêu cầu...</p>
            </div>
          {:else}
            {#each pendingMembers as member}
              <div class="table-row pending-members">
                <div class="col user-col">
                  <button
                    class="user-link"
                    onclick={() => push(`/user/${member.user?.username}`)}
                  >
                    <img
                      src={member.user?.avatar?.url || "/user.jpg"}
                      alt={member.user?.username || "User"}
                      class="user-avatar"
                    />
                    <span class="username"
                      >u/{member.user?.username || "Unknown"}</span
                    >
                  </button>
                </div>
                <div class="col">
                  {formatTime(member.created_at)}
                </div>
                <div class="col actions-col">
                  <button
                    class="action-btn approve"
                    onclick={() =>
                      handleApproveMember(
                        member.id,
                        member.user?.username || "Unknown",
                      )}
                    disabled={isProcessingMember === member.id}
                    title="Duyệt thành viên"
                  >
                    {#if isProcessingMember === member.id}
                      Đang xử lý...
                    {:else}
                      ✓ Duyệt
                    {/if}
                  </button>
                  <button
                    class="action-btn reject"
                    onclick={() =>
                      handleRejectMember(
                        member.id,
                        member.user?.username || "Unknown",
                      )}
                    disabled={isProcessingMember === member.id}
                    title="Từ chối yêu cầu"
                  >
                    ✕ Từ chối
                  </button>
                </div>
              </div>
            {:else}
              <div class="empty-state">
                <p>Không có yêu cầu nào đang chờ duyệt</p>
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
            <h1>Nội quy cộng đồng</h1>
            <button class="create-rule-btn" onclick={handleCreateRule}>
              Tạo nội quy
            </button>
          </div>

          <div class="rules-table">
            <div class="rules-table-header">
              <div class="col-name">TÊN</div>
              <div class="col-created">TẠO NGÀY</div>
            </div>

            {#each communityRules as rule, index}
              <div class="rule-row" onclick={() => handleEditRule(index)}>
                <div class="rule-info">
                  <span class="rule-number">{index + 1}</span>
                  <div class="rule-details">
                    <h3 class="rule-name">{rule.title}</h3>
                    <p class="rule-desc" style="color: #666;">
                      {rule.description}
                    </p>
                  </div>
                </div>
                <div class="rule-actions">
                  <button
                    class="icon-btn edit"
                    onclick={(e) => {
                      e.stopPropagation();
                      handleEditRule(index);
                    }}
                    title="Sửa quy tắc"
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
                      handleDeleteRule(index);
                    }}
                    title="Xóa quy tắc"
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
                    title="Thêm tùy chọn"
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
                <p>Chưa có nội quy. Tạo nội quy đầu tiên của bạn!</p>
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
                <h2>Đặt tên và mô tả nội quy của bạn</h2>
                <p class="sub">
                  Nội quy đặt ra kỳ vọng cho thành viên và người ghé thăm cộng
                  đồng
                </p>
              </div>
            </div>
          </div>

          <form class="rule-form">
            <div class="form-row">
              <label>Tên nội quy<span class="required">*</span></label>
              <div class="pill-input">
                <input
                  type="text"
                  placeholder="Tên nội quy"
                  maxlength="100"
                  bind:value={ruleName}
                />
                <span class="char-count">{100 - ruleName.length}</span>
              </div>
            </div>

            <div class="form-row">
              <label>Mô tả<span class="required">*</span></label>
              <div class="pill-input textarea">
                <textarea
                  placeholder="Mô tả"
                  maxlength="500"
                  bind:value={ruleDescription}
                ></textarea>
                <span class="char-count">{500 - ruleDescription.length}</span>
              </div>
            </div>

            <div class="form-row save-row">
              <button
                class="save-btn"
                class:enabled={isRuleFormValid}
                disabled={!isRuleFormValid}
                onclick={handleSaveRule}
                type="button"
              >
                Lưu
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
      <h2>Cấm người dùng</h2>

      <div class="form-group">
        <div class="search-input-wrapper-with-suggestions">
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
              placeholder="Tìm kiếm thành viên..."
              bind:value={banUsername}
              class="search-input"
              oninput={() => {
                showUserSuggestions = banUsername.trim().length > 0;
                selectedSuggestionIndex = -1;
              }}
              onkeydown={handleUserInputKeydown}
              onfocus={handleUserInputFocus}
              onblur={handleUserInputBlur}
            />
          </div>
          {#if showUserSuggestions && filteredMemberSuggestions().length > 0}
            <div class="user-suggestions">
              {#each filteredMemberSuggestions() as member, index}
                <button
                  type="button"
                  class="suggestion-item"
                  class:selected={index === selectedSuggestionIndex}
                  onclick={() =>
                    handleSelectUserSuggestion(member.user?.username || "")}
                >
                  <img
                    src={member.user?.avatar?.url || "/user.jpg"}
                    alt=""
                    class="suggestion-avatar"
                  />
                  <span class="suggestion-username"
                    >{member.user?.username}</span
                  >
                </button>
              {/each}
            </div>
          {:else if showUserSuggestions && banUsername.trim() && filteredMemberSuggestions().length === 0}
            <div class="user-suggestions">
              <div class="no-suggestions">Không tìm thấy thành viên</div>
            </div>
          {/if}
        </div>
        <p class="hint">Nhập tên để tìm thành viên trong cộng đồng</p>
      </div>

      <div class="form-group">
        <select bind:value={banRule} class="modal-select">
          <option value="">Rules</option>
          {#each communityRules as rule}
            <option value={rule.title}>{rule.title}</option>
          {/each}
        </select>
      </div>

      <div class="form-group">
        <select bind:value={banDuration} class="modal-select">
          <option value="">Thời hạn</option>
          <option value="1h">1 giờ</option>
          <option value="24h">24 giờ</option>
          <option value="7d">7 ngày</option>
          <option value="30d">30 ngày</option>
          <option value="permanent">Vĩnh viễn</option>
        </select>
      </div>

      <div class="form-group">
        <textarea
          placeholder="Lý do"
          bind:value={banReason}
          class="modal-textarea"
        ></textarea>
      </div>

      <div class="form-group">
        <textarea
          placeholder="Ghi chú"
          bind:value={banNote}
          class="modal-textarea"
        ></textarea>
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseBanModal}> Hủy </button>
        <button class="btn-danger" onclick={handleBanUser}>Cấm</button>
      </div>
    </div>
  </div>
{/if}

<!-- Mute User Modal -->
{#if showMuteModal}
  <div class="modal-overlay" onclick={handleCloseMuteModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>Tắt tiếng người dùng</h2>

      <div class="form-group">
        <div class="search-input-wrapper-with-suggestions">
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
              placeholder="Tìm kiếm thành viên..."
              bind:value={banUsername}
              class="search-input"
              oninput={() => {
                showUserSuggestions = banUsername.trim().length > 0;
                selectedSuggestionIndex = -1;
              }}
              onkeydown={handleUserInputKeydown}
              onfocus={handleUserInputFocus}
              onblur={handleUserInputBlur}
            />
          </div>
          {#if showUserSuggestions && filteredMemberSuggestions().length > 0}
            <div class="user-suggestions">
              {#each filteredMemberSuggestions() as member, index}
                <button
                  type="button"
                  class="suggestion-item"
                  class:selected={index === selectedSuggestionIndex}
                  onclick={() =>
                    handleSelectUserSuggestion(member.user?.username || "")}
                >
                  <img
                    src={member.user?.avatar?.url || "/user.jpg"}
                    alt=""
                    class="suggestion-avatar"
                  />
                  <span class="suggestion-username"
                    >{member.user?.username}</span
                  >
                </button>
              {/each}
            </div>
          {:else if showUserSuggestions && banUsername.trim() && filteredMemberSuggestions().length === 0}
            <div class="user-suggestions">
              <div class="no-suggestions">Không tìm thấy thành viên</div>
            </div>
          {/if}
        </div>
        <p class="hint">Nhập tên để tìm thành viên trong cộng đồng</p>
      </div>

      <div class="form-group">
        <select bind:value={banRule} class="modal-select">
          <option value="">Rules</option>
          {#each communityRules as rule}
            <option value={rule.title}>{rule.title}</option>
          {/each}
        </select>
      </div>

      <div class="form-group">
        <select bind:value={banDuration} class="modal-select">
          <option value="">Thời hạn</option>
          <option value="1h">1 giờ</option>
          <option value="24h">24 giờ</option>
          <option value="7d">7 ngày</option>
          <option value="30d">30 ngày</option>
          <option value="permanent">Vĩnh viễn</option>
        </select>
      </div>

      <div class="form-group">
        <textarea
          placeholder="Lý do"
          bind:value={banReason}
          class="modal-textarea"
        ></textarea>
      </div>

      <div class="form-group">
        <textarea
          placeholder="Ghi chú"
          bind:value={banNote}
          class="modal-textarea"
        ></textarea>
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseMuteModal}> Hủy </button>
        <button class="btn-danger" onclick={handleMuteUser}>Tắt tiếng</button>
      </div>
    </div>
  </div>
{/if}

<!-- Invite/Add User Modal -->
{#if showInviteModal}
  <div class="modal-overlay" onclick={handleCloseInviteModal}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>
        {inviteType === "mod"
          ? "Mời quản trị viên"
          : "Thêm người dùng được duyệt"}
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
            placeholder="Tìm kiếm người dùng"
            bind:value={inviteUsername}
            class="search-input"
          />
        </div>
        <p class="hint">Nhập tên người dùng để tìm</p>
      </div>

      {#if inviteType === "mod"}
        <div class="form-group">
          <label>Quyền hạn</label>
          <select bind:value={invitePermission} class="permission-select">
            <option value="Everything">Toàn bộ</option>
            <option value="Manage Posts & Comments"
              >Quản lý bài viết & bình luận</option
            >
            <option value="Manage Users">Quản lý người dùng</option>
            <option value="Manage Settings">Quản lý cài đặt</option>
          </select>
        </div>

        <div class="form-group">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={inviteCanEdit} />
            <span>Bạn có thể chỉnh sửa quản trị viên này</span>
          </label>
        </div>
      {/if}

      <div class="modal-actions">
        <button class="btn-cancel" onclick={handleCloseInviteModal}>
          Hủy
        </button>
        <button class="action-btn-primary" onclick={handleInviteUser}>
          {inviteType === "mod" ? "Mời" : "Thêm"}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Confirm Modals -->
<ConfirmModal
  show={showKickConfirm}
  title="Xác nhận xóa thành viên"
  message={`Bạn có chắc chắn muốn xóa ${confirmTargetUser?.username || ""} khỏi cộng đồng?`}
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmKickMember}
  onCancel={() => {
    showKickConfirm = false;
    confirmTargetUser = null;
  }}
/>

<ConfirmModal
  show={showRejectConfirm}
  title="Xác nhận từ chối"
  message={`Bạn có chắc chắn muốn từ chối yêu cầu của ${confirmTargetUser?.username || ""}?`}
  confirmText="Từ chối"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmRejectMember}
  onCancel={() => {
    showRejectConfirm = false;
    confirmTargetUser = null;
  }}
/>

<ConfirmModal
  show={showRemoveApprovedConfirm}
  title="Xác nhận xóa người dùng"
  message={`Bạn có chắc chắn muốn xóa ${confirmTargetUser?.username || ""} khỏi danh sách người dùng được duyệt?`}
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmRemoveApprovedUser}
  onCancel={() => {
    showRemoveApprovedConfirm = false;
    confirmTargetUser = null;
  }}
/>

<ConfirmModal
  show={showDeleteRuleConfirm}
  title="Xác nhận xóa quy tắc"
  message="Bạn có chắc chắn muốn xóa quy tắc này?"
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDeleteRule}
  onCancel={() => {
    showDeleteRuleConfirm = false;
    confirmTargetIndex = null;
  }}
/>

<ConfirmModal
  show={showUnbanConfirm}
  title="Xác nhận bỏ cấm"
  message="Bạn có chắc chắn muốn bỏ cấm người dùng này?"
  confirmText="Bỏ cấm"
  cancelText="Hủy"
  confirmVariant="primary"
  onConfirm={confirmUnbanUser}
  onCancel={() => {
    showUnbanConfirm = false;
    confirmTargetUser = null;
  }}
/>

<ConfirmModal
  show={showUnmuteConfirm}
  title="Xác nhận bỏ tắt tiếng"
  message="Bạn có chắc chắn muốn bỏ tắt tiếng người dùng này?"
  confirmText="Bỏ tắt tiếng"
  cancelText="Hủy"
  confirmVariant="primary"
  onConfirm={confirmUnmuteUser}
  onCancel={() => {
    showUnmuteConfirm = false;
    confirmTargetUser = null;
  }}
/>

<ConfirmModal
  show={showDeleteMemberConfirm}
  title={confirmMemberType === "mod"
    ? "Xác nhận xóa quản trị viên"
    : "Xác nhận xóa người dùng"}
  message={`Bạn có chắc chắn muốn xóa ${confirmMemberType === "mod" ? "quản trị viên" : "người dùng được duyệt"} này?`}
  confirmText="Xóa"
  cancelText="Hủy"
  confirmVariant="danger"
  onConfirm={confirmDeleteMember}
  onCancel={() => {
    showDeleteMemberConfirm = false;
    confirmTargetUser = null;
  }}
/>

<!-- Remove Post Modal -->
{#if showRemovePostModal}
  <div class="confirm-overlay" onclick={() => (showRemovePostModal = false)}>
    <div class="confirm-modal" onclick={(e) => e.stopPropagation()}>
      <div class="confirm-header">
        <h3>Xóa bài viết</h3>
      </div>
      <div class="confirm-body">
        <p>Lý do xóa bài viết (tùy chọn):</p>
        <input
          type="text"
          bind:value={removePostReason}
          placeholder="Nhập lý do..."
          class="reason-input"
        />
      </div>
      <div class="confirm-actions">
        <button
          class="btn-cancel"
          onclick={() => {
            showRemovePostModal = false;
            removePostId = null;
          }}
        >
          Hủy
        </button>
        <button class="btn-confirm danger" onclick={confirmRemovePost}>
          Xóa bài viết
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

  .search-input-wrapper-with-suggestions {
    position: relative;
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

  .user-suggestions {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 100;
    max-height: 200px;
    overflow-y: auto;
    margin-top: 4px;
  }

  .suggestion-item {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 10px 16px;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    transition: background 0.15s;
  }

  .suggestion-item:hover,
  .suggestion-item.selected {
    background: #f6f7f8;
  }

  .suggestion-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
  }

  .suggestion-username {
    font-size: 14px;
    color: #1c1c1c;
    font-weight: 500;
  }

  .no-suggestions {
    padding: 12px 16px;
    color: #7c7c7c;
    font-size: 14px;
    text-align: center;
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

  /* All Members Section */
  .all-members-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .all-members-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .all-members-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .action-btn.refresh {
    background: #f6f7f8;
    color: #1c1c1c;
    border: 1px solid #edeff1;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn.refresh:hover:not(:disabled) {
    background: #edeff1;
  }

  .action-btn.refresh:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .all-members-table {
    border: 1px solid #edeff1;
    border-radius: 8px;
    overflow: hidden;
  }

  .all-members-table .table-header.all-members {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 100px;
    background: #f6f7f8;
    padding: 12px 16px;
    font-size: 12px;
    font-weight: 700;
    color: var(--grayfont);
    text-transform: uppercase;
  }

  .all-members-table .table-row.all-members {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 100px;
    padding: 16px;
    border-top: 1px solid #edeff1;
    align-items: center;
    font-size: 14px;
    color: #1c1c1c;
  }

  .role-col {
    display: flex;
    align-items: center;
  }

  .role-badge {
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
  }

  .role-badge.creator {
    background: #fff3cd;
    color: #856404;
  }

  .role-badge.mod {
    background: #d4edda;
    color: #155724;
  }

  .role-badge.member {
    background: #e9ecef;
    color: #495057;
  }

  .kick-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .loading-spinner {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid #ddd;
    border-top-color: #666;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .loading-state {
    padding: 40px;
    text-align: center;
    color: #7c7c7c;
  }

  /* Pending Members Section */
  .pending-members-section {
    background: white;
    border-radius: 8px;
    padding: 24px 32px;
  }

  .pending-members-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .pending-members-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .pending-members-description {
    color: #7c7c7c;
    font-size: 14px;
    margin-bottom: 24px;
  }

  .pending-members-table {
    border: 1px solid #edeff1;
    border-radius: 8px;
    overflow: hidden;
  }

  .pending-members-table .table-header.pending-members {
    display: grid;
    grid-template-columns: 2fr 1fr 200px;
    background: #f6f7f8;
    padding: 12px 16px;
    font-size: 12px;
    font-weight: 700;
    color: var(--grayfont);
    text-transform: uppercase;
  }

  .pending-members-table .table-row.pending-members {
    display: grid;
    grid-template-columns: 2fr 1fr 200px;
    padding: 16px;
    border-top: 1px solid #edeff1;
    align-items: center;
    font-size: 14px;
    color: #1c1c1c;
  }

  .action-btn.approve {
    background: var(--blue--);
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn.approve:hover:not(:disabled) {
    background: #00008b;
    filter: brightness(0.85);
  }

  .action-btn.approve:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .action-btn.reject {
    background: #f6f7f8;
    color: #ff4500;
    border: 1px solid #ff4500;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .action-btn.reject:hover:not(:disabled) {
    background: #fff5f5;
  }

  .action-btn.reject:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .pending-badge {
    background: #ff4500;
    color: white;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 12px;
    font-weight: 600;
    margin-left: 8px;
  }

  /* Confirm Modal Styles */
  .confirm-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10000;
  }

  .confirm-modal {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 400px;
    padding: 24px;
  }

  .confirm-header h3 {
    margin: 0 0 16px 0;
    font-size: 18px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .confirm-body {
    margin-bottom: 24px;
  }

  .confirm-body p {
    margin: 0 0 12px 0;
    font-size: 14px;
    color: #576f76;
  }

  .reason-input {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
  }

  .confirm-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
  }

  .btn-cancel,
  .btn-confirm {
    padding: 10px 20px;
    border-radius: 9999px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    border: none;
  }

  .btn-cancel {
    background: rgba(214, 216, 222, 0.5);
    color: #1c1c1c;
  }

  .btn-cancel:hover {
    background: rgba(214, 216, 222, 0.7);
  }

  .btn-confirm.danger {
    background: #ff4500;
    color: white;
  }

  .btn-confirm.danger:hover {
    background: #e03d00;
  }
</style>
