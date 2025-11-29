import { Injectable } from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {environment} from "@env/environment";
import {Observable} from "rxjs";

@Injectable({ providedIn: 'root' })
export class ApiService {
  private baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  get<T>(endpoint: string, params?: HttpParams = {}): Observable<T> {
    return this.http.get<T>(`${this.baseUrl}/api/${endpoint}`, { params });
  }

  post<T>(endpoint: string, data: any) {
    return this.http.post<T>(`${this.baseUrl}/api/${endpoint}`, data);
  }
}
