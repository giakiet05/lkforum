export interface Community {
  id: string;
  name: string;
  description?: string;
  avatar?: string;
  banner?: string;
  settings: {
    isPrivate: boolean;
    allowPosts: boolean;
    allowComments: boolean;
    allowMedia: boolean;
    postRequireApproval: boolean;
    joinRequireApproval: boolean;
    maxPostLength?: number;
  };
  moderators: Array<{
    userId: string;
    username: string;
    assignedAt: string;
  }>;
  memberCount: number;
  postCount: number;
  createdAt: string;
  createdById: string;
  createdByName: string;
  createdByAvatar?: string;
}

export const mockCommunities: Community[] = [
  {
    id: 'tech',
    name: 'Technology',
    description: 'A community for tech enthusiasts and professionals',
    avatar: 'https://picsum.photos/200/200?tech',
    banner: 'https://picsum.photos/1200/300?tech',
    settings: {
      isPrivate: false,
      allowPosts: true,
      allowComments: true,
      allowMedia: true,
      postRequireApproval: false,
      joinRequireApproval: false
    },
    moderators: [
      {
        userId: 'mod1',
        username: 'techmoderator',
        assignedAt: '2025-01-01T00:00:00Z'
      }
    ],
    memberCount: 15000,
    postCount: 2500,
    createdAt: '2025-01-01T00:00:00Z',
    createdById: 'admin1',
    createdByName: 'AdminUser',
    createdByAvatar: 'https://i.pravatar.cc/150?img=4'
  },
  {
    id: 'programming',
    name: 'Programming',
    description: 'Share your coding experiences and get help',
    avatar: 'https://picsum.photos/200/200?programming',
    banner: 'https://picsum.photos/1200/300?programming',
    settings: {
      isPrivate: false,
      allowPosts: true,
      allowComments: true,
      allowMedia: true,
      postRequireApproval: false,
      joinRequireApproval: false,
      maxPostLength: 5000
    },
    moderators: [
      {
        userId: 'mod2',
        username: 'codeguru',
        assignedAt: '2025-02-01T00:00:00Z'
      }
    ],
    memberCount: 10000,
    postCount: 1800,
    createdAt: '2025-02-01T00:00:00Z',
    createdById: 'admin2',
    createdByName: 'TechAdmin',
    createdByAvatar: 'https://i.pravatar.cc/150?img=5'
  },
  {
    id: 'photography',
    name: 'Photography',
    description: 'Share your best shots and photography tips',
    avatar: 'https://picsum.photos/200/200?photo',
    banner: 'https://picsum.photos/1200/300?photo',
    settings: {
      isPrivate: false,
      allowPosts: true,
      allowComments: true,
      allowMedia: true,
      postRequireApproval: true,
      joinRequireApproval: false
    },
    moderators: [
      {
        userId: 'mod3',
        username: 'photopro',
        assignedAt: '2025-03-01T00:00:00Z'
      }
    ],
    memberCount: 8000,
    postCount: 1200,
    createdAt: '2025-03-01T00:00:00Z',
    createdById: 'admin3',
    createdByName: 'PhotoAdmin',
    createdByAvatar: 'https://i.pravatar.cc/150?img=6'
  }
];