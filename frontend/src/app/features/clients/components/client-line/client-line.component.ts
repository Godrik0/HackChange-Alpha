import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzFlexModule } from 'ng-zorro-antd/flex';
import { NzAvatarModule } from 'ng-zorro-antd/avatar';
import { Client } from '@core/models/client';
import {getFullName} from "@core/utils/utils";
import {NzButtonComponent} from "ng-zorro-antd/button";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-client-line',
  standalone: true,
  imports: [CommonModule, NzFlexModule, NzAvatarModule, NzButtonComponent, RouterLink],
  templateUrl: './client-line.component.html'
})
export class ClientLineComponent {
  @Input() client!: Client;
  protected readonly getFullName = getFullName;
}
