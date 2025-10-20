export interface Comment {
  id: string;
  postId: string;
  authorId: string;
  authorUsername: string;
  authorAvatar?: string;
  content: string;
  votesCount: {
    up: number;
    down: number;
  };
  replies: Comment[];
  createdAt: string;
  updatedAt?: string;
}

export const mockComments: Comment[] = [
  {
    id: 'comment1',
    postId: '1',
    authorId: 'user4',
    authorUsername: 'commentor1',
    authorAvatar: 'https://i.pravatar.cc/150?img=7',
    content: 'This is a great introduction to TypeScript! Really helpful for beginners.',
    votesCount: {
      up: 125,
      down: 8
    },
    replies: [
      {
        id: 'reply1-1',
        postId: '1',
        authorId: 'user5',
        authorUsername: 'replier1',
        authorAvatar: 'https://i.pravatar.cc/150?img=8',
        content: 'Agreed! The examples are very clear.',
        votesCount: {
          up: 43,
          down: 2
        },
        replies: [
          {
            id: 'reply1-1-1',
            postId: '1',
            authorId: 'user4',
            authorUsername: 'commentor1',
            authorAvatar: 'https://i.pravatar.cc/150?img=7',
            content: 'Thank you! I tried to make them as practical as possible.',
            votesCount: {
              up: 28,
              down: 0
            },
            replies: [
              {
                id: 'reply1-1-1-1',
                postId: '1',
                authorId: 'user7',
                authorUsername: 'deepreplier',
                authorAvatar: 'https://i.pravatar.cc/150?img=10',
                content: 'This thread is getting deep! But it\'s great to see the nested structure working.',
                votesCount: {
                  up: 15,
                  down: 1
                },
                replies: [],
                createdAt: '2025-10-19T09:15:00Z'
              }
            ],
            createdAt: '2025-10-19T09:00:00Z'
          }
        ],
        createdAt: '2025-10-19T08:30:00Z'
      },
      {
        id: 'reply1-2',
        postId: '1',
        authorId: 'user8',
        authorUsername: 'tsexpert',
        authorAvatar: 'https://i.pravatar.cc/150?img=11',
        content: 'One thing to add: TypeScript\'s type inference is really powerful. You don\'t always need to explicitly type everything!',
        votesCount: {
          up: 67,
          down: 3
        },
        replies: [],
        createdAt: '2025-10-19T10:00:00Z'
      }
    ],
    createdAt: '2025-10-19T08:15:00Z'
  },
  {
    id: 'comment2',
    postId: '1',
    authorId: 'user6',
    authorUsername: 'commentor2',
    authorAvatar: 'https://i.pravatar.cc/150?img=9',
    content: 'Could you explain more about TypeScript interfaces? I\'m still a bit confused about when to use interfaces vs types.',
    votesCount: {
      up: 89,
      down: 4
    },
    replies: [
      {
        id: 'reply2-1',
        postId: '1',
        authorId: 'user9',
        authorUsername: 'typescript_guru',
        authorAvatar: 'https://i.pravatar.cc/150?img=12',
        content: 'Great question! Interfaces are better for object shapes and can be extended. Types are more flexible and can represent unions, intersections, primitives, etc.',
        votesCount: {
          up: 112,
          down: 2
        },
        replies: [
          {
            id: 'reply2-1-1',
            postId: '1',
            authorId: 'user6',
            authorUsername: 'commentor2',
            authorAvatar: 'https://i.pravatar.cc/150?img=9',
            content: 'That makes sense! So interfaces for objects, types for more complex scenarios?',
            votesCount: {
              up: 34,
              down: 1
            },
            replies: [],
            createdAt: '2025-10-19T09:45:00Z'
          }
        ],
        createdAt: '2025-10-19T09:30:00Z'
      }
    ],
    createdAt: '2025-10-19T09:00:00Z'
  },
  {
    id: 'comment3',
    postId: '1',
    authorId: 'user10',
    authorUsername: 'controversial_user',
    authorAvatar: 'https://i.pravatar.cc/150?img=13',
    content: 'I don\'t think TypeScript is necessary for most projects. JavaScript is fine.',
    votesCount: {
      up: 45,
      down: 48
    },
    replies: [
      {
        id: 'reply3-1',
        postId: '1',
        authorId: 'user11',
        authorUsername: 'ts_defender',
        authorAvatar: 'https://i.pravatar.cc/150?img=14',
        content: 'Have to disagree. Type safety catches so many bugs at compile time!',
        votesCount: {
          up: 78,
          down: 5
        },
        replies: [],
        createdAt: '2025-10-19T10:30:00Z'
      }
    ],
    createdAt: '2025-10-19T10:15:00Z'
  },
  {
    id: 'comment4',
    postId: '1',
    authorId: 'user12',
    authorUsername: 'newbie_coder',
    authorAvatar: 'https://i.pravatar.cc/150?img=15',
    content: 'Just started learning TypeScript last week. This guide is exactly what I needed!',
    votesCount: {
      up: 56,
      down: 0
    },
    replies: [],
    createdAt: '2025-10-19T11:00:00Z'
  }
];