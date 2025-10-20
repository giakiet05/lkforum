export type PollOption = {
  id: number;
  text: string;
  votes: number;
};

export type PostData = {
  id: string;
  type: 'text' | 'image' | 'video' | 'poll';
  community: string;
  author: string;
  time: string;
  title: string;
  upvotes: number;
  downvotes: number;
  commentsCount: number;
  content?: string; // For text posts
  images?: string[]; // For image posts
  videoUrl?: string; // For video posts
  thumbnailUrl?: string; // For video posts
  poll?: {
    question: string;
    options: PollOption[];
    multipleChoice: boolean;
    totalVotes: number;
  };
};
