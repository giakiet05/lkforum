export interface JoinedCommunity {
  name: string;
  members: string;
  status: "Favorited" | "Subscribed";
  icon: string;
}

// Mock data cho các community đã join (for CreatePostModal community selector)
export const mockJoinedCommunities: JoinedCommunity[] = [
  {
    name: "3amjokes",
    members: "2,349,727 members",
    status: "Favorited",
    icon: "/LKlogo.jpg",
  },
  {
    name: "anime",
    members: "14,043,102 members",
    status: "Favorited",
    icon: "/LKlogo.jpg",
  },
  {
    name: "Animesuggiest",
    members: "1,064,940 members",
    status: "Favorited",
    icon: "/LKlogo.jpg",
  },
  {
    name: "30PlusSkinCare",
    members: "2,323,937 members",
    status: "Subscribed",
    icon: "/LKlogo.jpg",
  },
  {
    name: "acne",
    members: "1,731,543 members",
    status: "Subscribed",
    icon: "/LKlogo.jpg",
  },
  {
    name: "AdviceAnimals",
    members: "9,904,334 members",
    status: "Subscribed",
    icon: "/LKlogo.jpg",
  },
  {
    name: "AmItheAsshole",
    members: "15,234,567 members",
    status: "Subscribed",
    icon: "/LKlogo.jpg",
  },
];
