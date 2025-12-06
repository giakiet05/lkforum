# Script to KEEP or UNDO comment feature changes
# Usage: .\restore-comment-changes.ps1 -Action [Keep|Undo]

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("Keep", "Undo")]
    [string]$Action
)

$ErrorActionPreference = "Stop"

# Backup copies of ORIGINAL code (before changes)
$backups = @{
    "comment-dto-frontend" = @{
        path = "frontend\user-web\src\dtos\comment-dto.ts"
        original = @"
// --- Request DTOs ---

export interface CreateCommentRequest {
    user_id: string;
    post_id: string;
    parent_id?: string;
    content: string;
}
"@
        current = @"
// --- Request DTOs ---

export interface CreateCommentRequest {
    post_id: string;
    parent_id?: string;
    content: string;
}
"@
    }
    
    "comment-service" = @{
        path = "frontend\user-web\src\services\comment-service.ts"
        original = @"
/**
 * Create a new comment
 */
export async function createComment(data: CreateCommentRequest): Promise<CommentResponse> {
    const res = await authenticatedFetch("/api/comments", {
        method: "POST",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}
"@
        current = @"
/**
 * Create a new comment or reply
 */
export async function createComment(data: CreateCommentRequest): Promise<CommentResponse> {
    const res = await authenticatedFetch("/api/comments", {
        method: "POST",
        body: JSON.stringify(data),
    });

    return await handleApiResponse(res);
}
"@
    }
    
    "comment-section" = @{
        path = "frontend\user-web\src\components\CommentSection.svelte"
        original = @"
  const submitComment = async () => {
    if (!newCommentContent.trim()) return;
    if (!currentUser) {
      alert("Please login to comment");
      return;
    }

    try {
      isSubmitting = true;
      await createComment({
        user_id: currentUser.id,
        username: currentUser.username,
        user_avatar: currentUser.profile?.avatar?.url || "",
        post_id: postId,
        content: newCommentContent,
      });

      await loadComments();

      newCommentContent = "";
      selectedImage = null;
      imagePreview = null;
    } catch (error) {
      console.error("Failed to submit comment:", error);
      alert("Failed to post comment. Please try again.");
    } finally {
      isSubmitting = false;
    }
  };
"@
        current = @"
  const submitComment = async () => {
    if (!newCommentContent.trim()) return;
    if (!currentUser) {
      alert("Please login to comment");
      return;
    }

    try {
      isSubmitting = true;
      await createComment({
        post_id: postId,
        content: newCommentContent,
      });

      await loadComments();

      newCommentContent = "";
      selectedImage = null;
      imagePreview = null;
    } catch (error) {
      console.error("Failed to submit comment:", error);
      alert("Failed to post comment. Please try again.");
    } finally {
      isSubmitting = false;
    }
  };
"@
    }
    
    "comment-component-import" = @{
        path = "frontend\user-web\src\components\Comment.svelte"
        original = @"
<script lang="ts">
  import type { CommentResponse } from "../dtos/comment-dto";
  import CommentComponent from "./Comment.svelte";
  import { deleteComment } from "../services/comment-service";
  import { authStore } from "../stores/auth-store";
"@
        current = @"
<script lang="ts">
  import type { CommentResponse } from "../dtos/comment-dto";
  import CommentComponent from "./Comment.svelte";
  import { deleteComment, createComment } from "../services/comment-service";
  import { authStore } from "../stores/auth-store";
"@
    }
    
    "comment-component-reply" = @{
        path = "frontend\user-web\src\components\Comment.svelte"
        original = @"
  const submitReply = () => {
    if (replyContent.trim()) {
      // Mock reply - in real app, would add to parent state
      console.log("Replying to comment:", comment.id, replyContent);
      if (replyImage) {
        console.log("With image:", replyImage.name);
      }
      replyContent = "";
      replyImage = null;
      replyImagePreview = null;
      showReplyBox = false;
    }
  };
"@
        current = @"
  const submitReply = async () => {
    if (!replyContent.trim()) return;
    if (!currentUser) {
      alert("Please login to reply");
      return;
    }

    try {
      await createComment({
        post_id: comment.post_id,
        parent_id: comment.id,
        content: replyContent,
      });

      // Clear reply form
      replyContent = "";
      replyImage = null;
      replyImagePreview = null;
      showReplyBox = false;

      // Reload comments to show the new reply
      if (onUpdate) {
        onUpdate();
      }
    } catch (error) {
      console.error("Failed to submit reply:", error);
      alert("Failed to post reply. Please try again.");
    }
  };
"@
    }
    
    "backend-dto" = @{
        path = "backend\internal\dto\comment_dto.go"
        original = @"
type CreateCommentRequest struct {
	UserID   string  ``json:"user_id"``
	PostID   string  ``json:"post_id"``
	ParentID *string ``json:"parent_id,omitempty"``
	Content  string  ``json:"content"``
}
"@
        current = @"
type CreateCommentRequest struct {
	PostID   string  ``json:"post_id" binding:"required"``
	ParentID *string ``json:"parent_id,omitempty"``
	Content  string  ``json:"content" binding:"required"``
}
"@
    }
}

function Write-ColorOutput {
    param($Message, $Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

Write-ColorOutput "`n========================================" "Cyan"
Write-ColorOutput "  Comment Feature Change Manager" "Cyan"
Write-ColorOutput "========================================`n" "Cyan"

if ($Action -eq "Keep") {
    Write-ColorOutput "✅ KEEPING all changes (no action needed)" "Green"
    Write-ColorOutput "`nThe following changes are currently active:" "Yellow"
    Write-ColorOutput "  • Frontend DTOs updated (removed user_id)" "White"
    Write-ColorOutput "  • Comment service updated" "White"
    Write-ColorOutput "  • CommentSection component updated" "White"
    Write-ColorOutput "  • Comment reply functionality implemented" "White"
    Write-ColorOutput "  • Backend DTO validation added" "White"
    Write-ColorOutput "`n✨ Comment feature is fully functional!" "Green"
    Write-ColorOutput "`nBackend needs to be restarted if not already running." "Yellow"
    
} elseif ($Action -eq "Undo") {
    Write-ColorOutput "⏪ UNDOING all changes..." "Yellow"
    
    $changesApplied = 0
    
    foreach ($key in $backups.Keys) {
        $backup = $backups[$key]
        $filePath = $backup.path
        
        if (Test-Path $filePath) {
            Write-ColorOutput "  Restoring: $filePath" "Gray"
            
            $content = Get-Content $filePath -Raw
            
            # Replace current with original
            $newContent = $content -replace [regex]::Escape($backup.current), $backup.original
            
            if ($newContent -ne $content) {
                Set-Content -Path $filePath -Value $newContent -NoNewline
                $changesApplied++
                Write-ColorOutput "    ✓ Restored" "Green"
            } else {
                Write-ColorOutput "    ⚠ No changes found (may already be original)" "Yellow"
            }
        } else {
            Write-ColorOutput "    ⚠ File not found: $filePath" "Red"
        }
    }
    
    Write-ColorOutput "`n✅ Rollback complete! Applied $changesApplied changes." "Green"
    Write-ColorOutput "`n⚠️  IMPORTANT: Restart backend server for changes to take effect!" "Yellow"
    Write-ColorOutput "   Run: cd backend; go run main.go" "Gray"
}

Write-ColorOutput "`n========================================`n" "Cyan"
