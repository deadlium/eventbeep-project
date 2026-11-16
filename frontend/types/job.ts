export interface Job {
  id: string;
  hash_key: string;
  title: string;
  company: string;
  role: string | null;
  city: string | null;
  source_url: string;
  created_at: string;
}

export interface Pagination {
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface PaginatedJobsResponse {
  data: Job[];
  pagination: Pagination;
}

export interface JobFilterParams {
  role?: string;
  city?: string;
  search?: string;
  page?: number;
  limit?: number;
}
