export interface GetLogsResponse {
  logs: Log[];
  next_token: string;
}

export interface Log {
  id: string;
  tag: string;
  timestamp: string;
  data: any;
}
