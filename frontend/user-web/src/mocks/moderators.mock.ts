export interface Moderator {
  id: number;
  username: string;
  avatar: string;
  permissions: string;
  canEdit: boolean;
  joinedDate: string;
}

export interface ApprovedUser {
  id: number;
  username: string;
  avatar: string;
  joinedDate: string;
}

export const mockModerators: Moderator[] = [
  {
    id: 1,
    username: "u/DrinkConsistent9498",
    avatar: "/avatar.jpg",
    permissions: "Everything",
    canEdit: false,
    joinedDate: "1:17 AM\nSep 20, 2025",
  },
  {
    id: 2,
    username: "u/Agitated_Relief6165",
    avatar: "/avatar.jpg",
    permissions: "Everything",
    canEdit: true,
    joinedDate: "2:30 PM\nOct 15, 2025",
  },
];

export const mockApprovedUsers: ApprovedUser[] = [
  {
    id: 1,
    username: "user1",
    avatar: "/avatar.jpg",
    joinedDate: "18/9/2025",
  },
  {
    id: 2,
    username: "user2",
    avatar: "/avatar.jpg",
    joinedDate: "20/10/2025",
  },
];
