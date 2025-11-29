import { Component, Input } from '@angular/core';

import { ClientLineComponent } from '../client-line/client-line.component';
import { Client } from '@core/models/client';
import {NzFlexDirective} from "ng-zorro-antd/flex";

@Component({
  selector: 'app-clients-list',
  standalone: true,
  imports: [ClientLineComponent, NzFlexDirective],
  templateUrl: './clients-list.component.html'
})
export class ClientsListComponent {
  @Input() clients: Client[] = [];
}
