import {ChangeDetectorRef, Component, OnInit} from '@angular/core';

import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzCardModule } from 'ng-zorro-antd/card';
import { Client } from '@core/models/client';
import {NzDropDownModule} from "ng-zorro-antd/dropdown";
import { NzMenuModule } from 'ng-zorro-antd/menu';
import {ClientLineComponent} from "@features/clients/components/client-line/client-line.component";
import {AsyncPipe, NgForOf, NgOptimizedImage} from "@angular/common";
import {NzFlexDirective} from "ng-zorro-antd/flex";
import {NzOptionComponent, NzSelectComponent} from "ng-zorro-antd/select";
import {FormsModule} from "@angular/forms";
import {NzEmptyModule} from "ng-zorro-antd/empty";
import {NzInputDirective, NzInputSearchDirective} from "ng-zorro-antd/input";
import {RouterLink} from "@angular/router";
import {DataService} from "@core/services/data.service";
import {Observable} from "rxjs";

@Component({
  selector: 'app-clients-page',
  standalone: true,
  imports: [
    NzLayoutModule,
    NzCardModule,
    NzDropDownModule,
    ClientLineComponent,
    NgForOf,
    NzFlexDirective,
    NzMenuModule,
    NzSelectComponent,
    NzOptionComponent,
    FormsModule,
    NzEmptyModule,
    NzInputSearchDirective,
    NzInputDirective,
    NgOptimizedImage,
    RouterLink,
    AsyncPipe,
  ],
  templateUrl: './clients.component.html',
  styleUrls: ['./clients.component.scss'],
})
export class ClientsComponent implements OnInit {
  clients$: Observable<Client[]> = [] as unknown as Observable<Client[]>;
  displayedClients: Client[] = [];
  inputValue: string = '';
  selectedRange = 60000;
  rangeOptions = [60000, 70000, 100000];

  constructor(private dataService: DataService,  private cdr: ChangeDetectorRef) {}

  ngOnInit() {
    this.clients$ = this.dataService.getClients();
    this.clients$.subscribe(
      (clients: Client[]) => {
        this.displayedClients = clients;
        this.cdr.detectChanges();
      }
    );
  }

  get displayedClientsIsEmpty() {
    return this.displayedClients.length === 0;
  }

  changeLimit() {
    let clients: Client[] = [];
    const range = Number(this.selectedRange);
    this.clients$.subscribe(
      (client: Client[]) => {
         clients = client
      }
    );
    this.displayedClients = clients.filter(client => !!client.incoming && client.incoming <= range)
  }
}
