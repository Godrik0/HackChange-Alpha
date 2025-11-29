import { Component, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzSliderModule } from 'ng-zorro-antd/slider';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-clients-filter',
  standalone: true,
  imports: [
    CommonModule,
    NzCardModule,
    NzFormModule,
    NzInputModule,
    NzSelectModule,
    NzSliderModule,
    FormsModule
  ],
  templateUrl: './clients-filter.component.html'
})
export class ClientsFilterComponent {
  @Output() filterChange = new EventEmitter<any>();

  filter = {
    search: '',
    profitRange: [40000, 60000]
  };

  onFilterChange() {
    this.filterChange.emit(this.filter);
  }
}
