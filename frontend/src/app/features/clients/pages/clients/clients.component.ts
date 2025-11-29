import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzCardModule } from 'ng-zorro-antd/card';
import { ClientsListComponent } from '../../components/clients-list/clients-list.component';
import { Client } from '@core/models/client';

@Component({
  selector: 'app-clients-page',
  standalone: true,
  imports: [
    CommonModule,
    NzLayoutModule,
    NzCardModule,
    ClientsListComponent
  ],
  templateUrl: './clients.component.html'
})
export class ClientsComponent {
  clients: Client[] = [
    { id: 1, name: 'Алексей Петров', description: 'Постоянный клиент', profit: 70000 },
    { id: 2, name: 'Мария Сидорова', description: 'Новый клиент', profit: 65000 },
    { id: 3, name: 'Иван Козлов', description: 'VIP клиент', profit: 60000 },
    { id: 4, name: 'Елена Новикова', description: 'Корпоративный клиент', profit: 55000 },
    { id: 5, name: 'Дмитрий Волков', description: 'Частный клиент', profit: 50000 }
  ];

  onFilterChange(filter: any) {
    console.log('Filter changed:', filter);
    // Здесь будет логика фильтрации
  }
}
