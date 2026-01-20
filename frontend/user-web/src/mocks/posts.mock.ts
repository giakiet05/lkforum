export interface Post {
  id: string;
  authorId: string;
  authorUsername: string;
  authorAvatar?: string;
  communityId: string;
  communityName: string;
  type: 'text' | 'poll' | 'video' | 'image';
  title: string;
  content: {
    text?: string;
    poll?: {
      question: string;
      options: Array<{
        id: string;
        text: string;
        votes: number;
      }>;
      totalVotes: number;
      expiresAt?: string;
      allowMultiple: boolean;
    };
    video?: {
      thumbnail: string;
      url: string;
    };
    images?: Array<{
      id: string;
      url: string;
      uploadedAt: string;
    }>;
  };
  votesCount: {
    up: number;
    down: number;
  };
  commentsCount: number;
  createdAt: string;
  updatedAt?: string;
}

export const mockPosts: Post[] = [
  {
    id: '1',
    authorId: 'user1',
    authorUsername: 'johndoe',
    authorAvatar: 'https://i.pravatar.cc/150?img=1',
    communityId: 'tech',
    communityName: 'Technology',
    type: 'text',
    title: 'Getting Started with TypeScript',
    content: {
      text: 'TypeScript is a powerful superset of JavaScript that adds static types...',
    },
    votesCount: {
      up: 150,
      down: 10
    },
    commentsCount: 25,
    createdAt: '2025-10-19T08:00:00Z'
  },
  {
    id: '2',
    authorId: 'user2',
    authorUsername: 'techie',
    authorAvatar: 'https://i.pravatar.cc/150?img=2',
    communityId: 'programming',
    communityName: 'Programming',
    type: 'poll',
    title: 'What\'s your favorite programming language?',
    content: {
      poll: {
        question: 'Which programming language do you use the most?',
        options: [
          { id: 'opt1', text: 'JavaScript', votes: 120 },
          { id: 'opt2', text: 'Python', votes: 150 },
          { id: 'opt3', text: 'Java', votes: 80 },
          { id: 'opt4', text: 'C++', votes: 60 }
        ],
        totalVotes: 410,
        allowMultiple: false
      }
    },
    votesCount: {
      up: 200,
      down: 15
    },
    commentsCount: 45,
    createdAt: '2025-10-18T15:30:00Z'
  },
  {
    id: '3',
    authorId: 'user3',
    authorUsername: 'photoexpert',
    authorAvatar: 'https://i.pravatar.cc/150?img=3',
    communityId: 'photography',
    communityName: 'Photography',
    type: 'image',
    title: 'Sunset at the beach',
    content: {
      text: 'Captured this beautiful sunset at Miami Beach',
      images: [
        {
          id: 'img1',
          url: 'https://picsum.photos/800/600',
          uploadedAt: '2025-10-18T20:00:00Z'
        },
        {
          id: 'img2',
          url: 'https://picsum.photos/800/601',
          uploadedAt: '2025-10-18T20:00:00Z'
        }
      ]
    },
    votesCount: {
      up: 500,
      down: 20
    },
    commentsCount: 75,
    createdAt: '2025-10-18T20:00:00Z'
  }
];