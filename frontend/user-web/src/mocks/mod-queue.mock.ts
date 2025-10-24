export interface QueuePost {
  id: string;
  title: string;
  author: string;
  authorAvatar?: string;
  community: string;
  content: string;
  type: "text" | "image" | "link";
  createdAt: string;
  queueType: "unmoderated" | "edited" | "removed" | "reported";
  reportReason?: string;
  reportCount?: number;
  editedAt?: string;
  removedBy?: string;
  removedReason?: string;
}

// Mock data for moderation queue
export const mockQueuePosts: QueuePost[] = [
  {
    id: "1",
    title: "New anime recommendation thread",
    author: "anime_lover123",
    authorAvatar: "https://i.pravatar.cc/150?img=1",
    community: "anime",
    content: "Looking for some great anime to watch this season...",
    type: "text",
    createdAt: "2025-10-24T10:30:00Z",
    queueType: "unmoderated",
  },
  {
    id: "2",
    title: "Check out this funny meme",
    author: "meme_master",
    authorAvatar: "https://i.pravatar.cc/150?img=2",
    community: "anime",
    content: "https://example.com/funny-meme.jpg",
    type: "image",
    createdAt: "2025-10-24T09:15:00Z",
    queueType: "unmoderated",
  },
  {
    id: "3",
    title: "Updated post about best practices",
    author: "dev_guru",
    authorAvatar: "https://i.pravatar.cc/150?img=3",
    community: "anime",
    content: "I've updated my original post with new information...",
    type: "text",
    createdAt: "2025-10-23T15:20:00Z",
    queueType: "edited",
    editedAt: "2025-10-24T08:00:00Z",
  },
  {
    id: "4",
    title: "Spam post about products",
    author: "spammer99",
    authorAvatar: "https://i.pravatar.cc/150?img=4",
    community: "anime",
    content: "Buy our amazing products now!!!",
    type: "link",
    createdAt: "2025-10-24T07:45:00Z",
    queueType: "removed",
    removedBy: "mod_anime",
    removedReason: "Spam",
  },
  {
    id: "5",
    title: "Controversial opinion on recent episode",
    author: "hot_takes",
    authorAvatar: "https://i.pravatar.cc/150?img=5",
    community: "anime",
    content: "I think this episode was terrible and here's why...",
    type: "text",
    createdAt: "2025-10-24T06:30:00Z",
    queueType: "reported",
    reportReason: "Harassment",
    reportCount: 3,
  },
  {
    id: "6",
    title: "Toxic comment detected",
    author: "toxic_user",
    authorAvatar: "https://i.pravatar.cc/150?img=6",
    community: "anime",
    content: "This is a very rude and offensive post...",
    type: "text",
    createdAt: "2025-10-24T05:15:00Z",
    queueType: "reported",
    reportReason: "Hate speech",
    reportCount: 5,
  },
];
