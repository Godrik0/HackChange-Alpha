import { Injectable } from '@angular/core';
import {Observable} from "rxjs";
import {Client} from "@core/models/client";
import {ApiService} from "@core/services/api.service";
import {HttpParams} from "@angular/common/http";

type SearchData = Omit<Client, 'id' | 'incoming'>

@Injectable({ providedIn: 'root' })
export class DataService {

  constructor(private api: ApiService) {}

  getClients(): Observable<Client[]> {
    return this.api.get('clients')
  }

  getClient(clientID: number): Observable<Client> {
    return this.api.get(`clients/${clientID}`);
  }

  searchClient(searchData: SearchData): Observable<Client[]> {
    let params = new HttpParams();

    Object.entries(searchData).forEach(([key, value]) => {
      if (value != null) {
        params = params.set(key, value.toString());
      }
    });

    return params.keys().length !== 0 ? this.api.get('clients/search', params) : new Observable<Client[]>();
  }
}
