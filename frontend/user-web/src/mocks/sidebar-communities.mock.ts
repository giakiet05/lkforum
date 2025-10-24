export interface SidebarCommunity {
  id: string;
  name: string;
  icon: string;
  isFavorite: boolean;
  memberCount: number;
}

// Mock data for sidebar communities list
export const mockSidebarCommunities: SidebarCommunity[] = [
  {
    id: "1",
    name: "3amjokes",
    icon: "🌙",
    isFavorite: true,
    memberCount: 125000,
  },
  {
    id: "2",
    name: "anime",
    icon: "👧",
    isFavorite: true,
    memberCount: 850000,
  },
  {
    id: "3",
    name: "Animesuggest",
    icon: "💭",
    isFavorite: true,
    memberCount: 45000,
  },
  {
    id: "4",
    name: "30PlusSkinCare",
    icon: "🧴",
    isFavorite: false,
    memberCount: 23000,
  },
  {
    id: "5",
    name: "acne",
    icon: "🩺",
    isFavorite: false,
    memberCount: 18000,
  },
];
