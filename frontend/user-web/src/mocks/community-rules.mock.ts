export interface CommunityRule {
  id: number;
  name: string;
  description: string;
  reportReason?: string;
  createdAt: string;
}

export const mockCommunityRules: CommunityRule[] = [
  {
    id: 1,
    name: "rule1",
    description: "we dont talk about this group",
    createdAt: "2024-01-15T10:30:00Z",
  },
  {
    id: 2,
    name: "Be respectful",
    description: "Treat all members with respect. No harassment, hate speech, or personal attacks.",
    createdAt: "2024-01-20T14:15:00Z",
  },
  {
    id: 3,
    name: "No spam",
    description: "Do not post spam, promotional content, or excessive self-promotion.",
    createdAt: "2024-02-01T09:00:00Z",
  },
  {
    id: 4,
    name: "Stay on topic",
    description: "Keep posts relevant to the community. Off-topic content will be removed.",
    createdAt: "2024-02-10T16:45:00Z",
  },
  {
    id: 5,
    name: "No NSFW content",
    description: "This is a SFW community. NSFW content is not allowed.",
    createdAt: "2024-02-15T11:20:00Z",
  },
];
