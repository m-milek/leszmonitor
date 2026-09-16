export interface Timestamps {
  createdAt: Date;
  updatedAt: Date;
}

export interface Pagination {
  page: number;
  perPage: number;
}

export interface ApiErrorDetails {
  message: string;
}

export interface ApiError {
  error: ApiErrorDetails;
  status: number;
}
