import {Component, Input} from '@angular/core';
import {NzCardComponent} from "ng-zorro-antd/card";
import {NzIconDirective} from "ng-zorro-antd/icon";

@Component({
  selector: 'app-action-card',
  imports: [
    NzCardComponent,
    NzIconDirective
  ],
  templateUrl: './action-card.component.html',
  styleUrl: './action-card.component.less',
})
export class ActionCardComponent {
  @Input() recommendation!: string;
}
