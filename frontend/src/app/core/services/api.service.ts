import { Injectable, PLATFORM_ID, Inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { isPlatformBrowser } from '@angular/common';
import { environment } from "@env/environment";
import { Observable } from "rxjs";

@Injectable({ providedIn: 'root' })
export class ApiService {
  private baseUrl: string;

  constructor(
    private http: HttpClient,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {
    // В браузере используем localhost, на сервере - внутренний Docker URL
    if (isPlatformBrowser(this.platformId)) {
      this.baseUrl = 'http://localhost:8080/api';
    } else {
      // SSR: используем переменную окружения или дефолтный URL
      this.baseUrl = process.env['API_URL'] || 'http://backend:8080/api';
    }
  }

  get<T>(endpoint: string, params?: HttpParams): Observable<T> {
    return this.http.get<T>(`${this.baseUrl}/${endpoint}`, { params });
  }

  post<T>(endpoint: string, data: any) {
    return this.http.post<T>(`${this.baseUrl}/${endpoint}`, data);
  }
}
