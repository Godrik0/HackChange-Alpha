export interface Scoring {
  id: number;
  first_name: string;
  last_name: string;
  middle_name?: string;
  birth_date: string;
  income?: number;
  predict_income: number;
  credit_limit: number;
  max_credit_limit?: number;
  recommendations: string[];
  positive_factors: string[];
  negative_factors: string[];
}
