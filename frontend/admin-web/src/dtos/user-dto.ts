export interface UserResponse {
  id: string;
  username: string;
  email: string;
  avatar?: string;
  bio?: string;
  created_at: string;
  is_banned?: boolean;
}

export interface PaginatedUsersResponse {
  users: UserResponse[];
  total: number;
  page: number;
  page_size: number;
}
