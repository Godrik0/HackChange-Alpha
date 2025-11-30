import {ChangeDetectorRef, Component, OnInit, OnDestroy, PLATFORM_ID, Inject} from '@angular/core';

import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzCardModule } from 'ng-zorro-antd/card';
import { Client } from '@core/models/client';
import {NzDropDownModule} from "ng-zorro-antd/dropdown";
import { NzMenuModule } from 'ng-zorro-antd/menu';
import {ClientLineComponent} from "@features/clients/components/client-line/client-line.component";
import {AsyncPipe, NgForOf, NgOptimizedImage, isPlatformBrowser} from "@angular/common";
import {NzFlexDirective} from "ng-zorro-antd/flex";
import {NzOptionComponent, NzSelectComponent} from "ng-zorro-antd/select";
import {FormsModule} from "@angular/forms";
import {NzEmptyModule} from "ng-zorro-antd/empty";
import {NzInputDirective, NzInputSearchDirective} from "ng-zorro-antd/input";
import {RouterLink} from "@angular/router";
import {DataService} from "@core/services/data.service";
import {Subject, takeUntil} from "rxjs";

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
export class ClientsComponent implements OnInit, OnDestroy {
  allClients: Client[] = [];
  displayedClients: Client[] = [];
  inputValue: string = '';
  selectedRange = 60000;
  rangeOptions = [60000, 70000, 100000];
  private destroy$ = new Subject<void>();

  constructor(
    private dataService: DataService,
    private cdr: ChangeDetectorRef,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit() {
    // Загружаем данные только в браузере
    if (isPlatformBrowser(this.platformId)) {
      this.dataService.getClients()
        .pipe(takeUntil(this.destroy$))
        .subscribe({
          next: (clients: Client[]) => {
            this.allClients = clients;
            this.displayedClients = clients;
            this.cdr.detectChanges();
          },
          error: (error) => {
            console.error('Error loading clients:', error);
          }
        });
    }
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  get displayedClientsIsEmpty() {
    return this.displayedClients.length === 0;
  }

  changeLimit() {
    const range = Number(this.selectedRange);
    this.displayedClients = this.allClients.filter(client => !!client.income && client.income <= range);
  }
}
