export interface CommentData {
  id: string;
  author: string;
  avatar?: string;
  time: string;
  content: string;
  upvotes: number;
  downvotes: number;
  replies?: CommentData[];
  isCollapsed?: boolean;
  isEditing?: boolean;
  depth?: number;
}

export type CommentSortType = 'top' | 'newest' | 'oldest' | 'controversial';
