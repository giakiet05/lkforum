export interface RestrictedUser {
  id: number;
  username: string;
  avatar: string;
  type: "banned" | "muted";
  duration: string;
  durationEnd: string;
  date?: string;
  reason?: string;
  note?: string;
  ruleViolated?: string;
}

export const mockRestrictedUsers: RestrictedUser[] = [
  {
    id: 1,
    username: "user1",
    avatar: "/avatar.jpg",
    type: "banned",
    duration: "24 hour left",
    durationEnd: "2025-10-25T12:00:00Z",
    date: "18/9/2025",
    reason: "reason1",
    note: "rule1",
    ruleViolated: "rule1",
  },
  {
    id: 2,
    username: "user2",
    avatar: "/avatar.jpg",
    type: "banned",
    duration: "7 days left",
    durationEnd: "2025-10-31T12:00:00Z",
    date: "24/10/2025",
    reason: "Spam posting",
    note: "No spam rule violation",
    ruleViolated: "No spam",
  },
  {
    id: 3,
    username: "user3",
    avatar: "/avatar.jpg",
    type: "muted",
    duration: "24 hour left",
    durationEnd: "2025-10-25T15:30:00Z",
    note: "rule1",
  },
  {
    id: 4,
    username: "user4",
    avatar: "/avatar.jpg",
    type: "muted",
    duration: "3 days left",
    durationEnd: "2025-10-27T10:00:00Z",
    note: "Excessive off-topic comments",
  },
];
