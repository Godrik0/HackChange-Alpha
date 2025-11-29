import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzFlexModule } from 'ng-zorro-antd/flex';
import { NzAvatarModule } from 'ng-zorro-antd/avatar';
import { Client } from '@core/models/client';

@Component({
  selector: 'app-client-line',
  standalone: true,
  imports: [CommonModule, NzFlexModule, NzAvatarModule],
  templateUrl: './client-line.component.html'
})
export class ClientLineComponent {
  @Input() client!: Client;
}
