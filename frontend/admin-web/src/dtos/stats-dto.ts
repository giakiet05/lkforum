export interface PlatformOverview {
  total_users: number;
  active_users: number;
  banned_users: number;
  total_communities: number;
  active_communities: number;
  banned_communities: number;
  total_posts: number;
  total_comments: number;
  pending_reports: number;
}

export interface UserStats {
  total_users: number;
  new_users: number;
  active_users: number;
  banned_users: number;
  growth_rate: number;
  users_by_day?: DailyCount[];
}

export interface ContentStats {
  total_posts: number;
  new_posts: number;
  total_comments: number;
  new_comments: number;
  total_communities: number;
  growth_rate: number;
  posts_by_day?: DailyCount[];
}

export interface DailyCount {
  date: string;
  count: number;
}

export interface GetUserStatsQuery {
  period?: "day" | "week" | "month" | "year";
}

export interface GetContentStatsQuery {
  period?: "day" | "week" | "month" | "year";
}
