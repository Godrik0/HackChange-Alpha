import { Component } from '@angular/core';

import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzCardModule } from 'ng-zorro-antd/card';
import { Client } from '@core/models/client';
import {NzDropDownModule} from "ng-zorro-antd/dropdown";
import { NzMenuModule } from 'ng-zorro-antd/menu';
import {ClientLineComponent} from "@features/clients/components/client-line/client-line.component";
import {NgForOf, NgOptimizedImage} from "@angular/common";
import {NzFlexDirective} from "ng-zorro-antd/flex";
import {NzOptionComponent, NzSelectComponent} from "ng-zorro-antd/select";
import {FormsModule} from "@angular/forms";
import {NzEmptyModule} from "ng-zorro-antd/empty";
import {NzInputDirective, NzInputSearchDirective} from "ng-zorro-antd/input";
import {RouterLink} from "@angular/router";

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
  ],
  templateUrl: './clients.component.html',
  styleUrls: ['./clients.component.scss'],
})
export class ClientsComponent {
  clients: Client[] = [
    { id: 1, firstName: 'Сергей', lastName: 'Орлов', middleName: 'Николаевич', birthDate: '06-01-2005', incoming: 123456 },
    { id: 2, firstName: 'Роман', lastName: 'Орлов', middleName: 'Николаевич', birthDate: '06-01-2005', incoming: 123456 },
    { id: 3, firstName: 'Олег', lastName: 'Орлов', middleName: 'Николаевич', birthDate: '06-01-2005' },
  ];
  displayedClients: Client[];
  inputValue: string = '';
  selectedRange = 60000;
  rangeOptions = [60000, 70000, 100000];

  constructor() {
    this.displayedClients = this.clients;
  }

  get displayedClientsIsEmpty() {
    return this.displayedClients.length === 0;
  }

  changeLimit() {
    const range = Number(this.selectedRange);
    this.displayedClients = this.clients.filter(client => !!client.incoming && client.incoming <= range)
  }

  onChange(event: any) {

  }
}
